package core

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli/v3"
)

func messageCreate(session *discordgo.Session, messageCreate *discordgo.MessageCreate) {
	log.Printf("received message with content '%s'\n", messageCreate.Content)
}

func Run(ctx context.Context, cmd *cli.Command) error {
	botTokens := cmd.StringSlice("bot-token")

	sessions := []*discordgo.Session{}
	for _, botToken := range botTokens {
		session, err := discordgo.New("Bot " + botToken)
		if err != nil {
			log.Printf("couldn't create session: %s", err.Error())
			continue
		}

		httpTransport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient := &http.Client{Transport: httpTransport}
		session.Client = httpClient

		dialer := websocket.DefaultDialer
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		session.Dialer = dialer

		session.Identify.Intents = discordgo.IntentsGuildMessages

		session.AddHandler(messageCreate)

		err = session.Open()
		if err != nil {
			log.Printf("couldn't open session: %s", err.Error())
			continue
		}

		sessions = append(sessions, session)
	}

	if len(sessions) == 0 {
		return cli.Exit("couldn't create any sessions", 1)
	}

	log.Println("bots are running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	for _, session := range sessions {
		session.Close()
	}

	return nil
}
