package telegram

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
)

// listener is wrapper of net.Listener interface
type listener struct {
	inner                 net.Listener
	botToken              string
	inService, offService sync.Once
}

// newListener return wrapper for net.Listener interface
func newListener(inner net.Listener, botToken string) net.Listener {
	return &listener{
		inner:    inner,
		botToken: botToken,
	}
}

// Accept waits for and returns the next connection to the listener.
func (l *listener) Accept() (net.Conn, error) {
	l.inService.Do(func() {
		// Called once before we ready to accept incoming connections
		req := struct {
			WebHookURL string `json:"url,omitempty"`
		}{
			WebHookURL: fmt.Sprintf("https://star.vedga.com:88/%s", l.botToken),
		}

		response := struct {
			Processed   bool   `json:"ok,omitempty"`
			Result      bool   `json:"result,omitempty"`
			Code        int    `json:"error_code,omitempty"`
			Description string `json:"description,omitempty"`
		}{}

		httpClient.Call(context.Background(), MethodPOST, l.botToken, "setWebhook", &response, &req)

		log.Default().Print(response)
	})

	return l.inner.Accept()
}

// Close closes the listener.
func (l *listener) Close() error {
	l.offService.Do(func() {
		// Called once before stop accepting incoming connections
	})

	return l.inner.Close()
}

// Addr returns the listener's network address.
func (l *listener) Addr() net.Addr {
	return l.inner.Addr()
}
