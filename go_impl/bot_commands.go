package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getLatestOffer() Offer {
	return Offer{
		isMarathon: false,
		name:       "name",
		imageUrl:   "image",
		priceData:  []string{},
	}
}

func (bot *LbwBotData) isNewOffer(offer Offer) bool {
	if bot.lastOffer == nil {
		return true
	}

	return bot.lastOffer.name == offer.name
}

func (bot *LbwBotData) UpdateOffer() (bool, bool) {
	currentOffer := getLatestOffer()
	isNewOffer := bot.isNewOffer(currentOffer)
	marathonChange := bot.isMarathon != currentOffer.isMarathon

	bot.lastOffer = &currentOffer

	return isNewOffer, marathonChange && currentOffer.isMarathon
}

func (bot *LbwBotData) CreateNotificationEmbed(isNewOffer bool) discordgo.MessageEmbed {
	title := ""
	if isNewOffer {
		title = "New"
	} else {
		title = "Current"
	}

	tokenizedName := strings.Replace(bot.lastOffer.name, " ", "+", -1)
	urlTokenizedName := strings.Replace(bot.lastOffer.name, " ", "%20", -1)

	googleSearchLink := fmt.Sprintf("https://www.google.com/search?q=%s&oq=%s&aqs=chrome..69i57j69i61.1483j0j4&sourceid=chrome&ie=UTF-8", tokenizedName, tokenizedName)
	vivinoSearchLink := fmt.Sprintf("https://www.vivino.com/search/wines?q=%s", tokenizedName)
	wineSearcerSearchLink := fmt.Sprintf("https://www.wine-searcher.com/find/%s", tokenizedName)
	cellarTrackerSearchLink := fmt.Sprintf("https://www.cellartracker.com/list.asp?fInStock=0&Table=List&iUserOverride=0&szSearch=%s", tokenizedName)

	binnysSearchLink := fmt.Sprintf("https://www.binnys.com/search/?query=%s", urlTokenizedName)
	vinSearchLink := fmt.Sprintf("https://barrington.vinchicago.com/websearch_results.html?kw=%s", tokenizedName)

	return discordgo.MessageEmbed{
		Title: fmt.Sprintf(":wine_glass:%s Offer:wine_glass:", title),
		Color: 0xf03d44,
		URL:   "https://www.lastbottlewines.com",
		Image: &discordgo.MessageEmbedImage{
			URL: bot.lastOffer.imageUrl,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://www.lastbottlewines.com/favicon.png",
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  bot.lastOffer.name,
				Inline: false,
			},
			{
				Name:   bot.lastOffer.priceData[0],
				Value:  bot.lastOffer.priceData[1],
				Inline: true,
			},
			{
				Name: "Search Links",
				Value: fmt.Sprintf("[Google](%s)\n[Vivino](%s)\n[Wine Searcher](%s)\n[Cellar Tracker](%s)",
					googleSearchLink, vivinoSearchLink, wineSearcerSearchLink, cellarTrackerSearchLink),
				Inline: false,
			},
			{
				Name:   "Shop Search Links",
				Value:  fmt.Sprintf("[Binny's](%s)\n[Vin](%s)", binnysSearchLink, vinSearchLink),
				Inline: false,
			},
		},
	}
}
