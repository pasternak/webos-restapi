package main

// Register - used for initial registration with the WebOS
type Register struct {
	Type    string  `json:"type"`
	PayLoad *LGPair `json:"payload"`
}

// LGPair - json used to establish pairing with LG TV
type LGPair struct {
	ClientKey    string       `json:"client-key,omitempty"`
	ForcePairing bool         `json:"forcePairing"`
	PairingType  string       `json:"pairingType"`
	Manifest     manifestData `json:"manifest"`
}

type manifestData struct {
	ManifestVersion int         `json:"manifestVersion"`
	AppVersion      string      `json:"appVersion"`
	Permissions     []string    `json:"permissions"`
	Signed          signData    `json:"signed"`
	Signatures      []signature `json:"signatures"`
}

type signData struct {
	Created              string            `json:"created"`
	AppID                string            `json:"appId"`
	VendorID             string            `json:"vendorId"`
	Permissions          []string          `json:"permissions"`
	Serial               string            `json:"serial"`
	LocalizedAppNames    map[string]string `json:"localizedAppNames"`
	LocalizedVendorNames map[string]string `json:"localizedVendorNames"`
}

type signature struct {
	SignatureVersion int    `json:"signatureVersion"`
	Signature        string `json:"signature"`
}

// WebosHandshake - prepare authorization call for each request
func (data *LGPair) WebosHandshake() {
	// Prepare payload for initial

	data.Manifest = manifestData{
		ManifestVersion: 1,
		AppVersion:      "1.1",
		Permissions: []string{
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
		},
		Signed: signData{
			Created:  "20140509",
			AppID:    "com.lge.test",
			VendorID: "com.lge",
			Permissions: []string{
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
			},
			Serial: "2f930e2d2cfe083771f68e4fe7bb07",
			LocalizedAppNames: map[string]string{
				"":       "LG Remote App",
				"ko-KR":  "리모컨 앱",
				"zxx-XX": "ЛГ Rэмotэ AПП",
			},
			LocalizedVendorNames: map[string]string{
				"": "LG Electronics",
			},
		},
		Signatures: []signature{
			signature{
				Signature:        "eyJhbGdvcml0aG0iOiJSU0EtU0hBMjU2Iiwia2V5SWQiOiJ0ZXN0LXNpZ25pbmctY2VydCIsInNpZ25hdHVyZVZlcnNpb24iOjF9.hrVRgjCwXVvE2OOSpDZ58hR+59aFNwYDyjQgKk3auukd7pcegmE2CzPCa0bJ0ZsRAcKkCTJrWo5iDzNhMBWRyaMOv5zWSrthlf7G128qvIlpMT0YNY+n/FaOHE73uLrS/g7swl3/qH/BGFG2Hu4RlL48eb3lLKqTt2xKHdCs6Cd4RMfJPYnzgvI4BNrFUKsjkcu+WD4OO2A27Pq1n50cMchmcaXadJhGrOqH5YmHdOCj5NSHzJYrsW0HPlpuAx/ECMeIZYDh6RMqaFM2DXzdKX9NmmyqzJ3o/0lkk/N97gfVRLW5hA29yeAwaCViZNCP8iC9aO0q9fQojoa7NQnAtw==",
				SignatureVersion: 1,
			},
		},
	}
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

type request struct {
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
