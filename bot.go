package bocto

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Bot is a representation of a discord chatbot.
type Bot struct {
	Color          int
	Confused       string
	DisablePhrases bool
	Mentioned      string
	Name           string
	Prefix         string
	Self           *discordgo.User
	Session        *discordgo.Session
	commands       map[string]func(
		Bot,
		*discordgo.MessageCreate,
		[]string)
	phrases map[string]string
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
// - Mentioned is what your bot says when Mentioned. Use this to inform people
// of your bots command prefix.
//
// - Confused is what your bot says when it recieves an invalid command.
//
// - Color designates what color an embed should have by default.
func (b *Bot) New(name, prefix, token, men, confused string, color int) error {
	var err error
	b.Name = name
	b.Prefix = prefix
	b.Color = color
	b.Mentioned = men
	b.Confused = confused
	b.commands = make(map[string]func(Bot, *discordgo.MessageCreate, []string))
	b.phrases = make(map[string]string)
	b.Session, err = discordgo.New("Bot " + token)
	if err != nil {
		return err
	}
	// phrases must be manually disabled
	b.DisablePhrases = false
	b.Self, err = b.Session.User("@me")
	return err
}

// String is a debug command.
func (b Bot) String() string {
	return fmt.Sprintf("\tBot Name: %v\n\tBot Prefix: %v\n\tBot Color: %v"+
		"\n\tBot Commands: %v\n\tBot Phrases: %v\n\tBot Session: %v\n",
		b.Name, b.Prefix, b.Color, b.commands, b.phrases, b.Session)
}

// AddPhrase adds a quirky phrase for our bot to respond to.
// these are implicit
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
// this is the heart and soul of your discord bot, MessageCreate is run
// every time your bot can read a message in a channel. I
func (b Bot) MessageCreate(session *discordgo.Session,
	message *discordgo.MessageCreate) {

	// ignore bot users
	if message.Author.Bot == true {
		return
	}

	input := sliceStrings(message.Message.Content)
	id := message.ChannelID

	// old prefix behavior check
	// i'm going to hard code this, future!me please remove this
	// around June 15th, 2019.
	if input[0] == "jade:" && len(input) >= 2 {
		session.ChannelMessageSend(id,
			"you should try "+b.Prefix+input[1]+" instead! :B")
		return
	}
	if input[0] == "vriska:" && len(input) >= 2 {
		session.ChannelMessageSend(id,
			"Hey num8skull I only respond to "+b.Prefix+input[1]+" now!")
		return
	}

	// command check
	if beginsWith(input[0], b.Prefix) {
		confused := true
		for key := range b.commands {
			if key == input[1] {
				go b.commands[key](b, message, input[1:])
				confused = false
				break
			}
		}
		if confused {
			session.ChannelMessageSend(id, b.Confused)
		}
		return
	}

	// phrase check
	go func() {
		for key, value := range b.phrases {
			if strings.Contains(message.Content, key) && b.notDisabled(message.GuildID) {
				session.ChannelMessageSend(id, value)
			}
		}
	}()

	// mention check
	if isMentioned(message.Message.Mentions, b.Self) {
		session.ChannelMessageSend(id, b.Mentioned)
	}
}

// ReadyEvent occurs when the bot recieves a ready event.
func (b *Bot) ReadyEvent(session *discordgo.Session,
	rdy *discordgo.Ready) {

	b.Session.UpdateStatus(0, "Prefix: \""+b.Prefix+"\"")

	fmt.Printf("Ready event recieved. %v online.\nGuilds: %v\n",
		b.Self,
		len(rdy.Guilds))

}

func beginsWith(s string, substr string) bool {
	if strings.Index(s, substr) == 0 {
		return true
	}
	return false
}

func isMentioned(users [](*discordgo.User), self *discordgo.User) bool {
	for _, ele := range users {
		if ele.String() == self.String() {
			return true
		}
	}
	return false
}

// notDisabled checks is a phrases are not disabled
func (b *Bot) notDisabled(guild string) bool {
	// dummy proofing
	if !b.DisablePhrases || len(b.phrases) > 0 {
		return true
	}

	return false
}

func sliceStrings(s string) []string {

	s = strings.ToLower(s)
	return strings.Split(s, " ")
}
