package main

import (
	"./config"
	"./handlers"
	"crypto/tls"
	"flag"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"os"
)

var configFile = flag.String("config", "", "YAML configuration file.")

func main() {
	flag.Parse() // parses the logging flags.
	if *configFile == "" {
		println("No config file specified, exiting.")
		os.Exit(1)
	}

	conf, err := config.LoadConfig(*configFile)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	c := irc.SimpleClient(conf.BotNick)
	if conf.SSL {
		c.SSL = conf.SSL
		cert, err := tls.LoadX509KeyPair(conf.CertFile, conf.KeyFile)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		c.SSLConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
		c.SSLConfig.BuildNameToCertificate()
		c.SSLConfig.InsecureSkipVerify = true
	}

	c.AddHandler(irc.CONNECTED,
		handlers.CreateConnectedHandler(conf))

	quit := make(chan bool)
	c.AddHandler(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	c.AddHandler("PRIVMSG",
		handlers.CreateWrongNickHandler(conf))
	if conf.Password != "" {
		if err := c.Connect(fmt.Sprintf("%s:%d", conf.Network, conf.Port), conf.Password); err != nil {
			fmt.Printf("Connection error: %s\n", err.Error())
		}
	} else {
		println(fmt.Sprintf("%s:%d", conf.Network, conf.Port))
		if err := c.Connect(fmt.Sprintf("%s:%d", conf.Network, conf.Port)); err != nil {
			fmt.Printf("Connection error: %s\n", err.Error())
		}
	}

	<-quit
}
