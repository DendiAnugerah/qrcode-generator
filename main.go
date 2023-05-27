package main

import (
	"dend-qrcode/cmd/handler"
	"dend-qrcode/cmd/helper"
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	flag.Parse()

	router := mux.NewRouter()
	router.HandleFunc("/qrcode", handler.HandlerQRCode).Methods("POST")

	fmt.Println("Server started on http://localhost:8080")
	
	http.ListenAndServe(":8080", router)
	helper.Command()
}
