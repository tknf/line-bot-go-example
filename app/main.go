package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v7/linebot/httphandler"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()

		if err != nil {
			log.Print(err)
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					switch message.Text {
					case "Echo":
						// Echo bot
						if _, err := bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage(message.Text),
						).Do(); err != nil {
							log.Print(err)
							return
						}
					case "Flex":
						// Flex message
						container := &linebot.BubbleContainer{
							Type: linebot.FlexContainerTypeBubble,
							Body: &linebot.BoxComponent{
								Type:   linebot.FlexComponentTypeBox,
								Layout: linebot.FlexBoxLayoutTypeHorizontal,
								Contents: []linebot.FlexComponent{
									&linebot.TextComponent{
										Type: linebot.FlexComponentTypeText,
										Text: "Hello",
									},
									&linebot.TextComponent{
										Type: linebot.FlexComponentTypeText,
										Text: "World!",
									},
								},
							},
						}
						if _, err = bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewFlexMessage("Flex Message", container),
						).Do(); err != nil {
							log.Print(err)
							return
						}
					}
				}
			}
		}
	})

	http.Handle("/callback", handler)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
