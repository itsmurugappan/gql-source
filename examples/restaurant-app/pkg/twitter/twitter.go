package twitter

import (
	"io/ioutil"
	"log"
	"os"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/kelseyhightower/envconfig"
)

type authConfig struct {
	ConsumerKey    string `envconfig:"CONSUMER_KEY" required:"true"`
	ConsumerSecret string `envconfig:"CONSUMER_SECRET" required:"true"`
	AccessToken    string `envconfig:"ACCESS_TOKEN" required:"true"`
	AccessSecret   string `envconfig:"ACCESS_SECRET" required:"true"`
}

type Client struct {
	client *twitter.Client
}

func InitTwitterClient() *Client {
	var auth authConfig
	//parse env variables
	if err := envconfig.Process("", &auth); err != nil {
		log.Panicf("[ERROR] Failed to process env var: %s", err)
	}
	config := oauth1.NewConfig(auth.ConsumerKey, auth.ConsumerSecret)
	token := oauth1.NewToken(auth.AccessToken, auth.AccessSecret)
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	return &Client{client: twitter.NewClient(httpClient)}
}

func (c *Client) Tweet(msg, action string) {
	tweetParams := &twitter.StatusUpdateParams{}
	filepath := fmt.Sprintf("%s/images/%s.gif", os.Getenv("KO_DATA_PATH"), action)
	file, _ := ioutil.ReadFile(filepath)
	res, _, err := c.client.Media.Upload(file, "tweet_gif")
	if err != nil {
		log.Println(err)
	}
	if res.MediaID > 0 {
		tweetParams.MediaIds = []int64{res.MediaID}
	}

	if _, _, err := c.client.Statuses.Update(msg, tweetParams); err != nil {
		log.Printf("error sending tweet %v\n", err)
	}
}
