package main

import (
	"os"

	"github.com/codegangsta/cli"

	"bitbucket.com/aria.pqstudio.pl-api/utils"
	"bitbucket.com/aria.pqstudio.pl-api/utils/db"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"github.com/kelseyhightower/envconfig"
	"github.com/op/go-logging"

	_ "github.com/go-sql-driver/mysql"

	"bitbucket.com/aria.pqstudio.pl-api/user/model"
	"bitbucket.com/aria.pqstudio.pl-api/user/service"

	oauth2Service "bitbucket.com/aria.pqstudio.pl-api/oauth2/service"
	"github.com/RangelReale/osin"
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

func main() {
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

	// configure MySQL
	conn := config.MySQLUser + ":" + config.MySQLPassword + "@tcp(" +
		config.MySQLHost + ":" + config.MySQLPort + ")/" + config.MySQLDatabase + "?charset=utf8&parseTime=true"
	err = db.Connect("mysql", conn)

	if err != nil {
		Log.Critical(err.Error())
	} /*else {*/
	//// TODO: ping database
	////Log.Info("Connection to database acquired: %s", conn)
	/*}*/

	cliApp := cli.NewApp()

	cliApp.Commands = []cli.Command{
		{
			Name:      "oauth2:client:create",
			ShortName: "o:c:c",
			Usage:     "create new oauth2 client",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 1 {
					println("Missing redirect uri")
					return
				}

				client := &osin.DefaultClient{
					Id:          utils.NewUUID(),
					Secret:      utils.NewUUID(),
					RedirectUri: c.Args()[0],
				}

				oauth2Service.CreateClient(client)

				println("Created client, Id:", client.Id, ",Secret:", client.Secret)
			},
		},
		{
			Name:      "user:add",
			ShortName: "u:a",
			Usage:     "add user user:add email password",
			Action: func(c *cli.Context) {
				if len(c.Args()) < 2 {
					println("Missing email or password")
					return
				}

				user := &model.User{
					Email:    c.Args()[0],
					Password: c.Args()[1],
				}

				err = service.CreateUser(user)
				if err != nil {
					println("Error creating user. Possibly user already exists?")
					return
				}

				println("Added user:", c.Args()[0])
			},
		},
	}
	cliApp.Name = "aria"
	cliApp.Usage = "cli for administrator"
	cliApp.Run(os.Args)
}
