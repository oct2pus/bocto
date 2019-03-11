package bot

import (
	"fmt"
	"strings"

	"github.com/oct2pus/bot/embed"

	"github.com/bwmarrin/discordgo"
)

// Bot is a representation of our bot.
type Bot struct {
	Name      string
	Self      string
	Prefix    string
	Session   *discordgo.Session
	Mentioned string
	Confused  string
	Color     int
	commands  map[string]func(
		Bot,
		*discordgo.MessageCreate,
		[]string)
	phrases map[string]string
}

// New initializes a new Bot.
func (b *Bot) New(name, prefix, token, men, conf string, color int) error {
	var err error
	b.Name = name
	b.Prefix = prefix
	b.Color = color
	b.Mentioned = men
	b.Confused = conf
	b.commands = make(map[string]func(Bot, *discordgo.MessageCreate, []string))
	b.phrases = make(map[string]string)
	b.Session, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	b.Self = b.Session.User("@me")
	return nil
}

func (b Bot) String() string {
	return fmt.Sprintf("\tBot Name: %v\n\tBot Prefix: %v\n\tBot Color: %v"+
		"\n\tBot Commands: %v\n\tBot Phrases: %v\n\tBot Session: %v\n",
		b.Name, b.Prefix, b.Color, b.commands, b.phrases, b.Session)
}

// AddPhrase adds a quirky phrase for our bot to respond to.
func (b *Bot) AddPhrase(key, value string) {
	b.phrases[key] = value
}

// AddCommand adds a Command to a Bot.
func (b *Bot) AddCommand(key string,
	value func(
		Bot,
		*discordgo.MessageCreate,
		[]string)) {

	b.commands[key] = value
}

// MessageCreate occurs every time the bot recieves a message.
func (b Bot) MessageCreate(session *discordgo.Session,
	message *discordgo.MessageCreate) {

	// bot user gotcha; ignore bots
	if message.Author.Bot == true {
		return
	}

	input := sliceStrings(message.Message.Content)
	// split messages into chunks
	go func() {
		confused := true
		if input[0] == b.Prefix && len(input) >= 2 {
			for key := range b.commands {
				if key == input[1] {
					go b.commands[key](b, message, input[2:])
					confused = false
					break
				}
			}
			if confused {
				go embed.SendMessage(b.Session, message.ChannelID, b.Confused)
			}
		}
	}()

	go func() {
		for key, value := range b.phrases {
			if strings.Contains(message.Content, key) {
				go b.Session.ChannelMessageSend(message.ChannelID, value)
				break
			}
		}
	}()

	// check if mentioned
	if isMentioned(message.Message.Mentions, b.Self) {
		go b.Session.ChannelMessageSend(message.ChannelID,
			b.Mentioned)
	}
}

// ReadyEvent occurs when the bot recieves a ready event.
func (b *Bot) ReadyEvent(session *discordgo.Session,
	rdy *discordgo.Ready) {

	b.Session.UpdateStatus(0, "prefix: '"+b.Prefix+"'")

	fmt.Printf("Ready event recieved. %v online.\nGuilds: %v\n",
		b.Self,
		len(rdy.Guilds))

}

func sliceStrings(m string) []string {

	m = strings.ToLower(m)
	return strings.Split(m, " ")
}

func isMentioned(users [](*discordgo.User), self string) bool {
	for _, ele := range users {
		if ele.String() == self {
			return true
		}
	}
	return false
}
