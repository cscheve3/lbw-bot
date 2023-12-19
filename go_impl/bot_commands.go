package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bwmarrin/discordgo"
)

func getLatestOffer() Offer {
	res, err := http.Get("https://www.lastbottlewines.com")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	isMarathon := doc.Find("div.marquee-top").Length() > 0 || doc.Find("div.marathon").Length() > 0
	offerName := doc.Find("h1.offer-name").Text()
	imageSrc, _ := doc.Find("img#offer-image").Attr("src")
	offerImage := ""
	if imageSrc != "" {
		offerImage = "http:" + imageSrc
	}

	priceDataMap := make(map[string]PriceData)
	doc.Find("div.price-holder").Each(func(i int, s *goquery.Selection) {
		label := s.Find("p").Text()
		amount := s.Find("span.amount").Text()
		priceDataMap[label] = PriceData{
			label,
			amount,
		}
	})

	priceData := make([]PriceData, 0)
	for _, data := range priceDataMap {
		priceData = append(priceData, data)
	}

	return Offer{
		isMarathon: isMarathon,
		name:       offerName,
		imageUrl:   offerImage,
		priceData:  priceData,
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

func (bot *LbwBotData) CreateNotificationEmbed(isNewOffer bool) *discordgo.MessageEmbed {
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

	priceEmbedFields := make([]*discordgo.MessageEmbedField, 0)
	for _, price := range bot.lastOffer.priceData {
		priceEmbedFields = append(priceEmbedFields, &discordgo.MessageEmbedField{
			Name:   price.label,
			Value:  price.amount,
			Inline: true,
		})
	}

	embedFields := []*discordgo.MessageEmbedField{
		{
			Name:   "Name",
			Value:  bot.lastOffer.name,
			Inline: false,
		},
	}

	embedFields = append(embedFields, priceEmbedFields...)

	embedFields = append(embedFields, &discordgo.MessageEmbedField{
		Name: "Search Links",
		Value: fmt.Sprintf("[Google](%s)\n[Vivino](%s)\n[Wine Searcher](%s)\n[Cellar Tracker](%s)",
			googleSearchLink, vivinoSearchLink, wineSearcerSearchLink, cellarTrackerSearchLink),
		Inline: false,
	}, &discordgo.MessageEmbedField{
		Name:   "Shop Search Links",
		Value:  fmt.Sprintf("[Binny's](%s)\n[Vin](%s)", binnysSearchLink, vinSearchLink),
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf(":wine_glass:%s Offer:wine_glass:", title),
		Color: 0xf03d44,
		URL:   "https://www.lastbottlewines.com",
		Image: &discordgo.MessageEmbedImage{
			URL: bot.lastOffer.imageUrl,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://www.lastbottlewines.com/favicon.png",
		},
		Fields: embedFields,
	}
}
