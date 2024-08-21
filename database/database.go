package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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

	time.Sleep(70 * time.Second)

	var err error
	DB, err = sql.Open("postgres", info_string)
	if err != nil {
		panic(fmt.Sprintf("Error: Failed to connect to the database %s with error: %s", info_string, err))
	}
}

func Initalize() {
	loadDatabaseConnectionConfig()
	connectToDatabase()
}

func GetDeploymentHealth(app_name string) (HealthCheck, bool) {
	var healthCheck HealthCheck

	err := DB.QueryRow(`
		SELECT * 
		FROM HealthCheck
		WHERE app_name = $1;`,
		app_name).
		Scan(&healthCheck.ID,
			&healthCheck.App_name,
			&healthCheck.Failure_count,
			&healthCheck.Success_count,
			&healthCheck.Last_failure,
			&healthCheck.Last_success,
			&healthCheck.Created_at)

	if err != nil {
		log.Println("Error:", err)
		return healthCheck, false
	}

	return healthCheck, true
}
