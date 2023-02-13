package routers

import (
	"fmt"
	"jmtp/indexer/parsers"
	"net/http"
	_ "net/http/pprof"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RootRoute(r *chi.Mux, path, auth string, limit int) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Mount("/debug", middleware.Profiler())

	r.Get("/load_db", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("loading data...")
		parsers.DataLoader()
		fmt.Println("data loaded.")
	})
	r.Get("/load_db_cnt", func(w http.ResponseWriter, r *http.Request) {
		wg := &sync.WaitGroup{}
		wg.Add(2)
		go parsers.ReadMails(wg)
		go parsers.ParseMails(wg, w)
		// go parsers.SendMail(wg)
		wg.Wait()
		w.Write([]byte("--DONE--"))
	})
}
