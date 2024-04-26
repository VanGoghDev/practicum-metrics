package main

import (
	"net/http"
)

func main() {
	// config

	// logger

	// storage

	// router
	mux := setupMux()
	mux.HandleFunc(`/`, testHandler)
	mux.HandleFunc(`/update/`, updateHandler)

	http.ListenAndServe(`:8333`, mux)
}

func testHandler(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello"))
}

func updateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	// p := strings.Split(req.URL.Path, "/")

}

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	return mux
}
