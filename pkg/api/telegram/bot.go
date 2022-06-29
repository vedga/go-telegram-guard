package telegram

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/render"

	"github.com/vedga/go-telegram-guard/pkg/api/telegram/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envBindAddress = "BIND_ADDRESS"
	envCertificate = "HTTPS_CERTIFICATE"
	envPrivateKey  = "HTTPS_KEY"
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	requestProcessingTimeout = 60 * time.Second
)

// Bot represent the Telegram bot
type Bot struct {
	router                                       chi.Router
	botToken, bindAddress, fullChain, privateKey string
}

// NewBot return new bot instance
func NewBot(botToken string) *Bot {
	bot := &Bot{
		botToken: botToken,
		router:   chi.NewRouter(),
	}

	if value, present := os.LookupEnv(envBindAddress); present {
		bot.bindAddress = value
	} else {
		bot.bindAddress = ":7001"
	}

	if value, present := os.LookupEnv(envCertificate); present {
		bot.fullChain = value
	} else {
		bot.fullChain = "fullchain.pem"
	}

	if value, present := os.LookupEnv(envPrivateKey); present {
		bot.privateKey = value
	} else {
		bot.privateKey = "privatekey.pem"
	}

	// A good base middleware stack
	bot.router.Use(middleware.RequestID)
	bot.router.Use(middleware.RealIP)
	bot.router.Use(middleware.Logger)
	bot.router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	bot.router.Use(middleware.Timeout(requestProcessingTimeout))

	// Root router
	bot.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	// Router by webhook
	bot.router.Route("/{token}", func(r chi.Router) {
		r.Use(bot.PrepareWebhook)
		r.Post("/", onUpdate)
	})

	return bot
}

// Run the bot
func (bot *Bot) Run() error {
	server := &http.Server{Addr: bot.bindAddress, Handler: bot.router}

	tcpListener, e := net.Listen("tcp", server.Addr)
	if nil != e {
		return e
	}

	serverListener := newListener(tcpListener, bot.botToken)

	defer serverListener.Close()

	return server.ServeTLS(serverListener, bot.fullChain, bot.privateKey)
	//	http.ListenAndServeTLS(bot.bindAddress, bot.fullChain, bot.privateKey, bot.router)
}

// PrepareWebhook do prepare webhook context
func (bot *Bot) PrepareWebhook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Fetch webhook token
		token := chi.URLParam(r, "token")

		/*
			if err != nil {
				http.Error(w, http.StatusText(404), 404)
				return
			}
		*/

		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func onUpdate(w http.ResponseWriter, r *http.Request) {
	var update api.Update

	decodeJSON(r.Header, r.Body, &update)

	if source := update.Message; nil != source {
		method := "sendMessage"

		msg := &api.SendMessage{
			WebHookMessageReply: &method,
			Target:              strconv.Itoa(update.Message.ConversationPlace.ID),
			ReplyToMessageID:    &source.ID,
			Text:                "Усек!",
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, msg)
	}
}
