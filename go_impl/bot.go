package main

import (
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// for reference
const (
	NotificationsChannelId = 829129887905742858
	TestingChannelId       = 1180294950403969025
	CommandPrefix          = "!"
)

type NotificationType int

const (
	All = iota
	Offers
	Marathon
)

type Offer struct {
	name       string
	priceData  []string
	isMarathon bool
	imageUrl   string
}

type LbwBotData struct {
	isMarathon       bool
	notificationRule NotificationType
	lastOffer        *Offer
}

func (bot *LbwBotData) updateOffer(s *discordgo.Session, m *discordgo.MessageCreate) {
	// is_new_offer, is_new_marathon = update_stored_offer_info()

	// if is_new_marathon:
	//     await context.send('!!!!!!!!!!!!!!!!! Marathon has started !!!!!!!!!!!!!!!!!')

	// await context.send(embed=create_notification_embed(last_offer, is_new_offer))

}

// declare bot functions
func handleAllMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages sent by the bot itself
	if m.Author.ID == s.State.User.ID {
		log.Println("Ignoring message from myself")
		return
	}

	// Check if the message starts with the prefix
	if !strings.HasPrefix(m.Content, CommandPrefix) {
		return
	}

	// Remove the prefix from the message content
	content := strings.TrimPrefix(m.Content, CommandPrefix)

	// Handle the command
	switch content {
	case "update":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "is-marathon":
		s.ChannelMessageSend(m.ChannelID, "Hello, World!")
	case "get-notification-setting":
		s.ChannelMessageSend(m.ChannelID, "Hello, World!")
		// _, err := session.ChannelMessageSendEmbed(channelID, embed)
	case "start":
		s.ChannelMessageSend(m.ChannelID, "Hello, World!")
	case "stop":
		s.ChannelMessageSend(m.ChannelID, "Hello, World!")
	case "set-interval":
		s.ChannelMessageSend(m.ChannelID, "Hello, World!")
	case "set-notification-setting":
		s.ChannelMessageSend(m.ChannelID, "Hello, World!")
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command")
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

	// Register the message handler
	bot.AddHandler(handleAllMessages)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Gracefully shutting down.")
}
