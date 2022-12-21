package main

import "os"

func init() {
	if os.Getenv("PROFILE") != "PROD" {
		os.Setenv("PROFILE", "DEV")
	}

	if os.Getenv("PROFILE") == "DEV" {
		os.Setenv("PORT", "8081")
		os.Setenv("ZINC_SEARCH_HOST", "http://localhost:4080")
	}
}
