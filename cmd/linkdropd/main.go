package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "strings"
//    "runtime"
    "github.com/joho/godotenv"
)

const keyServerAddr = "serverAddr"

// remember: lowercase fields are private to the struct
type linkStruct struct {
    Link string
}

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

func links (w http.ResponseWriter, r *http.Request) {
    var l linkStruct
    
    err := json.NewDecoder(r.Body).Decode(&l)
    
    if err != nil {
        // display error on client side instead of server side
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    errDotEnv := godotenv.Load()
    
    if errDotEnv != nil {
        fmt.Fprintln(os.Stderr, "error sourcing bearer token: %s", errDotEnv)
    }    
    
    // authenticated := false

    // basic bearer token auth - check if incoming bearer token is the same as the one in .env    
    if os.Getenv("BEARER_TOKEN") == r.Header.Get("Authorization") {
        fmt.Println("authenticated!")
    } else {
        http.Error(w, "401 unauthorized", http.StatusUnauthorized)
        return
    }

    if l.Link == "" {
        http.Error(w, "400 bad request - need to include a link", http.StatusBadRequest)
        return
    }
    
    // check to see if links are valid
    if !strings.HasPrefix(l.Link, "http") && !strings.HasPrefix(l.Link, "https") {
        http.Error(w, "400 bad request - not a real link", http.StatusBadRequest)
        return
    }
    
    fmt.Printf("new link: %s\n", l.Link)
    io.WriteString(w, "link receieved!!\n")
    
    io.WriteString(w, "opening link...\n")
    errOpeningURL := exec.Command("xdg-open", l.Link).Start()
    
    if errOpeningURL != nil {
        io.WriteString(w, "500 error - error opening URL: ")
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    io.WriteString(w, "201 - link opened\n")    
} 

func main() {
    http.HandleFunc("/", throwErr)
    http.HandleFunc("/health", getHealth)
    http.HandleFunc("/links", links)    
                
    err := http.ListenAndServe(":4545", nil)
     
    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server closed \n")
    } else if err != nil {
        fmt.Printf("error starting server: %s\n", err)
        os.Exit(1)
    }
}