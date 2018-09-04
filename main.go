package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr = flag.String("addr", "localhost:8001", "http listen address")

	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type echoServer struct{}

func (es echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}

		if err = c.WriteMessage(mt, message); err != nil {
			break
		}
	}
	c.Close()
}

func main() {
	flag.Parse()

	var echo echoServer
	srv := &http.Server{
		Handler:      echo,
		Addr:         *addr,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
