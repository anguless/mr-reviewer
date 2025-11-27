package main

import (
	"context"
	"log"
	"net/http"

	"github.com/anguless/mr-reviewer/internal/api/handlers"
	"github.com/anguless/mr-reviewer/internal/config"
	"github.com/anguless/mr-reviewer/internal/db"
	"github.com/anguless/mr-reviewer/internal/migrator"
	"github.com/anguless/mr-reviewer/internal/repository"
	"github.com/anguless/mr-reviewer/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
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

	repo := repository.NewRepository(dbConn)

	srvc := service.NewService(repo)

	userHandler := &handlers.UserHandler{UserService: srvc.UserService}
	teamHandler := &handlers.TeamHandler{TeamService: srvc.TeamService, PRService: srvc.PrService}
	prHandler := &handlers.PRHandler{PRService: srvc.PrService}
	statsHandler := &handlers.StatisticsHandler{StatService: srvc.StatService}

	r := mux.NewRouter()

	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	r.HandleFunc("/api/v1/users/create", userHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/v1/users/{user_id}", userHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/v1/users/{user_id}", userHandler.UpdateUser).Methods("PUT")
	r.HandleFunc("/api/v1/users/{user_id}", userHandler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/api/v1/users/{user_id}/pull-requests", userHandler.GetUserPRs).Methods("GET")

	r.HandleFunc("/api/v1/team/add", teamHandler.CreateTeam).Methods("POST")
	r.HandleFunc("/api/v1/team/{team_id}", teamHandler.GetTeam).Methods("GET")
	r.HandleFunc("/api/v1/team/{team_id}", teamHandler.UpdateTeam).Methods("PUT")
	r.HandleFunc("/api/v1/team/{team_id}", teamHandler.DeleteTeam).Methods("DELETE")

	r.HandleFunc("/api/v1/pull-request/create", prHandler.CreatePR).Methods("POST")
	r.HandleFunc("/api/v1/pull-request", prHandler.GetAllPRs).Methods("GET")
	r.HandleFunc("/api/v1/pull-request/{pull_request_id}", prHandler.GetPR).Methods("GET")
	r.HandleFunc("/api/v1/pull-request/reassign", prHandler.ReassignReviewer).Methods("POST")
	r.HandleFunc("/api/v1/pull-request/merge", prHandler.MergePR).Methods("POST")

	r.HandleFunc("/api/v1/statistics", statsHandler.GetStatistics).Methods("GET")

	log.Printf("Starting server at %s", ":"+cfg.AppConfig.AppPort)

	if err := http.ListenAndServe(":"+cfg.AppConfig.AppPort, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
