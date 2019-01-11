package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/mongodb/mongo-go-driver/bson"
)

var storage Storage

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("WARN: No '.env' file was found, failing silently.")
	}

	err = storage.Connect(os.Getenv("DB_URI"), os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})
	r.HandleFunc("/videos", listVideos).Methods("GET")
	r.HandleFunc("/videos/{title}", viewVideo).Methods("GET")
	r.HandleFunc("/pubsubhubbub/subscribe", subscribe).Methods("POST")
	r.HandleFunc("/pubsubhubbub/feed", verify).Methods("GET")
	r.HandleFunc("/pubsubhubbub/feed", publish).Methods("POST")

	PORT := ":3000" // default server port
	if os.Getenv("PORT") != "" {
		PORT = ":" + os.Getenv("PORT")
	}

	log.Fatal(http.ListenAndServe(PORT, r))
}

func publish(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		handleError(w, err, "Converting xml to bytes")
		return
	}
	feed := Feed{}
	err = xml.Unmarshal(b, &feed)
	if err != nil {
		handleError(w, err, "Parsing feed")
		return
	}
	video := feed.Entry
	err = storage.InsertDoc("videos", video)
	if err != nil {
		handleError(w, err, "Inserting video")
		return
	}
	json.NewEncoder(w).Encode(video)
}

func verify(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	if challenge != "" {
		fmt.Fprintf(w, challenge)
	}
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	var d struct {
		HubCallback string `json:"hub.callback"`
		HubTopic    string `json:"hub.topic"`
	}
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		handleError(w, err, "Decoding subscription body")
		return
	}
	data := url.Values{
		"hub.callback": {d.HubCallback},
		"hub.topic":    {d.HubTopic},
		"hub.verify":   {"async"},
		"hub.mode":     {"subscribe"},
	}
	resp, err := http.PostForm("https://pubsubhubbub.appspot.com/subscribe", data)
	if err != nil {
		handleError(w, err, "Calling subscribe endpoint")
		return
	}
	defer resp.Body.Close()
	json.NewEncoder(w).Encode(struct {
		Success bool `json:"success"`
	}{
		resp.StatusCode != 200,
	})
}

func viewVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]
	video, err := storage.ViewDoc("videos", bson.M{"title": title})
	if err != nil {
		handleError(w, err, "Fetching video")
		return
	}
	json.NewEncoder(w).Encode(video)
}

func listVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := storage.ListDocs("videos")
	if err != nil {
		handleError(w, err, "Listing videos")
		return
	}
	json.NewEncoder(w).Encode(videos)
}

func handleError(w http.ResponseWriter, err error, from string) {
	log.Print(fmt.Errorf("[%v] while [%s]", err, from))
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Oups, I think you broke something...")
}
