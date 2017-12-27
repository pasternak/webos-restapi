package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// LGPair - json used to establish pairing with LG TV
type LGPair struct {
	ForcePairing bool   `json:"forcePairing"`
	PairingType  string `json:"pairingType"`
	ClientKey    string `json:"client-key"`
	Manifest     `json:"manifest"`
}

// Manifest -
type Manifest struct {
	ManifestVersion uint64       `json:"manifestVersion"`
	AppVersion      string       `json:"appVersion"`
	Permissions     []string     `json:"permissions"`
	Signatures      []Signatures `json:"signatures"`
	Signed          `json:"signed"`
}

// Signed -
type Signed struct {
	Created              string            `json:"created"`
	AppID                string            `json:"appId"`
	VendorID             string            `json:"vendorId"`
	LocalizedAppNames    map[string]string `json:"localizedAppNames"`
	LocalizedVendorNames map[string]string `json:"localizedVendorNames"`
	Permissions          []string          `json:"permissions"`
	Serial               string            `json:"serial"`
}

// Signatures -
type Signatures struct {
	SignatureVersion uint64 `json:"signatureVersion"`
	Signature        string `json:"signature"`
}

// Register - used for initial registration with the WebOS
type Register struct {
	Type    string  `json:"type"`
	PayLoad *LGPair `json:"payload"`
}

// Receiver - keep messages received back from LG websocket
type Receiver struct {
	Type    string                 `json:"type"`
	Error   string                 `json:"error,omitempty"`
	PayLoad map[string]interface{} `json:"payload,omitempy"`
}

func register(client *websocket.Conn, data *LGPair) {
	register := Register{
		Type:    "register",
		PayLoad: data,
	}
	receiver := Receiver{}

	message, err := json.Marshal(register)
	if err != nil {
		log.Fatal(err.Error())
	}

	if data.ClientKey == "" {
		usr, _ := user.Current()
		webosFile := path.Join(usr.HomeDir, ".webos")
		dat, err := ioutil.ReadFile(webosFile)
		if err != nil {
			log.Println("Requesting pairing with WebOS")
			client.WriteMessage(websocket.TextMessage, message)

			// Wait till user approve initial pairing
			for {
				_, message, err = client.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					break
				}
				log.Println("log: ", string(message))
				json.Unmarshal(message, &receiver)
				if receiver.Type == "registered" {
					ioutil.WriteFile(
						path.Join(webosFile),
						[]byte(receiver.PayLoad["client-key"].(string)),
						0644)
					log.Println("Paired")
				}
			}
		} else {
			data.ClientKey = string(dat)
			message, err = json.Marshal(Register{
				Type:    "register",
				PayLoad: data,
			})
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
	// after the session is initialized send the payload to request registration
	err = client.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("register failed on write:", err)
	}
	_, message, err = client.ReadMessage()
	if err != nil {
		log.Println("register failed on response:", err)
	}
}

func addrDiscovery() (addr *net.UDPAddr) {
	destination := &net.UDPAddr{
		IP:   net.ParseIP("239.255.255.250"),
		Port: 1900,
	}

	/*
		ssdp := []byte{
			77, 45, 83, 69, 65, 82, 67, 72, 32, 42,
			32, 72, 84, 84, 80, 47, 49, 46, 49, 13,
			10, 72, 79, 83, 84, 58, 32, 50, 51, 57,
			46, 50, 53, 53, 46, 50, 53, 53, 46, 50,
			53, 48, 58, 49, 57, 48, 48, 13, 10, 77,
			65, 78, 58, 32, 34, 115, 115, 100, 112,
			58, 100, 105, 115, 99, 111, 118, 101, 114,
			34, 13, 10, 77, 88, 58, 32, 50, 13, 10,
			83, 84, 58, 32, 117, 114, 110, 58, 100,
			105, 97, 108, 45, 109, 117, 108, 116, 105,
			115, 99, 114, 101, 101, 110, 45, 111, 114,
			103, 58, 115, 101, 114, 118, 105, 99, 101,
			58, 100, 105, 97, 108, 58, 49, 13, 10, 13, 10}*/

	ssdp := []string{
		"M-SEARCH * HTTP/1.1\r\n",
		"HOST: 239.255.255.250:1900\r\n",
		"MAN: \"ssdp:discover\"\r\n",
		"MX: 2\r\n",
		"ST: urn:dial-multiscreen-org:service:dial:1\r\n",
		"\r\n",
	}

	server, err := net.ListenUDP("udp4", &net.UDPAddr{})
	if err != nil {
		log.Fatal("Can not bind")
	}
	defer server.Close()

	// send dicovery message on the broadcast addr
	server.WriteToUDP([]byte(strings.Join(ssdp, "")), destination)

	// Read the messages back from the broadcast, till WebOS is received
	buffer := make([]byte, 1024)
	server.SetReadDeadline(time.Now().Add(30 * time.Second))
	for {
		_, addr, err := server.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("LG WebOS not found.")
		}
		if len(buffer) > 1 {
			if strings.Contains(string(buffer[:]), "WebOS") {
				log.Println("found WebOS: ", addr.IP)
				return addr
			}
		}
	}
}

func (pair LGPair) loadData() LGPair {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	raw, err := ioutil.ReadFile(path.Join(dir, "./pairing.json"))
	if err != nil {
		log.Fatal(err.Error())
	}
	json.Unmarshal(raw, &pair)
	return pair
}

func (pair LGPair) toJSON() string {
	bytes, err := json.Marshal(pair.loadData())
	if err != nil {
		log.Fatal(err.Error())
	}

	return string(bytes)
}
