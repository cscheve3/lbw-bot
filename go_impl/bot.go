package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type NotificationType int

const (
	All = iota
	Offers
	Marathon
)

// for reference
const (
	NotificationsChannelId = 829129887905742858
	TestingChannelId       = 1180294950403969025
	commandPrefix          = "!"
)

type Offer struct {
	name  string
	price int
}

type LbwBot struct {
	// bot interface{}
	// TODO
	isMarathon       bool
	notificationRule NotificationType
	lastOffer        Offer
}

// declare bot functions
func handleAllMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		log.Println("Ignoring message from myself")
		return
	}

	// Check if the message starts with the prefix
	if strings.HasPrefix(m.Content, commandPrefix) {
		// Remove the prefix from the message content
		content := strings.TrimPrefix(m.Content, commandPrefix)

		// Handle the command
		switch content {
		case "ping":
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		case "hello":
			s.ChannelMessageSend(m.ChannelID, "Hello, World!")
		default:
			s.ChannelMessageSend(m.ChannelID, "Unknown command")
		}
	}
}

func main() {
	// for local, every time do `export DISCORD_BOT_TOKEN=<value from token.txt>`
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	// notificationChannelId := TestingChannelId

	// init bot
	bot, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	defer bot.Close()
	bot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = bot.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot ", s.State.User.Username, "as ", s.State.User.Discriminator, "is connected to discord")
	})

	// Register the messageCreate handler
	bot.AddHandler(handleAllMessages)
	// bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 	if strings.Contains(m.Content, "ping") {
	// 		if ch, err := s.State.Channel(m.ChannelID); err != nil || !ch.IsThread() {
	// 			thread, err := s.MessageThreadStartComplex(m.ChannelID, m.ID, &discordgo.ThreadStart{
	// 				Name:                "Pong game with " + m.Author.Username,
	// 				AutoArchiveDuration: 60,
	// 				Invitable:           false,
	// 				RateLimitPerUser:    10,
	// 			})
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 			_, _ = s.ChannelMessageSend(thread.ID, "pong")
	// 			m.ChannelID = thread.ID
	// 		} else {
	// 			_, _ = s.ChannelMessageSendReply(m.ChannelID, "pong", m.Reference())
	// 		}
	// 		games[m.ChannelID] = time.Now()
	// 		<-time.After(timeout)
	// 		if time.Since(games[m.ChannelID]) >= timeout {
	// 			archived := true
	// 			locked := true
	// 			_, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
	// 				Archived: &archived,
	// 				Locked:   &locked,
	// 			})
	// 			if err != nil {
	// 				panic(err)
	// 			}
	// 		}
	// 	}
	// })

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
