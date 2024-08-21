package database

type HealthCheck struct {
	ID            uint   `json:"id"`
	App_name      string `json:"app_name"`
	Failure_count uint   `json:"failure_count"`
	Success_count uint   `json:"success_count"`
	Last_failure  string `json:"last_failure"`
	Last_success  string `json:"last_success"`
	Created_at    string `json:"created_at"`
}

type DatabaseInfo struct {
	host          string
	port          string
	user          string
	password      string
	database_name string
}
