package chat

import (
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"sync"
)

type Handler struct {
	log       *slog.Logger
	ws        *websocket.Upgrader
	wsClients map[*websocket.Conn]struct{}
	rwmutex   *sync.RWMutex
}

func NewHandler(log *slog.Logger, ws *websocket.Upgrader) *Handler {
	return &Handler{
		log:     log,
		ws:      ws,
		rwmutex: &sync.RWMutex{},
	}
}

func (h *Handler) WsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "Chat.Handler.WsHandler"

	log := h.log.With("op", op)

	ws, err := h.ws.Upgrade(w, r, nil)
	if err != nil {
		log.Error("failed upgrade", "err", err)

		http.Error(w, "failed websocket upgrade", http.StatusBadRequest)
		return
	}

	defer ws.Close()

	log.Info("websocket connected",
		"ip", ws.RemoteAddr().String(),
	)
	h.rwmutex.Lock()
	h.wsClients[ws] = struct{}{}
	h.rwmutex.Unlock()
	go h.readFromClient(ws)
}

func (h *Handler) readFromClient(conn *websocket.Conn) {
	for {
		msg := new(WSMessage)
		if err := conn.ReadJSON(msg); err != nil {
			slog.Error("failed read from websocket", "err", err)
			break
		}
	}
	h.rwmutex.Lock()
	delete(h.wsClients, conn)
	h.rwmutex.Unlock()
}
