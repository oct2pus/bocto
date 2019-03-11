package embed

import "github.com/bwmarrin/discordgo"

// CreditsEmbed is used for accreditation of all people involved in making
// one of my bots.
func CreditsEmbed(botName, artist, artist2, artist3 string, url string,
	color int) *discordgo.MessageEmbed {

	thanks := "Avatar by " + artist + "\n"

	// this is such a hack
	if artist2 != "" {
		thanks += "Original Design by " + artist2 + "\n"
	}

	thanks += "Emojis by " + artist3

	embed := &discordgo.MessageEmbed{
		Color: color,
		Type:  "About",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name: botName,
				Value: "Created by \\🐙\\🐙#0413" +
					" ( http://oct2pus.tumblr.com/ )\n" +
					botName + " uses the 'discordgo' library\n" +
					"( https://github.com/bwmarrin/discordgo/ )",
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name:   "Special Thanks",
				Value:  thanks,
				Inline: false,
			},
			&discordgo.MessageEmbedField{
				Name: "Disclaimer",
				Value: botName + " uses **Mutant Standard Emoji**" +
					" ( https://mutant.tech )\n**Mutant Standard Emoji** are " +
					" licensed under CC-BY-NC-SA 4.0\n" +
					"( https://creativecommons.org/licenses/by-nc-sa/4.0/ )",
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: url,
		},
	}

	return embed
}

// ImageEmbed returns an embeded image.
func ImageEmbed(title string, url string, image string,
	footer string, color int) *discordgo.MessageEmbed {

	if url == "" {
		url = image
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: color,
		URL:   url,
		Image: &discordgo.MessageEmbedImage{
			URL: image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
	}
	return embed
}

// SendMessage sends a message to a discord channel.
func SendMessage(s *discordgo.Session,
	channelID, message string) {
	s.ChannelMessageSend(channelID, message)
}

// SendEmbededMessage sends an embed message to a discord channel.
func SendEmbededMessage(s *discordgo.Session,
	channelID string,
	embed *discordgo.MessageEmbed) {

	s.ChannelMessageSendEmbed(channelID, embed)
}
