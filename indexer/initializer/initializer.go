package initializer

import (
	"jmtp/indexer/commons"
	"net/http"
	"time"
)

func Initializer() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	commons.HttpClient = &http.Client{
		Timeout:   15 * time.Second,
		Transport: t,
	}
}
