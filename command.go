package bocto

import "github.com/bwmarrin/discordgo"

type command struct {
	Function func(Bot, *discordgo.MessageCreate, []string)
	Prefix   bool
}
