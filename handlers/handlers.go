package handlers

import (
	"../config"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"log"
	"regexp"
)

type WatchList struct {
	Watch string
	Nick  string
}

func CreateConnectedHandler(conf *config.Config) irc.Handler {
	return func(conn *irc.Conn, line *irc.Line) {
		for _, channel := range conf.Channels {
			log.Printf("Joining channel: %s", channel.Name)
			conn.Join(channel.Name)
		}
	}
}

func CreateWrongNickHandler(conf *config.Config) irc.Handler {
	nicks := make(map[string]WatchList)
	for _, channel := range conf.Channels {
		nicks[channel.Name] = WatchList{
			Watch: channel.Watch,
			Nick:  channel.Nick,
		}
	}

	return func(conn *irc.Conn, line *irc.Line) {
		channel := line.Args[0]
		msg := line.Args[1]

		nickPair, ok := nicks[channel]
		// channel not watched
		if !ok {
			return
		}

		matched, err := regexp.MatchString(fmt.Sprintf("\\b%s\\b", nickPair.Watch), msg)
		if err != nil {
			log.Printf("%v", err)
		}
		if matched {
			conn.Privmsg(channel, fmt.Sprintf("%s: ^^^", nickPair.Nick))
		}
	}
}
