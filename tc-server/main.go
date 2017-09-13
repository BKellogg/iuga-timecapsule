package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	gmail "google.golang.org/api/gmail/v1"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/BKellogg/iuga-timecapsule/tc-server/handlers"
	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultHost = ""
	defaultPort = "443"

	capsulePath = "/capsule"
)

func main() {

	// HOSTING VARS
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
	tlsCert := os.Getenv("TLSCERT")
	if len(tlsCert) == 0 {
		log.Fatal("no TLSCERT environment variable set, please set a TLSCERT")
	}
	tlsKey := os.Getenv("TLSKEY")
	if len(tlsKey) == 0 {
		log.Fatal("no TLSKEY environment variable set, please set a TLSKEY")
	}

	// MYSQL VARS
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
		log.Fatal("no MYSQLDB environment variable set, please set a MYSQLDB")
	}

	// Reads in our credentials
	secret, err := ioutil.ReadFile("./secret/client_secret.json")
	if err != nil {
		log.Printf("Error: %v", err)
	}

	// Creates a oauth2.Config using the secret
	// The second parameter is the scope, in this case we only want to send email
	conf, err := google.ConfigFromJSON(secret, gmail.GmailSendScope)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	googleToken := os.Getenv("GOOGLE_REFRESH_TOKEN")
	if len(googleToken) == 0 {
		log.Fatal("no GOOGLE_REFRESH_TOKEN environment variable set, please set a GOOGLE_REFRESH_TOKEN")
	}
	var token *oauth2.Token
	if !token.Valid() {
		token = &oauth2.Token{RefreshToken: googleToken}
	}

	// Exchange the auth code for an access token
	client := conf.Client(oauth2.NoContext, token)

	gmailService, err := gmail.New(client)
	if err != nil {
		log.Fatal("error creating gmail client: " + err.Error())
	}

	// open a connection to the mysql db
	db := openMySQLConnectionOrStop("root", mysqlPass, mysqlAddr, mysqlDB)
	ensureTableOrStop(db, mysqlDB)

	// Create the handler context
	ctx := handlers.HandlerContext{
		DB:           db,
		GmailService: gmailService,
	}

	// create a new serve mux and attach the handlers
	mux := http.NewServeMux()
	mux.HandleFunc(capsulePath, ctx.HandleNewTimeCapsule)

	addr := fmt.Sprintf("%s:%s", host, port)
	fmt.Printf("TimeCapsule server listening on %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, handlers.NewCORSHandler(mux)))
}

// opens a connection to the mysql database with the given credentials
// fatally logs if there was an error opening the connection
func openMySQLConnectionOrStop(user, password, addr, dbName string) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, addr, dbName))
	if err != nil {
		log.Fatal("error opening connection to mysql: " + err.Error())
	}
	err = db.Ping()
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
