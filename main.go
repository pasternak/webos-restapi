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
	"strings"
	"time"

	"github.com/gorilla/websocket"
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

func command(addr net.IP, cmd *cmd, data *LGPair) []byte {
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
	}
	log.Printf("Response from WebOS: %s", message)
	return message
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

	// Prepare payload for initial
	data := LGPair{
		ForcePairing: false,
		PairingType:  "PROMPT",
	}
	data.Manifest.Signed.Permissions = []string{
		"TEST_SECURE",
		"CONTROL_INPUT_TEXT",
		"CONTROL_MOUSE_AND_KEYBOARD",
		"READ_INSTALLED_APPS",
		"READ_LGE_SDX",
		"READ_NOTIFICATIONS",
		"SEARCH",
		"WRITE_SETTINGS",
		"WRITE_NOTIFICATION_ALERT",
		"CONTROL_POWER",
		"READ_CURRENT_CHANNEL",
		"READ_RUNNING_APPS",
		"READ_UPDATE_INFO",
		"UPDATE_FROM_REMOTE_APP",
		"READ_LGE_TV_INPUT_EVENTS",
		"READ_TV_CURRENT_TIME",
	}
	data.Manifest.Permissions = []string{
		"LAUNCH",
		"LAUNCH_WEBAPP",
		"APP_TO_APP",
		"CLOSE",
		"TEST_OPEN",
		"TEST_PROTECTED",
		"CONTROL_AUDIO",
		"CONTROL_DISPLAY",
		"CONTROL_INPUT_JOYSTICK",
		"CONTROL_INPUT_MEDIA_RECORDING",
		"CONTROL_INPUT_MEDIA_PLAYBACK",
		"CONTROL_INPUT_TV",
		"CONTROL_POWER",
		"READ_APP_STATUS",
		"READ_CURRENT_CHANNEL",
		"READ_INPUT_DEVICE_LIST",
		"READ_NETWORK_STATE",
		"READ_RUNNING_APPS",
		"READ_TV_CHANNEL_LIST",
		"WRITE_NOTIFICATION_TOAST",
		"READ_POWER_STATE",
		"READ_COUNTRY_INFO",
	}

	// Addr of webos device
	addr := webosDiscovery()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		user, pass, _ := req.BasicAuth()
		if user != *username || pass != *password {
			http.Error(w, "Unauthorized.", 401)
			return
		}

		/// Parse JSON's POST request
		decoder := json.NewDecoder(req.Body)

		var post map[string]interface{}
		err := decoder.Decode(&post)
		if err != nil {
			log.Println("Can not decode json, ", err)
			return
		}
		ssap := fmt.Sprintf("ssap://%v", req.URL.Path[1:])
		cmd := &cmd{Type: "request", ID: getID(32), URI: ssap}

		for key := range post {
			switch key {
			case "volume":
				if volume, ok := post[key].(float64); ok {
					cmd.Payload.Volume = int(volume)
				}
			case "id":
				if id, ok := post[key].(string); ok {
					cmd.URI = "ssap://com.webos.applicationManager/listLaunchPoints"
					apps := &Receiver{}
					json.Unmarshal(command(addr.IP, cmd, &data), apps)
					for _, app := range apps.Payload.LaunchPoints {
						if strings.Contains(strings.ToLower(app.Title), strings.ToLower(id)) {
							cmd.Payload.ID = app.ID
							cmd.URI = ssap
							break
						}
					}
				}
			case "mute":
				if mute, ok := post[key].(bool); ok {
					cmd.Payload.Mute = mute
				}
			case "message":
				if message, ok := post[key].(string); ok {
					cmd.Payload.Message = message
				}
			}
		}
		command(addr.IP, cmd, &data)
	})

	err := http.ListenAndServeTLS(":8888", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
