package main

import (
	"context"
	"log"
	"net/http"

	"github.com/itsmurugappan/gql-source/examples/restaurant-app/pkg/smbot"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	ctx := context.Background()
	bot := smbot.InitBot()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	h, err := cloudevents.NewHTTPReceiveHandler(ctx, p, bot.ProcessCE)
	if err != nil {
		log.Fatalf("failed to create handler: %s", err.Error())
	}

	http.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}))

	http.ListenAndServe(":8080", nil)
}
