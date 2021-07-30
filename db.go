package main

import (
	"database/sql"
	"fmt"
	// "log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type albumDB struct {
	sql *sql.DB
}

var db albumDB

// func main() {
// 	if err := initDB(); err != nil {
// 		log.Fatal(err)
// 	}
//
// 	alb, err := db.albumByID(2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Album found: %v\n", alb)
//
// 	albID, err := db.addAlbum(Album{
// 		Title:  "The Modern Sound of Betty Carter",
// 		Artist: "Betty Carter",
// 		Price:  49.99,
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("ID of added album: %v\n", albID)
// }

// initDB initializes the recordings database.
func initDB() error {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}
	// Get a database handle.
	var err error
	sqlDB, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return fmt.Errorf("initDB %d: unable to initialize database", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("initDB %d: unable to verify connection", err)
	}
	fmt.Println("Connected!")
	db.sql = sqlDB
	return nil
}

// albumByID queries for the album with the specified ID.
func (db *albumDB) albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.sql.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func (db *albumDB) addAlbum(alb Album) (int64, error) {
	result, err := db.sql.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
