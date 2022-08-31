package main

import (
	"Edtech_Golang/api"
	db "Edtech_Golang/db/sqlc"
	"Edtech_Golang/util"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func main() {
	path := flag.String("./", "config", "config file location")
	writeToFile := flag.Bool("f", false, "write logs to file")
	flag.Parse()
	config, err := util.LoadConfig(*path)
	if err != nil {
		panic("cannot load  config:" + err.Error())
	}

	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)

	if *writeToFile {
		f, err := os.OpenFile("/var/log/edtech.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	}

	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name)

	conn, err := sql.Open("mysql", str)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn, &config)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.Server.Listen)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
