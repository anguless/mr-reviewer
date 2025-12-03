package model

type User struct {
	ID       string `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
}
