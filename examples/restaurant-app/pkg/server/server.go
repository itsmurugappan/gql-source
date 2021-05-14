package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/allegro/bigcache"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"

	"github.com/itsmurugappan/gql-source/examples/restaurant-app/schema"
)

type graphQLServer struct {
	itemChannels map[uuid.UUID]chan *schema.Item
	infoChannels map[uuid.UUID]chan *schema.Info
	mutex        sync.Mutex
	cache        *bigcache.BigCache
}

func NewGraphQLServer() (*graphQLServer, error) {
	c, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	loadInitialMenu(c)

	return &graphQLServer{
		mutex:        sync.Mutex{},
		cache:        c,
		itemChannels: map[uuid.UUID]chan *schema.Item{},
		infoChannels: map[uuid.UUID]chan *schema.Info{},
	}, nil
}

func loadInitialMenu(c *bigcache.BigCache) {
	items, _ := json.Marshal(map[string]*schema.Item{
		"Appetizer-FriedMozzarella": {ItemType: schema.ItemTypesAppetizer, Name: "Fried Mozzarella"},
		"Appetizer-ChickenWings":    {ItemType: schema.ItemTypesAppetizer, Name: "Chicken Wings"},
		"Entree-ChickenScampi":      {ItemType: schema.ItemTypesEntree, Name: "Chicken Scampi"},
		"Entree-CheeseRavioli":      {ItemType: schema.ItemTypesEntree, Name: "Cheese Ravioli"},
		"Dessert-Tiramisu":          {ItemType: schema.ItemTypesDessert, Name: "Tiramisu"},
	})
	info, _ := json.Marshal(&schema.Info{
		Address: "1000 Orange Ave Cypress California 90630",
		Hours:   "15:00 — 22:00 on All Days",
	})
	c.Set("item", items)
	c.Set("info", info)
}

func (s *graphQLServer) Serve(route string, port int) error {
	mux := http.NewServeMux()
	mux.Handle(
		route,
		handler.GraphQL(schema.NewExecutableSchema(schema.Config{Resolvers: s}),
			handler.WebsocketKeepAliveDuration(10*time.Second),
			handler.WebsocketUpgrader(websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
				HandshakeTimeout: 5 * time.Second,
			}),
		),
	)
	mux.Handle("/playground", handler.Playground("GraphQL", route))

	handler := cors.AllowAll().Handler(mux)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}

func (s *graphQLServer) MutateItem(ctx context.Context, itemType schema.ItemTypes, name string, action schema.Actions) (*schema.Item, error) {
	i := &schema.Item{
		CreatedAt: time.Now().UTC(),
		ItemType:  itemType,
		Name:      name,
		Action:    action,
	}
	if err := s.loadData("item", i); err != nil {
		return nil, err
	}
	// Notify new message
	s.mutex.Lock()
	for _, ch := range s.itemChannels {
		ch <- i
	}
	s.mutex.Unlock()
	return i, nil
}

func (s *graphQLServer) UpdateInfo(ctx context.Context, hours, address string) (*schema.Info, error) {
	i := &schema.Info{
		Hours:   hours,
		Address: address,
	}
	if err := s.loadData("info", i); err != nil {
		return nil, err
	}
	// Notify new message
	s.mutex.Lock()
	for _, ch := range s.infoChannels {
		ch <- i
	}
	s.mutex.Unlock()
	return i, nil
}

func (s *graphQLServer) Items(ctx context.Context) ([]*schema.Item, error) {
	data, err := s.getData("item")
	if err != nil {
		return nil, err
	}
	return data.([]*schema.Item), nil
}

func (s *graphQLServer) AllInfo(ctx context.Context) (*schema.Info, error) {
	data, err := s.getData("info")
	if err != nil {
		return nil, err
	}
	return data.(*schema.Info), nil
}

func (s *graphQLServer) ItemChanged(ctx context.Context) (<-chan *schema.Item, error) {
	// Create new channel for request
	itemsCh := make(chan *schema.Item)
	id, _ := uuid.NewUUID()
	log.Printf("start item subscription for id %v\n", id)
	s.mutex.Lock()
	s.itemChannels[id] = itemsCh
	s.mutex.Unlock()

	go func() {
		<-ctx.Done()
		log.Printf("close subscription for id %v\n", id)
		s.mutex.Lock()
		delete(s.itemChannels, id)
		s.mutex.Unlock()
	}()

	return itemsCh, nil
}

func (s *graphQLServer) InfoChanged(ctx context.Context) (<-chan *schema.Info, error) {
	// Create new channel for request
	infoCh := make(chan *schema.Info)
	id, _ := uuid.NewUUID()
	log.Printf("start info subscription for id %v\n", id)
	s.mutex.Lock()
	s.infoChannels[id] = infoCh
	s.mutex.Unlock()

	go func() {
		<-ctx.Done()
		log.Printf("close subscription for id %v\n", id)
		s.mutex.Lock()
		delete(s.infoChannels, id)
		s.mutex.Unlock()
	}()

	return infoCh, nil
}

func (s *graphQLServer) getData(key string) (interface{}, error) {
	data, err := s.cache.Get(key)
	if err != nil {
		return nil, err
	}
	log.Printf("fetched data %s\n", string(data))
	switch key {
	case "item":
		var iMap map[string]*schema.Item
		if err := json.Unmarshal(data, &iMap); err != nil {
			return nil, err
		}
		var values []*schema.Item
		for _, v := range iMap {
			values = append(values, v)
		}
		return values, nil
	case "info":
		var i *schema.Info
		if err := json.Unmarshal(data, &i); err != nil {
			return nil, err
		}
		return i, nil
	}
	return nil, nil
}

func (s *graphQLServer) loadData(key string, data interface{}) error {
	switch key {
	case "item":
		cd, err := s.cache.Get(key)
		if err != nil {
			return err
		}
		var iMap map[string]*schema.Item
		if err := json.Unmarshal(cd, &iMap); err != nil {
			return err
		}
		itemData := data.(*schema.Item)
		itemKey := fmt.Sprintf("%s-%s", itemData.ItemType, strings.Replace(itemData.Name, " ", "", -1))
		if itemData.Action == schema.ActionsAdd {
			iMap[itemKey] = itemData
		} else {
			delete(iMap, itemKey)
		}
		newValue, _ := json.Marshal(iMap)
		log.Printf("menu list after delete %s\n", string(newValue))
		s.cache.Set(key, newValue)
		return nil
	case "info":
		newValue, _ := json.Marshal(data)
		s.cache.Set(key, newValue)
		return nil
	}
	return nil
}

func (s *graphQLServer) Mutation() schema.MutationResolver {
	return s
}

func (s *graphQLServer) Query() schema.QueryResolver {
	return s
}

func (s *graphQLServer) Subscription() schema.SubscriptionResolver {
	return s
}
