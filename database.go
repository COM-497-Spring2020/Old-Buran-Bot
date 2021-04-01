package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func openDB() *sql.DB {
	LogMsg("Opening database")
	data, err := sql.Open("mysql", fmt.Sprintf("%+v:%+v@%+v",
		config.DBUsername, config.DBPassword, config.DatabaseInfo))
	if err != nil {
		fmt.Printf("%+v:", err)
	}
	LogMsg("Database opened.")
	return data
}

func (s *ScoreRow) Insert() {
	db := openDB()
	defer db.Close()
	rType, rImage := boolToInt(s)
	LogMsg("Inserting %+v, rType: %+v, rImage: %+v", s, rType, rImage)
	insert, err := db.Query(fmt.Sprintf("INSERT INTO scores VALUES ('%+v',%+v,%+v,%v,'%+v')",
		s.DiscordID, rType, s.RatingScore, rImage, time.Now()))
	if err != nil {
		fmt.Printf("%+v", err)
	}
	LogMsg("Insert complete. Closing!")
	insert.Close()
}

func (s *ScoreRow) Retrieve() {
	db := openDB()
	defer db.Close()
	rType, _ := boolToInt(s)
	LogMsg("Trying to retrieve %+v\n", s)
	LogMsg("rType: %+v\nDiscordID: %+v\n", rType, s.DiscordID)
	query := "SELECT DiscordID, RatingType, RatingScore, RatingImage, TimeStamp from scores where DiscordID=? and RatingType=?"
	err := db.QueryRow(query, s.DiscordID, rType).Scan(&s.DiscordID,
		&s.RatingType, &s.RatingScore, &s.RatingImage, &s.TimeStamp)
	if err != nil {
		LogMsg("%+v", err)
	}
	LogMsg("Returning %+v", s)
}
func (s *ScoreRow) Update() {
	db := openDB()
	defer db.Close()
	rType, rImage := boolToInt(s)
	LogMsg("Updating %+v, rType: %+v, rImage: %+v", s, rType, rImage)
	query := fmt.Sprintf("UPDATE scores set RatingScore=%+v, RatingImage=%+v, TimeStamp='%+v' where DiscordID='%+v' and RatingType=%+v",
		s.RatingScore, rImage, time.Now(), s.DiscordID, rType)
	update, err := db.Query(query)
	if err != nil {
		LogMsg("%+v", err)
	}
	LogMsg("Closing update.")
	update.Close()
}

func boolToInt(s *ScoreRow) (int, int) {
	rType := 0
	rImage := 0
	if s.RatingType {
		rType = 1
	}
	if s.RatingImage {
		rImage = 1
	}
	return rType, rImage
}
