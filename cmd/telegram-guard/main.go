package main

import (
	"os"

	"github.com/vedga/go-telegram-guard/pkg/api/telegram"
)

const (
	envBotToken = "BOT_TOKEN"
)

func main() {
	botToken := ""

	if value, present := os.LookupEnv(envBotToken); present {
		botToken = value
	}

	bot := telegram.NewBot(botToken)
	bot.Run()
	/*
		r := chi.NewRouter()

		// A good base middleware stack
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		r.Use(middleware.Timeout(60 * time.Second))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("hi"))
		})

		http.ListenAndServeTLS(":7001", "fullchain.pem", "privatekey.pem", r)

	*/
}
