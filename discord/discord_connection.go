package discord

import (
	config "discord_crypto/util"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	BotID        string
	goBotSession *discordgo.Session
)

func StartBot() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Config Error: ", err)
	}
	fmt.Println("Config OK!")

	goBotSession, err := discordgo.New("Bot " + config.DISCORD_TOKEN)

	if err != nil {
		log.Fatal("Discord Error: ", err)
	}
	fmt.Println("Discord OK!")

	user, err := goBotSession.User("@me")

	if err != nil {
		log.Fatal("Bot Error: ", err)
	}
	fmt.Println("Bot OK!")

	BotID = user.ID

	goBotSession.AddHandler(forward)
	goBotSession.AddHandler(genesis)
	// goBotSession.AddHandler(test)

	// Open connection
	err = goBotSession.Open()

	if err != nil {
		log.Fatal("Session Error: ", err)
	}
	fmt.Println("Session OK!")
	fmt.Println("Bot is running!")

	// Always listening
	<-make(chan struct{})

}
