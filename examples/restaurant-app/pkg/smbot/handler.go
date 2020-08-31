package smbot

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/itsmurugappan/gql-source/examples/restaurant-app/pkg/twitter"
	"log"
)

const (
	NewItem     = "We are very excited to let you know that we are introducing a new %s in the menu, '''%s'''."
	RemoveItem  = "We are sorry to let you know that we are removing '''%s''' from the menu"
	InfoChanges = "PLEASE NOTE!!!: There is an update on our location and hours of operation. \n Location : %s, Hours : %s"
)

type SmBot struct {
	TwitterClient *twitter.Client
}

func InitBot() *SmBot {
	return &SmBot{TwitterClient: twitter.InitTwitterClient()}
}

// CloudEventReceived handles the cloud event post
func (bot *SmBot) ProcessCE(ctx context.Context, event cloudevents.Event) {
	log.Printf("cloud event received %v\n", event)
	var data map[string]interface{}
	err := event.DataAs(&data)
	if err != nil {
		log.Printf("Error marshaling data %s\n", err)
		return
	}
	dataNode := (data["data"]).(map[string]interface{})
	log.Printf("data %v\n", dataNode)
	if dataNode["itemChanged"] != nil {
		bot.handleItemChanges(dataNode["itemChanged"])
		return
	}
	bot.handleInfoChanges(dataNode["infoChanged"])
	return
}

func (bot *SmBot) handleItemChanges(ic interface{}) {
	log.Println("handle item changes")
	data := ic.(map[string]interface{})
	action := (data["action"]).(string)
	var msg string
	switch action {
	case "Add":
		msg = fmt.Sprintf(NewItem, (data["itemType"]).(string), (data["name"]).(string))
	case "Remove":
		msg = fmt.Sprintf(RemoveItem, (data["name"]).(string))
	}
	bot.TwitterClient.Tweet(msg, action)
}

func (bot *SmBot) handleInfoChanges(ic interface{}) {
	log.Println("handle info changes")
	data := ic.(map[string]interface{})
	msg := fmt.Sprintf(InfoChanges, (data["address"]).(string), (data["hours"]).(string))
	bot.TwitterClient.Tweet(msg, "Info")
}
