package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func openDB() *sql.DB {
	fmt.Println("Opening database")
	data, err := sql.Open("mysql", fmt.Sprintf("%+v:%+v@%+v",
		config.DBUsername, config.DBPassword, config.DatabaseInfo))
	if err != nil {
		fmt.Printf("%+V", err)
	}
	return data
}

func (s *ScoreRow) Insert() {
	db := openDB()
	defer db.Close()
	rType, rImage := boolToInt(s)
	insert, err := db.Query(fmt.Sprintf("INSERT INTO scores VALUES ('%+v',%+v,%+v,%v,'%+v')",
		s.DiscordID, rType, s.RatingScore, rImage, time.Now()))
	if err != nil {
		fmt.Printf("%+v", err)
	}
	insert.Close()
}

func (s *ScoreRow) Retrieve() {
	db := openDB()
	defer db.Close()
	query := "SELECT DiscordID, RatingType, RatingScore, RatingImage, TimeStamp from scores where DiscordID=?"

	err := db.QueryRow(query, s.DiscordID).Scan(&s.DiscordID,
		&s.RatingImage, &s.RatingScore, &s.RatingImage, &s.TimeStamp)
	if err != nil {
		fmt.Printf("%+v", err)
	}
}
func (s *ScoreRow) Update() {
	db := openDB()
	defer db.Close()
	rType, rImage := boolToInt(s)
	query := fmt.Sprintf("UPDATE scores set RatingType=%+v, RatingScore=%+v, RatingImage=%+v, TimeStamp='%+v' where DiscordID='%+v'",
		rType, s.RatingScore, rImage, time.Now(), s.DiscordID)
	update, err := db.Query(query)
	if err != nil {
		fmt.Printf("%+v", err)
	}
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
