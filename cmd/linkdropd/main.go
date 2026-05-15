package main

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
)

// function runs when we hit root endpoint
func getRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got / request\n")
    io.WriteString(w, "wassup\n")
}

// function runs when we hit hello endpoint
func getHello(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /hello request\n")
    io.WriteString(w, "hello world\n")
}

func main() {
    http.HandleFunc("/", getRoot)
    http.HandleFunc("/hello", getHello)
    
    err := http.ListenAndServe(":4545", nil)
     
    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server closed \n")
    } else if err != nil {
        fmt.Printf("error starting server: %s\n", err)
        os.Exit(1)
    }
}