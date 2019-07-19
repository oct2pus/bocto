package bocto

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Contributor represents metadata about an Contributor, used for accreditation.
type Contributor struct {
	Name    string
	Message string
	URL     string
	Type 	string
}

// CreditsEmbed is used for accreditation of all people involved in making
// the bot possible. or at least, most of them.
func CreditsEmbed(botName, url string,
	color int, usesMutantStandard bool, artists ...Contributor) *discordgo.MessageEmbed {

	embed := &discordgo.MessageEmbed{
		Color: color,
		Type:  "About",
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: url,
		},
	}
	if len(artists) > 0 {
		for _, artist := range artists {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   artist.Type,
				Value:  fmt.Sprintf(artist.Message, artist.Name, artist.URL),
				Inline: false,
			})
		}
	}
	// this is such a hack
	if usesMutantStandard {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Disclaimer",
			Value: botName + " uses **Mutant Standard Emoji**" +
				" ( https://mutant.tech )\n**Mutant Standard Emoji** are " +
				" licensed under CC-BY-NC-SA 4.0\n" +
				"( https://creativecommons.org/licenses/by-nc-sa/4.0/ )",
			Inline: false,
		})
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

// TextEmbed returns an embeded Text.
func TextEmbed(title, name, text, url,
	footer string, color int) *discordgo.MessageEmbed {

	if url == "" {
		url = ""
	}

	embed := &discordgo.MessageEmbed{
		Title: title,
		Color: color,
		URL:   url,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  name,
				Value: text,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
	}
	return embed
}
