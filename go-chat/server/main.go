package main

import (
  "fmt"
  "context"
  "time"
  "log"
  "net"
  "net/http"
  "os"
  "os/signal"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error occurred", err)
    return
	}
	log.Printf("listening on http://%v", l.Addr())

  s := &http.Server{
    Handler: echoServer{},
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
  }

  errc := make(chan error, 1)
  go func() {
    errc <- s.Serve(l)
  }()

  sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
  	case err := <-errc:
  		log.Printf("failed to serve: %v", err)
  	case sig := <-sigs:
  		log.Printf("terminating: %v", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	s.Shutdown(ctx)
}
