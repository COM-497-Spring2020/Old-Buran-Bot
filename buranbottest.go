package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	token      string
	configFile string
	dg         *discordgo.Session
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")

	flag.StringVar(&configFile, "c", "", "Config file")
	flag.Parse()
}

func main() {
	if token == "" {
		fmt.Printf("No token provided.")
	}

	dg, _ = discordgo.New("Bot " + token)
	dg.AddHandler(messageCreate)
	_ = dg.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}

	if strings.Contains(m.Content, "!hi") {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%+v Hello!", m.Author.Mention()))
	}
}
