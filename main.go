package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type cmd struct {
	Type    string  `json:"type"`
	ID      string  `json:"id"`
	URI     string  `json:"uri"`
	Payload payload `json:"payload,omitempty"`
}

type payload struct {
	Volume  int    `json:"volume,omitempty"`
	Message string `json:"message,omitempty"`
	ID      string `json:"id,omitempty"`
	Mute    bool   `json:"mute,omitempty"`
}

var (
	dir, dirErr = filepath.Abs(filepath.Dir(os.Args[0]))
)

// generate uniq id for each call
func getID(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func command(addr net.IP, cmd *cmd, data *LGPair) {
	client := connect(addr, data)
	defer client.Close()

	dict, err := json.Marshal(cmd)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = client.WriteMessage(websocket.TextMessage, dict)
	if err != nil {
		log.Println("Try again...", err)
	}
	_, message, err := client.ReadMessage()
	if err != nil {
		log.Println("Error while reading the message from server:", err)
		return
	}
	log.Printf("Response from WebOS: %s", message)
}

// Initialize websocket on each command call
func connect(addr net.IP, data *LGPair) *websocket.Conn {
	dst := fmt.Sprintf("%s:3000", addr)

	u := url.URL{Scheme: "ws", Host: dst}
	log.Printf("connecting to %s", u.String())
	client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("dial:", err)
	}

	register(client, data)
	return client
}

func main() {
	username := flag.String("username", "username", "REST API username")
	password := flag.String("password", "password", "REST API password")
	flag.Parse()

	pair := &LGPair{}

	// get absolute location
	if dirErr != nil {
		log.Fatal(dirErr)
	}

	// Prepare payload for initial
	data := pair.loadData()

	// Addr of webos device
	addr := addrDiscovery().IP

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		user, pass, _ := req.BasicAuth()
		if user != *username || pass != *password {
			http.Error(w, "Unauthorized.", 401)
			return
		}

		/// Parse JSON's POST request
		decoder := json.NewDecoder(req.Body)

		var p map[string]interface{}
		err := decoder.Decode(&p)
		if err != nil {
			log.Println("Can not decode json, ", err)
			return
		}
		ssap := fmt.Sprintf("ssap://%v", req.URL.Path[1:])
		cmd := &cmd{Type: "request", ID: getID(32), URI: ssap}

		// handle parameters for ( in the order ):
		// - audio/setVolume
		// - system.launcher/launch
		// - audio/setMute
		// - system.notifications/createToast

		for key := range p {
			switch key {
			case "volume":
				if volume, ok := p[key].(float64); ok {
					cmd.Payload = payload{Volume: int(volume)}
				}
			case "id":
				if id, ok := p[key].(string); ok {
					cmd.Payload = payload{ID: strings.ToLower(id)}
				}
			case "mute":
				if mute, ok := p[key].(bool); ok {
					cmd.Payload = payload{Mute: mute}
				}
			case "message":
				if message, ok := p[key].(string); ok {
					cmd.Payload = payload{Message: message}
				}
			}
		}
		command(addr, cmd, &data)
	})

	/*
		$ openssl genrsa -out server.key 2048
		or
		$ openssl ecparam -genkey -name secp384r1 -out server.key

		$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	*/
	err := http.ListenAndServeTLS(":8888", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
