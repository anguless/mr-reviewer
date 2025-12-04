package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/anguless/mr-reviewer/internal/api"
	"github.com/anguless/mr-reviewer/internal/config"
	"github.com/anguless/mr-reviewer/internal/db"
	"github.com/anguless/mr-reviewer/internal/migrator"
	"github.com/anguless/mr-reviewer/internal/repository"
	"github.com/anguless/mr-reviewer/internal/service"
	mrV1 "github.com/anguless/mr-reviewer/pkg/openapi/mr/v1"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	dbConn, err := db.NewDbPool(ctx, cfg.DbConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer dbConn.Close()

	err = dbConn.Ping(ctx)
	if err != nil {
		log.Fatalf("База данных недоступна: %v\n", err)
		return
	}

	migratorRunner := migrator.NewMigrator(dbConn.ToSqlDB(), cfg.MigrationConfig.MigrationsDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v\n", err)
		return
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	repo := repository.NewRepository(dbConn)

	srvc := service.NewService(repo, rnd)

	mrHandler := api.NewMrHandler(srvc)

	mrServer, err := mrV1.NewServer(mrHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", mrServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("0.0.0.0", cfg.AppConfig.AppPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", cfg.AppConfig.AppPort)
		log.Printf("Адрес: %s\n", server.Addr)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
