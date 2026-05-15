package main

import (
    "fmt"
    "os"
    //"os/exec"
    "net/http"
    //"strings"
    //"io"    
    "encoding/json"
    "bytes"
    "time"
)

type jsonBody struct {
    Link string `json:"link"`
}

func main() {
    linkArg := os.Args[1]
    
    payload := jsonBody{
        Link: linkArg,
    }

    l, err := json.Marshal(payload)
            
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to marshal JSON: %v\n", err)
    }

    bodyReader := bytes.NewReader(l)
        
    requestURL := fmt.Sprintf("http://localhost:4545/links")
    request, reqErr := http.NewRequest(http.MethodPost, requestURL, bodyReader)
   
    if reqErr != nil {
        fmt.Fprintf(os.Stderr, "request failed: %v\n", err)
    }
    
    request.Header.Set("Content-Type", "application/json")
    
    client := http.Client{
        Timeout: 5 * time.Second,
    }
    
    res, err := client.Do(request)
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "client: error making http request: %s\n", err)
        os.Exit(1)
    }
    
    defer res.Body.Close()
}