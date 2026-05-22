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

func inbox(bearer string) {
        requestURL := fmt.Sprintf("http://localhost:4545/links")
        request, reqErr := http.NewRequest(http.MethodGet, requestURL, nil)
        
        if reqErr != nil {
            log.Fatal(reqErr)
        }
        
        request.Header.Set("Authorization", bearer)
        
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

func send(link string, bearer string) {
        payload := jsonBody{
            Link: link,
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
        request.Header.Set("Authorization", bearer)
        request.Header.Set("Content-Type", "application/json")
        
        client := http.Client{
            Timeout: 5 * time.Second,
        }
        
        res, err := client.Do(request)
        if err != nil {
        	log.Fatal(err)
        }
        
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

func main() {
	if len (os.Args) >= 2 {
		cmdArg := os.Args[1]
		bearerArg := os.Args[2]

		if cmdArg == "inbox" {
			inbox(bearerArg)
		} else if cmdArg == "send" {
			if len (os.Args) == 2 {
				fmt.Println("error - please provide a link")
				return
			} else {
				linkArg := os.Args[3]
				send(linkArg, bearerArg)	
			}
		}
	} else {
		fmt.Println("error - not enough args given")
		return
	}
}
