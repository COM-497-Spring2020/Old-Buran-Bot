package main

// Config struct for Old Buran bot
type Config struct {
	// GuildID is the ID of the Discord Guild the bot will interact with
	GuildID string
	// DatabaseInfo is the info for the database the bot will connect to
	DatabaseInfo string
	// DBUsername is the username for connecting to the database
	DBUsername string
	// DBPassword is the password for connecting to the database
	DBPassword string
	// BotToken is the token of the Old Buran bot
	BotToken string
}
