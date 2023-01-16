package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"vasterd/max/chatter/models"
)

type echoServer struct{}

func (s echoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		fmt.Println("Error occurred: ", err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	var v models.ChatMessage
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		fmt.Println("Error occurred: ", err)
		return
	}

	log.Printf("received on channel: %d, message: %s", v.Channel, v.Msg)
	err = wsjson.Write(ctx, c, &v)
	if err != nil {
		fmt.Println("Error occurred: ", err)
		return
	}
	c.Close(websocket.StatusNormalClosure, "")
}
