package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
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

	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		fmt.Println("Error occurred: ", err)
		return
	}

	log.Printf("received: %v", v)
	err = wsjson.Write(ctx, c, &v)
	if err != nil {
		fmt.Println("Error occurred: ", err)
		return
	}
	c.Close(websocket.StatusNormalClosure, "")
}
