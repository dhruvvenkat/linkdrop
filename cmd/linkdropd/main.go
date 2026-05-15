package main

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
)

// function runs when we hit root endpoint
func throwErr(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got incorrect request\n")
    io.WriteString(w, "404 - page not found\n")
}

// function runs when we hit hello endpoint
func getHealth(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /health request\n")
    io.WriteString(w, "healthy\n")
}

func main() {
    http.HandleFunc("/", throwErr)
    http.HandleFunc("/health", getHealth)
    
    
    err := http.ListenAndServe(":4545", nil)
     
    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server closed \n")
    } else if err != nil {
        fmt.Printf("error starting server: %s\n", err)
        os.Exit(1)
    }
}