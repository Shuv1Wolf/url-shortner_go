package main

import (
	"url-shortener/internal/config"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()
	//TODO: init loger: slog

	//TODO: init storage: sqlite3

	//TODO: init router: chi, "chi render"

	//TODO: run server
}
