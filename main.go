package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.ParseInt(r.Header.Get("sleep"), 10, 64)
	if err != nil {
		n = 0
	}
	time.Sleep(time.Duration(n) * time.Second)
	_, _ = fmt.Fprintf(w, "Hello")
}

func main () {
	http.HandleFunc("/hello", HelloHandler)
	_ = http.ListenAndServe(":8878", nil)
}

