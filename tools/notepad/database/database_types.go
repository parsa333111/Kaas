package database

type Note struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
}

type DatabaseInfo struct {
	host          string
	port          string
	user          string
	password      string
	database_name string
}
