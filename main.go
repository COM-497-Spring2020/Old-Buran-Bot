package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	configFile string
	config     Config
	dg         *discordgo.Session
	debug      bool
)

func init() {
	flag.StringVar(&configFile, "c", "", "Config file")
	flag.BoolVar(&debug, "d", false, "Debug flag")
	flag.Parse()
}

func main() {
	LogMsg("Debug mode enabled.")
	if configFile == "" {
		fmt.Println("No config specified.")
		return
	} else {
		loadConfig()
	}
	LogMsg("Config: %+v\n", config)
	dg, _ = discordgo.New("Bot " + config.BotToken)
	dg.AddHandler(messageCreate)
	_ = dg.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	LogMsg("Detected incoming message.")
	// Return if message came from a bot, or doesn't mention this bot
	if m.Author.Bot || !strings.Contains(m.Content, s.State.User.ID) || !strings.Contains(config.GuildID, m.GuildID) {
		LogMsg("Ignoring message.")
		return
	}
	// Split input for use in command functions
	parts := strings.Split(m.Content, " ")
	b := BotCommand{
		Session:   s,
		Channel:   m.ChannelID,
		Message:   m,
		Command:   parts[1],
		DiscordID: m.Author.ID,
		Parts:     parts,
	}
	LogMsg("Command detected: %+v", b)
	if strings.Contains(b.Command, "iaadd") {
		IAadd(b)
		return
	}
	if strings.Contains(b.Command, "pvpadd") {
		PvPadd(b)
		return
	}
	if strings.Contains(b.Command, "iacheck") {
		IAcheck(b)
		return
	}
	if strings.Contains(b.Command, "pvpcheck") {
		PvPcheck(b)
		return
	}

	// No valid command found
	b.Reply(fmt.Sprintf("Unknown command: %+v. Ping me with **iaadd** followed by a number or a link to the in-game screenshot of your score to record your Infinity Arena high score. Ping me with **pvpadd** followed by a number or a link to the in-game screenshot of your leaderboard high score to record your PvP leaderboard score. Ping me with **iacheck** or **pvpcheck** to request the information you have inserted.", b.Command))

}

// IAadd will add an IA score to the database
func IAadd(b BotCommand) {
	// Return if there aren't enough parts
	if len(b.Parts) < 3 {
		b.Reply(fmt.Sprintf("Invalid use of iaadd command <@%+v>. Not enough arguments.", b.DiscordID))
		return
	}
	// For determining how to store in the database
	isImage := true
	// Score for echo if int or string
	score := 0
	// Check if message is a link to an image on discord's CDN
	if !strings.HasPrefix(b.Parts[2], "https://cdn.discordapp.com/attachments/") {
		// Message was not an image
		isImage = false
		// Try to convert to an int

		if !strings.HasPrefix(b.Parts[2], "https://media.discordapp.net/attachments/") {
			isImage = false

			if intScore, ok := strconv.Atoi(b.Parts[2]); ok != nil {
				// Could not convert to an int, invalid score!
				b.Reply("Invalid use of iaadd command. Not enough arguments.")
				return
			} else {
				score = intScore
			}
		}
	}
	// Save score for printing
	s := ScoreRow{
		DiscordID:   b.DiscordID,
		RatingType:  true,
		RatingImage: isImage,
		RatingScore: score,
	}

	x := ScoreRow{
		DiscordID:  b.DiscordID,
		RatingType: true,
	}

	if isImage {
		storeImage(b, s)
		b.Response = fmt.Sprintf("Updating your Infinity Arena High Score with the provided image.")
	}

	x.Retrieve()
	if x.TimeStamp != "" {
		b.Response = fmt.Sprintf("Updating Infinity Arena High Score with the provided information.")
		s.Update()
	} else {
		b.Response = fmt.Sprintf("Inserting Infinity Arena High Score.")
		s.Insert()

	}
	b.Reply("")

}

// PvPadd will add a PvP score to the database
func PvPadd(b BotCommand) {
	// Return if there aren't enough parts
	if len(b.Parts) < 3 {
		b.Reply("Invalid use of pvpadd command <@%+v>. Not enough arguments.")
		return
	}
	// For determining how to store in the database
	isImage := true
	// Score for echo if int or string
	score := 0
	// Check if message is a link to an image on Discord's CDN
	if !strings.HasPrefix(b.Parts[2], "https://cdn.discordapp.com/attachments/") {
		// Message was not an image
		isImage = false
		if !strings.HasPrefix(b.Parts[2], "https://media.discordap..net/attachments/") {
			isImage = false
			// Try to convert to an int
			if intScore, ok := strconv.Atoi(b.Parts[2]); ok != nil {
				// Could not convert to an int, invalid score!
				b.Reply("Invalid use of pvpadd command.")
				return
			} else {
				score = intScore
			}
		}
	}
	// Save score for printing

	s := ScoreRow{
		DiscordID:   b.DiscordID,
		RatingType:  false,
		RatingImage: isImage,
		RatingScore: score,
	}
	x := ScoreRow{
		DiscordID:  b.DiscordID,
		RatingType: false,
	}
	if isImage {
		storeImage(b, s)
		b.Response = fmt.Sprintf("Updating your PvP Leaderboard High Score with the provided image.")

	}

	x.Retrieve()
	if x.TimeStamp != "" {
		b.Response = fmt.Sprintf("Updating PvP Leaderboard High Score with the provided information.")
		s.Update()
	} else {
		b.Response = fmt.Sprintf("Inserting PvP Leaderboard High Score.")
		s.Insert()

	}
	b.Reply("")
}

