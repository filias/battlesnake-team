package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type ContentType int

const (
	Food ContentType = iota
	Head
	Body
	Empty
)

type CoordExtended struct {
	score   int
	content ContentType
}

type GameBoardExtended [][]CoordExtended

// --

type GameState struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type Game struct {
	ID      string  `json:"id"`
	Ruleset Ruleset `json:"ruleset"`
	Timeout int32   `json:"timeout"`
}

type Ruleset struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`

	// Used in non-standard game modes
	Hazards []Coord `json:"hazards"`
}

type Battlesnake struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Health  int32   `json:"health"`
	Body    []Coord `json:"body"`
	Head    Coord   `json:"head"`
	Length  int32   `json:"length"`
	Latency string  `json:"latency"`

	// Used in non-standard game modes
	Shout string `json:"shout"`
	Squad string `json:"squad"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Response Structs

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type BattlesnakeMoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

// HTTP Handlers

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := info()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode start json, %s", err)
		return
	}

	start(state)

	// Nothing to respond with here
}

func HandleMove(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}

	response := move(state)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}

func HandleEnd(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode end json, %s", err)
		return
	}

	end(state)

	// Nothing to respond with here
}

// Main Entrypoint

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	log.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
