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
    "log"
    "strconv"
)

type jsonBody struct {
    Link string `json:"link"`
}

type ServerResponse struct {
    Response string `json:"DBResponse"`
}

type InboxResponse struct {
    Inbox []string `json:"Inbox"`
}

func main() {
    linkArg := os.Args[1]
    bearerArg := os.Args[2]        

    if linkArg == "inbox" {
        requestURL := fmt.Sprintf("http://localhost:4545/links")
        request, reqErr := http.NewRequest(http.MethodGet, requestURL, nil)
        
        if reqErr != nil {
            log.Fatal(reqErr)
        }
        
        request.Header.Set("Authorization", bearerArg)
        
        client := http.Client {
            Timeout: 5* time.Second,
        }
        
        resp, err := client.Do(request)
        defer resp.Body.Close()
        
        if err != nil {
            log.Fatal(err)
        }
        
        var responseMsg InboxResponse
        err = json.NewDecoder(resp.Body).Decode(&responseMsg)
        if err != nil {
            fmt.Println("error decoding response JSON: ", err)
            return
        }
        
        counter := 0
        
        for _, link := range responseMsg.Inbox {
            counter += 1
            fmt.Printf(strconv.Itoa(counter) + ": " + link + "\n")
        }
        
        return        
    }           
                                 
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
        fmt.Fprintf(os.Stderr, "request failed: %v\n", reqErr)
    }
    
    // request.Header.Set("Authorization", "Bearer " + bearerArg)
    request.Header.Set("Authorization", bearerArg)
    request.Header.Set("Content-Type", "application/json")
    
    client := http.Client{
        Timeout: 5 * time.Second,
    }
    
    res, err := client.Do(request)
    defer res.Body.Close()
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "client: error making http request: %s\n", err)
        os.Exit(1)
    }
    
    var responseMsg ServerResponse
    err = json.NewDecoder(res.Body).Decode(&responseMsg)
    if err != nil {
        fmt.Println("error decoding response JSON: ", err)
        return
    }
    
    fmt.Printf(responseMsg.Response)
}