package main

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils/db"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/op/go-logging"

	"flag"
	"github.com/zenazn/goji"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	app     = "aria"
	version = "0.1-alpha"
)

var (
	config Config
)

type Config struct {
	MySQLHost     string `envconfig:"mysql_host"`
	MySQLPort     string `envconfig:"mysql_port"`
	MySQLUser     string `envconfig:"mysql_user"`
	MySQLPassword string `envconfig:"mysql_password"`
	MySQLDatabase string `envconfig:"mysql_database"`
	RedisHost     string `envconfig:"redis_host"`
	RedisPort     int    `envconfig:"redis_port"`
	LogLevel      string `envconfig:"log_level"`
}

func init() {
	config.LogLevel = "debug"
	err := envconfig.Process(app, &config)

	if err != nil {
		Log.Critical(err.Error())
	}

	switch config.LogLevel {
	case "debug":
		SetLevel(logging.DEBUG)
	case "info":
		SetLevel(logging.INFO)
	case "notice":
		SetLevel(logging.NOTICE)
	case "warning":
		SetLevel(logging.WARNING)
	case "error":
		SetLevel(logging.ERROR)
	case "critical":
		SetLevel(logging.CRITICAL)
	}

	config.MySQLHost = "mariadb.service"
	conn := config.MySQLUser + ":" + config.MySQLPassword + "@tcp(" +
		config.MySQLHost + ":" + config.MySQLPort + ")/" + config.MySQLDatabase + "?charset=utf8"
	DB, err = sql.Open("mysql", conn)

	if err != nil {
		Log.Critical(err.Error())
	} else {
		Log.Info("Connection to database acquired: %s", conn)
	}
}

func main() {
	setupRoutes()

	flag.Set("bind", ":3000")
	goji.Serve()
}