func IAcheck(b BotCommand) {

	s := ScoreRow{
		DiscordID:  b.DiscordID,
		RatingType: true,
	}
	s.Retrieve()
	if s.RatingImage {
		b.Response = fmt.Sprintf("See image above.")
		upperDirPattern := fmt.Sprintf("./buran_users/%+v-%+v-*", b.DiscordID, s.RatingType)

		matches, err := filepath.Glob(upperDirPattern)

		if err != nil {
			fmt.Println(err)
		}
		if len(matches) != 1 {
			b.Session.ChannelMessageSend(b.Channel, fmt.Sprintf("Error with your thing"))
			return
		}

		fmt.Println(matches)

		RatingImage, err := os.Open(matches[0])
		if err != nil {
			LogMsg("Some log words here")
			return

		}
		msg := fmt.Sprintf("I have received your request.")
		b.Session.ChannelFileSendWithMessage(b.Channel, msg, matches[0], RatingImage)
	}
	b.Reply(fmt.Sprintf("Your Infinity Arena High Score is: %+v.", s.RatingScore))
}

func PvPcheck(b BotCommand) {
	s := ScoreRow{
		DiscordID:  b.DiscordID,
		RatingType: false,
	}
	s.Retrieve()
	if s.RatingImage {
		b.Response = fmt.Sprintf("See image above.")
		upperDirPattern := fmt.Sprintf("./buran_users/%+v-%+v-*", b.DiscordID, s.RatingType)

		matches, err := filepath.Glob(upperDirPattern)

		if err != nil {
			fmt.Println(err)
		}
		if len(matches) != 1 {
			b.Session.ChannelMessageSend(b.Channel, fmt.Sprintf("Error with your things"))
			return
		}

		fmt.Println(matches)

		RatingImage, err := os.Open(matches[0])
		if err != nil {
			LogMsg("Some log words here")
			return
		}
		msg := fmt.Sprintf("I have received your request.")
		b.Session.ChannelFileSendWithMessage(b.Channel, msg, matches[0], RatingImage)
	}
	b.Reply(fmt.Sprintf("Your PvP Leaderboard High Score is: %+v.", s.RatingScore))
}

// Reply will reply to the BotCommand.Message, tagging the sender. If b.Response is set, it will use that otherwise the string will be used
func (b BotCommand) Reply(s string) {
	if len(b.Response) > 0 {
		b.Session.ChannelMessageSend(b.Channel, fmt.Sprintf("<@%+v>: %+v", b.DiscordID, b.Response))
	} else {
		b.Session.ChannelMessageSend(b.Channel, fmt.Sprintf("<@%+v>: %+v", b.DiscordID, s))
	}
}

func storeImage(b BotCommand, s ScoreRow) {
	LogMsg("Input detected", b)
	if !strings.HasPrefix(b.Parts[2], "https://cdn.discordapp.com/attachments/") {
		return
	}

	LogMsg("Input detected", b)
	fileURL, err := url.Parse(b.Parts[2])
	if err != nil {
		LogMsg("Unable to parse attachment.")
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")

	filename := segments[len(segments)-1]
	upperDirPattern := fmt.Sprintf("./buran_users/%+v-%+v-*", b.DiscordID, s.RatingType)
	matches, err := filepath.Glob(upperDirPattern)

	if err != nil {
		fmt.Println(err)
	}

	if len(matches) != 0 {
		os.Remove(matches[0])
	}

	LogMsg("Input detected", b)
	file, err := os.Create(fmt.Sprintf("./buran_users/%+v-%+v-%+v", s.DiscordID, s.RatingType, filename))
	if err != nil {
		LogMsg("Unable to create file.")
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	resp, err := client.Get(b.Parts[2])
	if err != nil {
		LogMsg("Unable to download verification %+v-%+v-%+v, s.DiscordID, s.RatingType, filename")
	}
	defer resp.Body.Close()
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		LogMsg("Unable to store verification %+v-%+v", b.DiscordID, filename)
	}
}
