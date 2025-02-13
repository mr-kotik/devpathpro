package config

// Program structure holds information about a development tool
type Program struct {
	Name           string   `json:"name"`
	ExecutableName string   `json:"executableName"`
	CommonPaths    []string `json:"commonPaths"`
	Category       string   `json:"category"`
}

// Configuration holds the global configuration
type Configuration struct {
	Programs []Program
	LogFile  string
} 