package db

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

//go:embed migrations/*.sql
var fs embed.FS

type Database struct {
	*sql.DB
}

func New() *Database {
	sqliteDb, err := sql.Open("sqlite3", "database.sqlite")
	if err != nil {
		panic("Failed to open sqlite DB")
	}

	return &Database{sqliteDb}
}

// RunMigrateScripts
//
// TODO: Если Drive изменится, то изменить реализацию данной функции
func RunMigrateScripts(db *sql.DB) {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		panic(err)
	}

	d, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithInstance("iofs", d, "main", driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Println(err)
	}
}

func (s Database) UpdateUsers(users []map[string]interface{}) {
	for _, user := range users {
		var data int
		err := s.QueryRow("SELECT count(*) FROM users as u WHERE u.name = ? and u.surname = ? and u.patronymic = ?", user["name"], user["surname"], user["patronymic"]).Scan(&data)
		if err != nil {
			log.Println(err)
		}

		if data == 0 {
			_, err := s.Exec("INSERT INTO users (name, surname, patronymic, telegram, birthday) VALUES (?, ?, ?, ?, ?)", user["name"], user["surname"], user["patronymic"], user["telegram"], user["birthday"])
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err := s.Exec("UPDATE users SET telegram = ?, birthday = ? WHERE name = ? and surname = ? and patronymic = ?", user["telegram"], user["birthday"], user["name"], user["surname"], user["patronymic"])
			if err != nil {
				log.Println(err)
			}
		}
	}
}
