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
	"database/sql"
	"log"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const keyServerAddr = "serverAddr"

// remember: lowercase fields are private to the struct
type linkStruct struct {
	Link string
}

type response struct {
	DBResponse string `json:"DBResponse"`
}

type inboxLinkBucket struct {
	Id   int    `json:"entryID"`
	Link string `json:"entryLink"`
}

type inboxResponse struct {
	Inbox []inboxLinkBucket
}

type recieveID struct {
	Id int `json:"ID"`
}

type GenericResponse struct {
	Msg string `json:"msg"`
}

// function runs when we hit root endpoint
func throwErr(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got incorrect request: " + r.URL.Path)
	io.WriteString(w, "404 - page not found\n")
}

// function runs when we hit hello endpoint
func getHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /health request\n")
	io.WriteString(w, "healthy\n")
}

func addLink(w http.ResponseWriter, l linkStruct) {
	dbPath := os.Getenv("LINKDROP_DB_PATH")
	if dbPath == "" {
		dbPath = "myLinks.db" // if database doesn't exist after running install.sh, just create it in the installation folder
	}

	db, dbErr := sql.Open("sqlite3", dbPath)

	if dbErr != nil {
		log.Fatal(dbErr)
	}

	defer db.Close()

	// statement to generate the table if it doesn't exist
	// ** for future: extract site name from the site metadata and store that in the name column
	sqlCreate := `
        CREATE TABLE IF NOT EXISTS links (
            id INTEGER PRIMARY KEY,
            name TEXT,
            url TEXT NOT NULL
        );   
    `

	_, dbCreateErr := db.Exec(sqlCreate)

	if dbCreateErr != nil {
		log.Fatal(dbCreateErr)
	}
	log.Println("Table 'links' created!")

	if l.Link == "" {
		http.Error(w, "400 bad request - need to include a link", http.StatusBadRequest)
		return
	}

	// check to see if links are valid
	if !strings.HasPrefix(l.Link, "http") && !strings.HasPrefix(l.Link, "https") {
		http.Error(w, "400 bad request - not a real link", http.StatusBadRequest)
		return
	}

	// confirm
	fmt.Printf("new link: %s\n", l.Link)
	// io.WriteString(w, "link receieved!!\n")

	// open link on user pc
	//io.WriteString(w, "opening link...\n")
	errOpeningURL := exec.Command("xdg-open", l.Link).Start()

	_, err := db.Exec("INSERT INTO links(url) VALUES(?)", l.Link)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("new link inserted into links table!!")

	if errOpeningURL != nil {
		//io.WriteString(w, "500 error - error opening URL: ")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//io.WriteString(w, "201 - link opened\n")

	var lastID (int)
	err = db.QueryRow("SELECT MAX(id) FROM links").Scan(&lastID)

	returned := response{DBResponse: "success! link entered at id " + strconv.Itoa(lastID) + "\n"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(returned)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func showInbox(w http.ResponseWriter) {
	db, dbErr := sql.Open("sqlite3", "myLinks.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id,url FROM links")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	linkSlice := []inboxLinkBucket{}

	for rows.Next() {
		var id int
		var url string
		if err := rows.Scan(&id, &url); err != nil {
			log.Fatal(err)
		}

		bucketToAdd := inboxLinkBucket{Id: id, Link: url}

		linkSlice = append(linkSlice, bucketToAdd)

	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	returned := inboxResponse{Inbox: linkSlice}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(returned)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func deleteLink(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("BEARER_TOKEN") == r.Header.Get("Authorization") {
		fmt.Println("authenticated!")
	} else {
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
		return
	}

	var d recieveID

	db, dbErr := sql.Open("sqlite3", "myLinks.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("deleting id %d\n", d.Id)

	_, err = db.Exec("DELETE FROM links WHERE id = ?", d.Id)
	if err != nil {
		log.Fatal(err)
	}

	dr := GenericResponse{Msg: "link " + strconv.Itoa(d.Id) + " deleted successfully!"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(dr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func links(w http.ResponseWriter, r *http.Request) {

	errDotEnv := godotenv.Load()

	if errDotEnv != nil {
		log.Fatal(errDotEnv)
	}

	// basic bearer token auth - check if incoming bearer token is the same as the one in .env
	if os.Getenv("BEARER_TOKEN") == r.Header.Get("Authorization") {
		fmt.Println("authenticated!")
	} else {
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodPost {
		var l linkStruct

		err := json.NewDecoder(r.Body).Decode(&l)
		if err != nil {
			// display error on client side instead of server side
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		addLink(w, l)
	} else if r.Method == http.MethodGet {
		showInbox(w)
	} else {
		http.Error(w, "unsuppored request type", http.StatusBadRequest)
		return
	}

}

func open(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("BEARER_TOKEN") == r.Header.Get("Authorization") {
		fmt.Println("authenticated!")
	} else {
		http.Error(w, "401 unauthorized", http.StatusUnauthorized)
		return
	}

	var o recieveID
	db, dbErr := sql.Open("sqlite3", "myLinks.db")
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	err := json.NewDecoder(r.Body).Decode(&o)
	if err != nil {
		// display error on client side instead of server side
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//fmt.Println(o.Id)

	var link string
	err = db.QueryRow("SELECT url FROM links WHERE id = ?", o.Id).Scan(&link)
	if err != nil {
		log.Fatal(err)
	}

	err = exec.Command("xdg-open", link).Start()

	or := GenericResponse{Msg: "link at id " + strconv.Itoa(o.Id) + " opened!"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(or)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", throwErr)
	http.HandleFunc("/health", getHealth)
	http.HandleFunc("/links", links)
	http.HandleFunc("/links/open", open)
	http.HandleFunc("/links/delete", deleteLink)

	err := http.ListenAndServe("0.0.0.0:4545", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed \n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
