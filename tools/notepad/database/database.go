package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var DBInfo *DatabaseInfo

func loadDatabaseConnectionConfig() {
	host, ok := os.LookupEnv("PQ_HOST")
	if !ok {
		panic("Missing enviroment variable PQ_HOST.")
	}

	port, ok := os.LookupEnv("PQ_PORT")
	if !ok {
		panic("Missing enviroment variable PQ_PORT.")
	}

	user, ok := os.LookupEnv("PQ_USER")
	if !ok {
		panic("Missing enviroment variable PQ_USER.")
	}

	password, ok := os.LookupEnv("PQ_PASSWORD")
	if !ok {
		panic("Missing enviroment variable PQ_PASSWORD.")
	}

	database_name, ok := os.LookupEnv("PQ_DBNAME")
	if !ok {
		panic("Missing enviroment variable PQ_DBNAME.")
	}

	DBInfo = &DatabaseInfo{
		host,
		port,
		user,
		password,
		database_name,
	}
}

func connectToDatabase() {
	var info_string string = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DBInfo.host, DBInfo.port, DBInfo.user, DBInfo.password, DBInfo.database_name)

	var err error
	DB, err = sql.Open("postgres", info_string)
	if err != nil {
		panic(fmt.Sprintf("Error: Failed to connect to the database %s with error: %s", info_string, err))
	}
}

func createNotesTable() {
	query, err := DB.Prepare(`
		CREATE TABLE IF NOT EXISTS Notes (
		id integer PRIMARY KEY generated always as identity, 
		description varchar(30) not null);`,
	)
	if err != nil {
		panic(err)
	}

	_, err = query.Exec()
	if err != nil {
		panic(err)
	}
}

func Initalize() {
	loadDatabaseConnectionConfig()
	connectToDatabase()
	createNotesTable()
}

func AddNote(description string) bool {
	_, err := DB.Exec(`
		INSERT INTO
		Notes(description)
		VALUES($1);`,
		description)

	if err != nil {
		log.Println("Error:", err)
		return false
	}

	return true
}

func GetNotes() ([]Note, bool) {
	rows, err := DB.Query(`
		SELECT * 
		FROM Notes;`,
	)

	if err != nil {
		log.Println("Error:", err)
		return []Note{}, false
	}

	var notes []Note

	for rows.Next() {
		var note Note

		if err := rows.Scan(
			&note.ID,
			&note.Description,
		); err != nil {
			log.Println("Error:", err)
			return []Note{}, false
		} else {
			notes = append(notes, note)
		}
	}

	return notes, true
}
