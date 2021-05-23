package bocto

import "github.com/bwmarrin/discordgo"

// command is a wrapper for a Function, and it specifies if its requires a prefix to used.
// this is depreciated as of 5/23/2021, as discord has added slash commands.
// will be supported for at least ? months while existing servers add the slash
// command privlege.
type command struct {
	Function func(Bot, *discordgo.MessageCreate, []string)
	Prefix   bool
}

// response is a message thats activated by a message without a prefix.
type response func(Bot, *discordgo.MessageCreate, []string)
