package main

import (
	"database/sql"
	logging "gorestboilerplate/utils/logger"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kelseyhightower/envconfig"
	goLogging "github.com/op/go-logging"
)

var logger goLogging.Logger
var db *sql.DB

var opt struct {
	LogLevel string `default:"INFO" split_words:"true"`
	Env      string `default:"prod"`
	DBUrl    string `default:"root:password1@tcp(127.0.0.1:3306)/test"`
	Port     int    `default:"8080" split_words:"true"`
}

func main() {

	logger.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	logger.Info("Restful API starting...")

	// Configuration options from environment variables
	logger.Info("Parsing Options...")
	err := envconfig.Process("Rest", &opt)
	if err != nil {
		logger.Fatalf("Failed to parse command line arguments: %s", err.Error())
	}

	logger = *logging.New("rest", opt.LogLevel)

	db, err = sql.Open("mysql", opt.DBUrl)
	if err != nil {
		logger.Fatalf("Error starting mysql: %s", err.Error())
	}
	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// Launch HTTP Server, used for Admin reasons only
	logger.Infof("Starting server on port %d...", opt.Port)
	setupRoutes(strconv.Itoa(opt.Port))
}
