package main

import (
	"birthday-bot/client"
	"birthday-bot/db"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	inst := db.New()
	defer inst.Close()
	db.RunMigrateScripts(inst.DB)
}

func main() {
	c := client.New()
	c.InitData()

	inst := db.New()
	defer inst.Close()
	
	inst.UpdateUsers(c.Users)
}
