package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token      string
	configFile string
	dg         *discordgo.Session
)

func init() {
	// To be removed after config is implemented
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

	// Return if message came from a bot, or doesn't mention this bot
	if m.Author.Bot || !strings.Contains(m.Content, s.State.User.ID) {
		return
	}
	// Split input for use in command functions
	parts := strings.Split(m.Content, " ")
	b := BotCommand{
		Session:   s,
		Message:   m,
		Command:   parts[1],
		DiscordID: m.Author.ID,
		Parts:     parts,
	}
	if strings.Contains(b.Command, "iaadd") {
		// Reassign parts[0] from Bot name to message Channel ID
		// Reassign parts[1] from command (since we already know it) to Author ID (for mentions)
		parts[0] = m.ChannelID
		parts[1] = m.Author.ID
		IAadd(s, parts)
		return
	}
	if strings.Contains(b.Command, "pvpadd") {
		parts[0] = m.ChannelID
		parts[1] = m.Author.ID
		PvPadd(s, parts)
		return
	}
	if strings.Contains(b.Command, "iacheck") {
		parts[0] = m.ChannelID
		parts[1] = m.Author.ID
		IAcheck(s, parts)
		return
	}
	if strings.Contains(parts[1], "pvpcheck") {
		PvPcheck(b)
		return
	}

	// No valid command found
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown command: %+v", parts[1]))

}

// IAadd will add an IA score to the database
func IAadd(s *discordgo.Session, cmd []string) {
	// Return if there aren't enough parts
	if len(cmd) < 3 {
		s.ChannelMessageSend(cmd[0], fmt.Sprintf("Invalid use of iaadd command <@%+v>. Not enough arguments.", cmd[1]))
		return
	}
	// For determining how to store in the database
	isImage := true
	// Score for echo if int or string
	score := "0"
	// Check if message is a link to an image on discord's CDN
	if !strings.HasPrefix(cmd[2], "https://cdn.discordapp.com/attachments/") {
		// Message was not an image
		isImage = false
		// Try to convert to an int
		if _, ok := strconv.Atoi(cmd[2]); ok != nil {
			// Could not convert to an int, invalid score!
			s.ChannelMessageSend(cmd[0], fmt.Sprintf("Invalid use of iaadd command <@%+v>.", cmd[1]))
			return
		}
	}
	// Save score for printing
	score = cmd[2]
	if isImage {
		s.ChannelMessageSend(cmd[0], fmt.Sprintf("Image score detected: %+v", score))
	} else {
		s.ChannelMessageSend(cmd[0], fmt.Sprintf("Integer score detected: %+v", score))
	}

}

// PvPadd will add a PvP score to the database
func PvPadd(s *discordgo.Session, cmd []string) {
	// Return if there aren't enough parts
	if len(cmd) < 3 {
		s.ChannelMessageSend(cmd[0], fmt.Sprintf("Invalid use of pvpadd command <@%+v>. Not enough arguments.", cmd[1]))
		return
	}
	// For determining how to store in the database
	isImage := true
	// Score for echo if int or string
	score := "0"
	// Check if message is a link to an image on Discord's CDN
	if !strings.HasPrefix(cmd[2], "https://cdn.discordapp.com/attachments/") {
		// Message was not an image
		isImage = false
		// Try to convert to an int
		if _, ok := strconv.Atoi(cmd[2]); ok != nil {
			// Could not convert to an int, invalid score!
			s.ChannelMessageSend(cmd[0], fmt.Sprintf("Invalid use of pvpadd command <@%+v>.", cmd[1]))
			return
		}
	}
	// Save score for printing
	score = cmd[2]
	if isImage {
		s.ChannelMessageSend(cmd[0], fmt.Sprintf("Image score detected: %+v", score))
	} else {
		s.ChannelMessageSend(cmd[0], fmt.Sprintf("Integer score detected: %+v", score))
	}
}

func IAcheck(s *discordgo.Session, cmd []string) {
	s.ChannelMessageSend(cmd[0], fmt.Sprintf("I have received your request."))
}

func PvPcheck(b BotCommand) {
	b.Reply("I have received your request.")
}

// Reply will reply to the BotCommand.Message, tagging the sender. If b.Response is set, it will use that otherwise the string will be used
func (b BotCommand) Reply(s string) {
	if len(b.Response) > 0 {
		b.Session.ChannelMessageSend(b.Channel, fmt.Sprintf("<@%+v>: %+v", b.DiscordID, b.Response))
	} else {
		b.Session.ChannelMessageSend(b.Channel, fmt.Sprintf("<@%+v>: %+v", b.DiscordID, s))
	}
}
