package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error occurred while running the client: %s\n", err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("please provide an address to listen on as the first argument")
	}

	address := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "http://"+address, nil)
	if err != nil {
		return fmt.Errorf("connection cannot be established: %s", err)
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	for {
		var message string

		fmt.Scanln(&message)
		err = wsjson.Write(ctx, conn, &message)
		if err != nil {
			return fmt.Errorf("message could not be sent: %s", err)
		}

		fmt.Println("Waiting for incoming message...")
		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			return fmt.Errorf("message could not be read: %s", err)
		}
		fmt.Printf("Received message: %s", v)
		break
	}

	conn.Close(websocket.StatusNormalClosure, "")
	return nil
}
