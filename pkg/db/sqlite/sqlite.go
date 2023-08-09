package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	NickName  string `json:"nickName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Dob       string `json:"dob"`
	Image     string `json:"image"`
	About     string `json:"about"`
	Public    bool   `json:"public"`
}

type Post struct {
	ID        int    `json:"id"`
	Author    int    `json:"userId"`
	CreatedAt string `json:"createdAt"`
	Message   string `json:"message"`
	Image     string `json:"image"`
	Privacy   bool   `json:"privacy"`
}

var db *sql.DB

var m *migrate.Migrate

// run migrations

func RunMigration() *migrate.Migrate {
	file := "file://pkg/db/migration/sqlite"
	sqlite := "sqlite3://pkg/db/database.db"
	if runtime.GOOS == "darwin" {
		file = "file://../../pkg/db/migration/sqlite"
		sqlite = "sqlite3://../../pkg/db/database.db"
	}
	m, err := migrate.New(file, sqlite)

	if err != nil {
		fmt.Print(err.Error())
	}

	m.Up()

	return m
}

// remove migrations

func RemoveMigration(m *migrate.Migrate) {
	m.Migrate(16)
}

// connect to database

func DbConnect() *sql.DB {
	if runtime.GOOS == "darwin" {
		db, err := sql.Open("sqlite3", "../../pkg/db/database.db")
		if err != nil {
			log.Fatal(err)
		}
		return db
	} else {
		db, err := sql.Open("sqlite3", "pkg/db/database.db")
		if err != nil {
			log.Fatal(err)
		}
		return db
	}
}

// insert mock user data

func InsertMockUserData() {

	// fetch api for mock data

	var res *http.Response

	res, _ = http.Get("https://63f35a0e864fb1d60014de90.mockapi.io/users")

	resData, _ := ioutil.ReadAll(res.Body)

	// Unmarshall http response

	var responseObject []User

	json.Unmarshal(resData, &responseObject)

	// insert into database

	db := DbConnect()

	for _, user := range responseObject {

		stmt, err := db.Prepare("INSERT INTO user(first_name, last_name, nick_name, email, password_, dob, image_, about, public) VALUES(?,?,?,?,?,?,?,?,?);")
		if err != nil {
			log.Fatal(err)
		}

		defer stmt.Close()

		stmt.Exec(user.FirstName, user.LastName, user.NickName, user.Email, user.Password, user.Dob, user.Image, user.About, 1)
	}
}

// insert mock post data

func InsertMockPostData() {

	// fetch api for mock data

	// iteration for 50 users
	for i := 1; i < 51; i++ {
		var res *http.Response

		res, _ = http.Get("https://63f35a0e864fb1d60014de90.mockapi.io/users/" + strconv.Itoa(i) + "/posts")

		resData, _ := ioutil.ReadAll(res.Body)

		// Unmarshall http response

		var responseObject []Post

		json.Unmarshal(resData, &responseObject)

		// insert into database
		db := DbConnect()

		for _, post := range responseObject {

			stmt, err := db.Prepare("INSERT INTO post(author, message_, image_, created_at, privacy) VALUES(?,?,?,?,?);")
			if err != nil {
				log.Fatal(err)
			}

			defer stmt.Close()

			stmt.Exec(i, post.Message, post.Image, post.CreatedAt, 0)
		}
	}

}
