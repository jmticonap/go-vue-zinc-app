package main

import (
	"flag"
	"jmtp/indexer/parsers"
	"jmtp/indexer/routers"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	var limit int64
	var path, auth, speed string
	var newData bool

	flag.Int64Var(&limit, "limit", int64(500), "limit for the quantity of rows")
	flag.StringVar(&path, "path", "", "path of the emails db")
	flag.StringVar(&auth, "auth", "admin:Complexpass#123", "credentials for the end-point")
	flag.StringVar(&speed, "speed", "slow", "painting speed")
	flag.BoolVar(&newData, "new_data", true, "detereminate if new data is load to zincsearch")
	flag.Parse()

	if newData {
		parsers.LoadNewData(path, auth, int(limit), parsers.DataParserV2)
	}

	port := os.Getenv("PORT")

	router := chi.NewRouter()

	routers.RootRoute(router, path, auth, int(limit))

	if port == "" {
		port = "8081"
	}

	log.Printf("Server is running\n\thttp://localhost:%s/debug", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
