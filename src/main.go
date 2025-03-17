package main

import (
	"fmt"
	"github/skrebnevf/hs_code/pkg/handlers"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/supabase-community/supabase-go"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	b, err := BotInit(config.Token)
	if err != nil {
		log.Fatalf("cannot init bot, error: %v", err)
	}

	client, err := supabase.NewClient(config.DB.Url, config.DB.Key, &supabase.ClientOptions{})
	if err != nil {
		log.Printf("DB error: %v", err)
	}

	handlers.ChatID, err = strconv.ParseInt(config.ChatId, 10, 64)
	if err != nil {
		log.Printf("cannot convert chat id value to int64, err: %v", err)
	}

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Bot is running")
		})

		port := os.Getenv("PORT")
		if port == "" {
			port = "8000"
		}

		log.Printf("Listening on port %s for health checks...", port)
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		for {
			resp, err := http.Get(config.BotUrl)
			if err != nil {
				log.Printf("Request error: %v", err)
			} else {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("cannot get body request: %v", err)
				} else {
					log.Printf("Resp: %s, Time %v", body, time.Now())
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	handlers.CommandHandlers(b, client)
	handlers.TextHandlers(b, client)
	b.Start()
}
