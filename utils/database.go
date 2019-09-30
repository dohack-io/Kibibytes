package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

// structs

type User struct {
	Id         int64
	Username   string
	Location   string
	State      int
	Travelmode string
}

type Notify struct {
	Id            int
	UserId        int64
	Context       string
	Executiontime int64
}

type Weather struct {
	Id int64
}

// Init

func init() {
	db := openDatabase()

	statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY,
			username VARCHAR(50),
			location VARCHAR(50),
			travelmode VARCHAR(5),
			state int
		);
	`)
	if err != nil {
		log.Panic(err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Panic(err)
	}

	statement, err = db.Prepare(`
		CREATE TABLE IF NOT EXISTS notifies (
			id INTEGER PRIMARY KEY,
			userId INT,
			executiontime VARCHAR(255),
			context VARCHAR(50)
		);
	`)
	if err != nil {
		log.Panic(err)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Panic(err)
	}

	defer statement.Close()
	defer db.Close()
}

// Notify

func InsertNotify(notify Notify) {
	db := openDatabase()

	stmt, _ := db.Prepare("INSERT INTO notifies (userId, context, executiontime) VALUES (?, ?, ?)")
	_, execError := stmt.Exec(notify.UserId, notify.Context, notify.Executiontime)

	if execError != nil {
		log.Panic(execError)
	}

	err := stmt.Close()
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
}

func GetNextNotifies(id int64) []Notify {
	var notifies []Notify

	db := openDatabase()

	stmt, err := db.Prepare("SELECT id, userId, context, executiontime FROM notifies WHERE executiontime < ? ORDER BY executiontime DESC")
	if err != nil {
		log.Panic(err)

		return nil
	}

	rows, err := stmt.Query(GetUnixtimestamp(0) + 30)
	if err != nil {
		log.Panic("Notify not found")
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var notify Notify
		err = rows.Scan(&notify.Id, &notify.UserId, &notify.Context, &notify.Executiontime)
		if err != nil {
			log.Fatal(err)
		}

		notifies = append(notifies, notify)
	}

	defer stmt.Close()

	return notifies
}

func DeleteNotify(id int) {
	db := openDatabase()

	stmt, err := db.Prepare("DELETE FROM notifies WHERE id = ?")
	if err != nil {
		log.Panic(err)
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Panic("Notify not found")
	}

	defer stmt.Close()
	defer db.Close()
}

// User

func InsertUser(user User) {
	db := openDatabase()
	checkUser := getUser(user.Id, db)

	if checkUser.Id == 0 {
		stmt, _ := db.Prepare("INSERT INTO users (id, username, location, state, travelmode) VALUES (?, ?, ?, ?, ?)")

		_, execError := stmt.Exec(strconv.FormatInt(user.Id, 10), user.Username, user.Location, user.State, user.Travelmode)
		if execError != nil {
			log.Panic(execError)
		}

		err := stmt.Close()
		if err != nil {
			log.Panic(err)
		}
	}

	defer db.Close()
}

func UpdateUser(user User) bool {
	db := openDatabase()

	tx, err := db.Begin()
	if err != nil {
		log.Panic(err)
		return false
	}

	stmt, err := tx.Prepare("UPDATE users SET username = ?, location = ?, state = ?, travelmode = ? WHERE id = ?")
	if err != nil {
		log.Panic(err)
		return false
	}

	_, err = stmt.Exec(user.Username, user.Location, user.State, user.Travelmode, user.Id)
	if err != nil {
		log.Panic(err)
		return false
	}

	err = tx.Commit()
	if err != nil {
		log.Panic(err)
	}

	defer stmt.Close()
	defer db.Close()

	return true
}

func GetUser(id int64) User {
	db := openDatabase()
	user := getUser(id, db)

	defer db.Close()

	return user
}

// Private Functions

func openDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Panic(fmt.Sprintf("%s - %s", "openDatabase", err))
	}

	return db
}

// Same function to just open the database connection once
func getUser(id int64, db *sql.DB) User {
	var user User

	stmt, err := db.Prepare("SELECT id, username, location, travelmode, state FROM users WHERE id = ?")
	if err != nil {
		log.Panic(err)

		return User{}
	}

	row, err := stmt.Query(id)

	if err != nil {
		log.Panic("User not found")

		return User{}
	}

	row.Next()
	_ = row.Scan(&user.Id, &user.Username, &user.Location, &user.Travelmode, &user.State)

	defer row.Close()
	defer stmt.Close()

	return user
}
