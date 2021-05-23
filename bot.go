package bocto

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Bot is a representation of a discord chatbot.
type Bot struct {
	Color     int
	Name      string
	Prefix    string
	Mentioned func(Bot, *discordgo.MessageCreate, []string)
	Confused  func(Bot, *discordgo.MessageCreate, []string)
	Self      *discordgo.User
	Session   *discordgo.Session
	commands  map[string]command
}

// New initializes a new Bot.
//
// - Name is the name of your bot.
//
// - Prefix is what must be prepended to a command to inform your bot that
// this message is a command.
//
// - Token is your bot token.
//
// - Color designates what color an embed should have by default.
func (b *Bot) New(name, prefix, token string, color int) error {
	var err error
	b.Name = name
	b.Prefix = prefix
	b.Color = color
	b.commands = make(map[string]command)
	b.Session, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	b.Self, err = b.Session.User("@me")
	b.Mentioned = func(Bot, *discordgo.MessageCreate, []string) {}
	b.Confused = func(Bot, *discordgo.MessageCreate, []string) {}
	return err
}

// String is a debug command.
func (b Bot) String() string {
	return fmt.Sprintf("\tBot Name: %v\n\tBot Prefix: %v\n\tBot Color: %v"+
		"\n\tBot Commands: %v\n\tBot Session: %v\n",
		b.Name, b.Prefix, b.Color, b.commands, b.Session)
}

// AddCommand adds a Command to a Bot.
func (b *Bot) AddCommand(key string,
	value func(
		Bot,
		*discordgo.MessageCreate,
		[]string),
	prefix bool) {

	if prefix {
		key = b.Prefix + key
	}

	b.commands[key] = command{Function: value, Prefix: prefix}
}

// MessageCreate occurs every time the bot recieves a message.
// this is the heart and soul of your discord bot, MessageCreate is run
// every time your bot can read a message in a channel. I
func (b Bot) MessageCreate(session *discordgo.Session,
	message *discordgo.MessageCreate) {

	// ignore bot users
	if message.Author.Bot == true {
		return
	}

	input := sliceStrings(message.Message.Content)

	// command check
	if strings.HasPrefix(input[0], b.Prefix) {
		confused := true
		for key := range b.commands {
			if key == input[0] {
				go b.commands[key].Function(b, message, input[1:])
				confused = false
				break
			}
		}
		if confused {
			b.Confused(b, message, input)
		}
		return
	}

	// command check, no prefix
	for key, value := range b.commands {
		if !value.Prefix && strings.Contains(message.Content, key) {
			go b.commands[key].Function(b, message, input)
			return
		}
	}

	// mention check
	if isMentioned(message.Message.Mentions, b.Self) {
		b.Mentioned(b, message, input)
	}
}

// ReadyEvent occurs when the bot recieves a ready event.
func (b *Bot) ReadyEvent(session *discordgo.Session,
	rdy *discordgo.Ready) {

	// TODO: This no longer seems to exist as in discordgo v0.23.2
	//	b.Session.UpdateStatus(0, "Prefix: \""+b.Prefix+"\"")

	fmt.Printf("Ready event recieved. %v online.\nGuilds: %v\n",
		b.Self,
		len(rdy.Guilds))

}

func isMentioned(users [](*discordgo.User), self *discordgo.User) bool {
	for _, ele := range users {
		if ele.String() == self.String() {
			return true
		}
	}
	return false
}

func sliceStrings(s string) []string {

	s = strings.ToLower(s)
	return strings.Split(s, " ")
}
