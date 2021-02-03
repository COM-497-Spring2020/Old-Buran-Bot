# Old-Buran-Bot
Group 1's Discord Bot, Shane Irons and Jordan Mason.

## Proposal

This bot should allow users of the official Monster Sanctuary Discord to input and later call their high scores in both the PvP leaderboard and the Infinity Arena. A user will be able to ping the bot to input this information in either plaintext or an optional screenshot, preferably, as this has more verifiable information. For PvP, this will allow users to also keep track of their successful PvP teams, to contribute to discussion within the PvP section of the Discord. Upon inputting a new record, the previous will be overwritten. This bot will be hosted on a gentoo-based cloud instance ran through a hypervisor.

## Design Decisions

* Written in Go
* Using a SQL Database
* Allowing users to upload screenshots
* Using Gentoo to host the database
* Using portable code that can be hosted in cloud environments
* Allowing interaction with database via Discord
* Use [bwmarrin/discordgo](https://godoc.org/github.com/bwmarrin/discordgo) as a Discord library
* Use [go-sql-driver/mysql](https://pkg.go.dev/github.com/go-sql-driver/mysql) as the MySQL Driver for the database
* Web interface using Discord OAuth to interact with Database (if there is time at the end)


## Database Schema
```
+-------------+------------+------+-----+---------------------+-------------------------------+
| Field       | Type       | Null | Key | Default             | Extra                         |
+-------------+------------+------+-----+---------------------+-------------------------------+
| DiscordID   | bigint(18) | YES  |     | NULL                |                               |
| RatingType  | bit(1)     | YES  |     | NULL                |                               |
| RatingScore | int(3)     | YES  |     | NULL                |                               |
| RatingImage | bit(1)     | YES  |     | NULL                |                               |
| TimeStamp   | timestamp  | NO   |     | 0000-00-00 00:00:00 | on update current_timestamp() |
+-------------+------------+------+-----+---------------------+-------------------------------+
```
