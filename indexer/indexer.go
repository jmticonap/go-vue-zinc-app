package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
		loadNewData(path, auth, int(limit), DataParserV2)
	}

	port := os.Getenv("PORT")

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Mount("/debug", middleware.Profiler())

	router.Get("/new_data_v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("loading data...")
		loadNewData(path, auth, int(limit), DataParser)
		fmt.Println("data loaded.")
	})

	router.Get("/new_data_v2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("loading data...")
		loadNewData(path, auth, int(limit), DataParserV2)
		fmt.Println("data loaded.")
	})

	if port == "" {
		port = "8081"
	}
	log.Printf("Server is running\n\thttp://localhost:%s/debug", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
