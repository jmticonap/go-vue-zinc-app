package routers

import (
	"fmt"
	"jmtp/indexer/parsers"
	"net/http"
	_ "net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RootRoute(r *chi.Mux, path, auth string, limit int) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/debug", middleware.Profiler())

	r.Get("/new_data_v1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("loading data...")
		parsers.LoadNewData(path, auth, int(limit), parsers.DataParserV1)
		fmt.Println("data loaded.")
	})

	r.Get("/new_data_v2", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("loading data...")
		parsers.LoadNewData(path, auth, int(limit), parsers.DataParserV2)
		fmt.Println("data loaded.")
	})
}
