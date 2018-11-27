package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/manucorporat/try"
)

func main() {
	fmt.Println("Starting ...")
	var router = mux.NewRouter()
	router.HandleFunc("/{function}", executeFunction).Methods("GET")
	fmt.Println("Started!")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", router))
}

func executeFunction(w http.ResponseWriter, r *http.Request) {
	urlVars := mux.Vars(r)
	functionName := urlVars["function"]
	queryVars := r.URL.Query()
	object := funcMap[functionName]
	if object != nil {
		try.This(func() {
			result := object.(func(url.Values) string)(queryVars)
			json.NewEncoder(w).Encode(map[string]string{"functionName": functionName, "args": fmt.Sprint(queryVars), "message": result})
		}).Catch(func(e try.E) {
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprint(e)})
		})

	} else {
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid_function_name"})
	}
}

func noArgs(function func() string) func(url.Values) string {
	return func(_ url.Values) string { return function() }
}
