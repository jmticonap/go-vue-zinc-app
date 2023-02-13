package main

import (
	"flag"
	"fmt"
	"io/fs"

	// "io/ioutil"
	"jmtp/indexer/commons"
	"jmtp/indexer/initializer"
	"jmtp/indexer/parsers"
	"jmtp/indexer/routers"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
)

func init() {
	if os.Getenv("PROFILE") != "PROD" {
		os.Setenv("PROFILE", "DEV")
	}

	if os.Getenv("PROFILE") == "DEV" {
		os.Setenv("PORT", "8081")
		os.Setenv("ZINC_SEARCH_HOST", "http://localhost:4080")
	}
	initializer.Initializer()
}

func main() {
	flag.Int64Var(&commons.Limit, "limit", int64(500), "limit for the quantity of rows")
	flag.Int64Var(&commons.MaxSizeBatch, "max_size_batch", int64(4_000_000), "top limit for emails batch that will be send in MB")
	flag.StringVar(&commons.DbPath, "db_path", "", "path of the emails db")
	flag.StringVar(&commons.Auth, "auth", "admin:Complexpass#123", "credentials for the end-point")
	flag.BoolVar(&commons.NewData, "new_data", false, "detereminate if new data is load to zincsearch")
	flag.Parse()

	//commons.MailInfo, _ = ioutil.ReadDir(commons.DbPath)
	commons.FilesQuantity = fileCount(commons.DbPath)

	fmt.Printf("FILES_QUANTITY = %v\n", commons.FilesQuantity)

	if commons.NewData {
		parsers.DataLoader()
	}

	go parsers.SendMail()

	port := os.Getenv("PORT")

	router := chi.NewRouter()

	routers.RootRoute(router, commons.DbPath, commons.Auth, int(commons.Limit))

	if port == "" {
		port = "8081"
	}
	green := color.New(color.Bold, color.BgGreen).SprintFunc()
	fmt.Printf("%s\n", green(" SERVER "))
	fmt.Printf("http://localhost:%s/debug\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func fileCount(path string) int {
	i := 0
	filepath.WalkDir(commons.DbPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			i++
		}
		return nil
	})

	return i
}
