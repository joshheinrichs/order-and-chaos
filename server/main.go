package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joshheinrichs/order-and-chaos/server/config"
)

var mainConfig *config.Config
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	NewUser(ws)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("https://%s%s%s", mainConfig.Website.URL, mainConfig.Website.HTTPSPort, r.RequestURI)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	mainConfig, err = config.ReadFile("config.gcfg")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/ws", wsHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(mainConfig.Website.Directory)))
	http.Handle("/", r)

	go func() {
		log.Printf("Serving HTTP on %s\n", mainConfig.Website.HTTPPort)
		err := http.ListenAndServe(mainConfig.Website.HTTPPort, http.HandlerFunc(redirectHandler))
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Printf("Serving HTTPS on %s\n", mainConfig.Website.HTTPSPort)
	err = http.ListenAndServeTLS(mainConfig.Website.HTTPSPort, mainConfig.Website.Cert, mainConfig.Website.Key, nil)
	if err != nil {
		log.Fatal(err)
	}
}
