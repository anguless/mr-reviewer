package service

type RandomGenerator interface {
	Intn(n int) int
}
