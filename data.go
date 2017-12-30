package main

// LGPair - json used to establish pairing with LG TV
type LGPair struct {
	ForcePairing bool   `json:"forcePairing"`
	PairingType  string `json:"pairingType"`
	Manifest     struct {
		Signed struct {
			Permissions []string `json:"permissions"`
		} `json:"signed"`
		Permissions []string `json:"permissions"`
	} `json:"manifest"`
	ClientKey string `json:"client-key,omitempty"`
}

// Register - used for initial registration with the WebOS
type Register struct {
	Type    string  `json:"type"`
	PayLoad *LGPair `json:"payload"`
}

// Receiver - keep messages received back from LG websocket
type Receiver struct {
	Type    string `json:"type"`
	ID      string `json:"id,omitempty"`
	Error   string `json:"error,omitempty"`
	Payload struct {
		Subscribed   bool `json:"subscribed,omitempty"`
		LaunchPoints []struct {
			SystemApp      bool   `json:"systemApp"`
			ID             string `json:"id"`
			Title          string `json:"title"`
			AppDescription string `json:"appDescription"`
		} `json:"launchPoints,omitempty"`
		ClientKey   string `json:"client-key,omitempty"`
		ReturnValue bool   `json:"returnValue"`
	} `json:"payload,omitempty"`
}

type cmd struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	URI     string `json:"uri"`
	Payload struct {
		Volume  int    `json:"volume,omitempty"`
		Message string `json:"message,omitempty"`
		ID      string `json:"id,omitempty"`
		Mute    bool   `json:"mute,omitempty"`
	} `json:"payload,omitempty"`
}
