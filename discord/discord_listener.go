package discord

import (
	"discord_crypto/db"
	"discord_crypto/dto"
	config "discord_crypto/util"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// genesis
func genesis(s *discordgo.Session, m *discordgo.MessageCreate) {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Config Error: ", err)
	}

	if m.ChannelID != config.RECEIVER_CHANNEL {
		return
	}

	if m.Author.ID != config.CORNIX_ID {
		return
	}

	messageDTO := dto.MessageToDTO(m.Content)
	messageDTO.MsgId = m.Message.ID

	// Check have reply or not
	messageReference := m.MessageReference

	if messageReference != nil {
		replyID := m.MessageReference.MessageID
		messageDTO.InReply = replyID
	}

	if messageDTO.InReply == "-1" {
		fmt.Println("Cornix: Geneisis Captured")
		db.AddGenesis(messageDTO, config)
	} else {
		if messageDTO.Status == "edited" {
			db.AddGenesis(messageDTO, config)
		}
		fmt.Println("Cornix: Update Captured")
		db.AddFollowUp(messageDTO, config)
	}
}

// Forward
func forward(s *discordgo.Session, m *discordgo.MessageCreate) {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Config Error: ", err)
	}

	if m.ChannelID != config.SENDER_CHANNEL {
		return
	}

	if m.Author.ID == BotID {
		return
	}

	finalMessage := m.Content + fmt.Sprintf("\nBy=%s (%s)", m.Author.Username, m.Author.ID)

	_, _ = s.ChannelMessageSend(config.RECEIVER_CHANNEL, finalMessage)

	return

}

// Test
func test(s *discordgo.Session, m *discordgo.MessageCreate) {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Config Error: ", err)
	}

	if m.ChannelID != config.RECEIVER_CHANNEL {
		return
	}

	if m.Author.ID == BotID {
		return
	}

	author, authorName, replyID := m.Author.ID, m.Author.Username, ""

	// Check have reply or not
	messageReference := m.MessageReference

	if messageReference != nil {
		replyID = m.MessageReference.MessageID
	}

	fmt.Println("------------------------Test----------------------------")
	fmt.Printf("Author ID: %s\n", author)
	fmt.Printf("authorName ID: %s\n", authorName)
	fmt.Printf("Reply ID: %s\n", replyID)
	fmt.Println("--------------------------------------------------------")
	return

}
