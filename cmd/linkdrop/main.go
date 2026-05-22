package main

import (
	"fmt"
	"os"

	//"os/exec"
	"net/http"
	//"strings"
	//"io"
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type jsonBody struct {
	Link string `json:"link"`
}

type linkIDBody struct {
	ID int `json:"linkID"`
}

type ServerResponse struct {
	Response string `json:"DBResponse"`
}

type inboxLinkBucket struct {
	Id   int    `json:"entryID"`
	Link string `json:"entryLink"`
}

type InboxResponse struct {
	Inbox []inboxLinkBucket `json:"Inbox"`
}

type OpenLink struct {
	Id int `json:"ID`
}

func inbox(bearer string) {
	requestURL := fmt.Sprintf("http://localhost:4545/links")
	request, reqErr := http.NewRequest(http.MethodGet, requestURL, nil)

	if reqErr != nil {
		log.Fatal(reqErr)
	}

	request.Header.Set("Authorization", bearer)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
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

	for _, entry := range responseMsg.Inbox {
		counter += 1
		fmt.Printf("%d: %s\n", entry.Id, entry.Link)
	}

	//return
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

func open(id int, bearer string) {
	payload := OpenLink{
		Id: id,
	}

	l, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	bodyReader := bytes.NewReader(l)

	requestURL := fmt.Sprintf("http://localhost:4545/links/open")
	request, reqErr := http.NewRequest(http.MethodPost, requestURL, bodyReader)

	if reqErr != nil {
		fmt.Fprintf(os.Stderr, "request failed: %v\n", reqErr)
	}

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

// MAKE A HELPER FUNCTION TO DO ALL REQUESTS
// THIS WAY WE CAN GET RID OF ALL THE REPETITIVE STUFF BETWEEN SEND/INBOX/OPEN
func main() {
	if len(os.Args) >= 2 {
		cmdArg := os.Args[1]

		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		bearer := os.Getenv("BEARER_TOKEN")

		if cmdArg == "inbox" {
			inbox(bearer)
		} else if cmdArg == "send" {
			if len(os.Args) < 3 {
				fmt.Println("error - please provide a link")
				return
			} else {
				linkArg := os.Args[2]
				send(linkArg, bearer)
			}
		} else if cmdArg == "open" {
			idToOpen := os.Args[2]
			idInt, err := strconv.Atoi(idToOpen)
			if err != nil {
				log.Fatal(err)
			}
			open(idInt, bearer)
		}
	} else {
		fmt.Println("error - not enough args given")
		return
	}
}
