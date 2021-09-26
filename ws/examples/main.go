package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RadekD/go-kit/ws"
	"github.com/RadekD/go-kit/ws/hooks/online"

)
/*
func main() {
	router := mux.NewRouter()

	router.HandleFunc("/websocket", socklace.ServeHTTP)
	router.PathPrefix("/").HandlerFunc(handleIndex)

	socklace.RegisterHandler("echo", wsEcho)

	weberr := http.ListenAndServe(":8081", router)
	log.Fatal(weberr)
}

func init() {
	socklace.AddPreHook(wsHelloHook)
	socklace.AddPreHook(func(c *socklace.Connection) {
		fmt.Println("pre online: ", online.Get())

		c.Values.Store("user", "this is user")
	})
	socklace.AddPostHook(func(c *socklace.Connection) {

		fmt.Println("post online: ", online.Get())

		//fmt.Println(c.Store.Get("user"))
	})
}

func wsEcho(r *socklace.Request) {
	_, x := pubsub.Publish("global.test", "test")

	fmt.Println("sent to:", x)

	//r.Respond(r.Data)
}
func wsHelloHook(c *socklace.Connection) {
	c.Send(0, "hello", "world")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "test")
}
*/
