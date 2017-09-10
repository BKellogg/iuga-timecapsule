package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BKellogg/iuga-timecapsule/tc-server/handlers"
	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultHost = ""
	defaultPort = "80"

	capsulePath = "/capsule"
)

func main() {

	// get environment variables
	host := os.Getenv("HOST")
	if len(host) == 0 {
		fmt.Printf("no HOST environment variable set, defaulting to %s\n", defaultHost)
		host = defaultHost
	}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		fmt.Printf("no PORT environment variable set, defaulting to %s\n", defaultPort)
		port = defaultPort
	}
	mysqlAddr := os.Getenv("MYSQLADDR")
	if len(mysqlAddr) == 0 {
		log.Fatal("no MYSQLADDR environment variable set, please set a MYSQLADDR")
	}
	mysqlPass := os.Getenv("MYSQLPASS")
	if len(mysqlPass) == 0 {
		log.Fatal("no MYSQLPASS environment variable set, please set a MYSQLPASS")
	}
	mysqlDB := os.Getenv("MYSQLDB")
	if len(mysqlDB) == 0 {
		log.Fatal("no MYSQLPASS environment variable set, please set a MYSQLDB")
	}

	// open a connection to the mysql db
	db := openMySQLConnectionOrStop("root", mysqlPass, mysqlDB)
	ensureTableOrStop(db, mysqlDB)

	// Create the handler context
	ctx := handlers.HandlerContext{
		DB: db,
	}

	// create a new serve mux and attach the handlers
	mux := http.NewServeMux()
	mux.HandleFunc(capsulePath, ctx.HandleNewTimeCapsule)

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("TimeCapsule server listening on %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// opens a connection to the mysql database with the given credentials
// fatally logs if there was an error opening the connection
func openMySQLConnectionOrStop(user, password, dbName string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", user, password, dbName))
	if err != nil {
		log.Fatal("error opening connection to mysql: " + err.Error())
	}
	return db
}

// ensure that the given sql db has a Capsules table
// fatally logs if there was an error creating the table
func ensureTableOrStop(db *sql.DB, dbName string) {
	_, err := db.Exec("SELECT 1 FROM Capsules LIMIT 1")
	if err != nil {
		fmt.Println(dbName + ".capsules does not exist; creating...")
		_, err := db.Exec(`CREATE TABLE Capsules (
									CapsuleID INT AUTO_INCREMENT PRIMARY KEY,
									NetID VARCHAR(25) NOT NULL,
									GradDate Date NOT NULL,
									Message VARCHAR(500) NOT NULL
								)`)
		if err != nil {
			log.Fatal("error creating Capsules table: " + err.Error())
		}
		fmt.Printf("successfully created %s.Capsules\n", dbName)
	}
}
