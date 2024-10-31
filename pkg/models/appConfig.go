package models

type Config struct {
	Application struct {
		LogPath  string `json:"logPath"`
		DataPath string `json:"dataPath"`
	} `json:"application"`
}
