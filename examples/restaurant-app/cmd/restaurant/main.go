package main

import (
	"log"

	"github.com/itsmurugappan/gql-source/examples/restaurant-app/pkg/server"
)

func main() {
	s, err := server.NewGraphQLServer()
	if err != nil {
		log.Fatal(err)
	}

	err = s.Serve("/graphql", 8080)
	if err != nil {
		log.Fatal(err)
	}
}
