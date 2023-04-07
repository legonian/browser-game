package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const reqBodyLimit = 100_000
const passwordMaxSize = 100

type Server struct {
	serveMux http.ServeMux
	db       *DB
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Error  *string `json:"error"`
	Result *string `json:"result"`
}

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`

	Data json.RawMessage `json:"data"`
}

func NewServer() (*Server, error) {
	db, err := NewDB()
	if err != nil {
		return nil, fmt.Errorf("NewDB: %w", err)
	}
	gs := &Server{
		db: db,
	}
	gs.serveMux.Handle("/", http.FileServer(http.Dir("static")))
	gs.serveMux.HandleFunc("/ws", gs.ws)
	gs.serveMux.HandleFunc("/auth", gs.auth)
	return gs, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveMux.ServeHTTP(w, r)
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var login Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		wrapError(w, "pasing error")
		return
	}

	if login.Username == "" {
		wrapError(w, "no username")
		return
	}
	if login.Password == "" {
		wrapError(w, "no password")
		return
	}

	if passwordMaxSize < len(login.Password) {
		wrapError(w, "password is too long")
		return
	}

	err = s.db.CreateUser(login.Username, login.Password)
	if errors.Is(err, ErrUserExistsInvalid) {
		wrapError(w, "incorrect password")
		return
	}

	cookie := http.Cookie{Name: "username", Value: login.Username}
	http.SetCookie(w, &cookie)

	if errors.Is(err, ErrUserExistsValid) {
		log.Printf("Login: %v", login.Username)
		wrapResult(w, "done")
		return
	}

	log.Printf("SignIn: %v", login.Username)
	wrapResult(w, "done")
}

func (s *Server) ws(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "closed")

	conn.SetReadLimit(reqBodyLimit)

	for {
		var msg Message
		err = wsjson.Read(r.Context(), conn, &msg)
		if err != nil {
			log.Println("wsjson.Read:", err)
			return
		}

		switch msg.Type {
		case "chat":
			log.Printf("chat: %v", string(msg.Data))
		case "draw":
			log.Printf("draw: len=%v", len(msg.Data))
		default:
			log.Printf("bad msg.Type: %v", msg.Type)
		}
	}
}

func wrapError(w http.ResponseWriter, errText string) {
	w.Header().Set("Content-Type", "application/json")
	var resp Response
	resp.Error = &errText
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
}

func wrapResult(w http.ResponseWriter, res string) {
	w.Header().Set("Content-Type", "application/json")
	var resp Response
	resp.Result = &res
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
	}
}
