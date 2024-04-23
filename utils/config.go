package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"os/user"
	"runtime"
	"strings"
)

func defaultColors() map[string]string {
	colors := map[string]string{}
	colors["emergency"] = "\033[1;101m" // Emergency: system is unusable
	colors["alert"] = "\033[31m"        // Alert: action must be taken immediately
	colors["critical"] = "\033[1;31m"   // Critical: critical conditions
	colors["error"] = "\033[31m"        // Error: error conditions
	colors["warning"] = "\033[33m"      // Warning: warning conditions
	colors["notice"] = "\033[1;33m"     // Notice: normal but significant condition
	colors["info"] = "\033[36m"         // Informational: informational messages
	colors["debug"] = "\033[0;32m"      // Debug: debug messages

	return colors
}

type Config struct {
	colors       map[string]string // INFO: Colors mapping
	defaultColor string
}

func NewConfig() *Config {
	return &Config{
		colors:       defaultColors(),
		defaultColor: "\033[0m",
	}
}

func GetConfigPath(dir ...string) string {
	osType := runtime.GOOS
	n := map[string]string{"darwin": ".config", "windows": "AppData\\Roaming", "linux": ".config"}
	confDir := n[osType]

	USER, _ := user.Current()
	newPath := []string{USER.HomeDir, confDir}
	for i := 0; i < len(dir); i++ {
		newPath = append(newPath, dir[i])
	}

	return strings.Join(newPath, string(os.PathSeparator))
}

var config *Config

type Color struct {
	err []string `ini:"error"`
}

func readConfiguration(path string) {

	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	var color *Color

	err = cfg.Section("Colors").MapTo(color)
	if err != nil {
		return
	}

	fmt.Println(color)

	// fmt.Println(fmt.Sprintf("%stest%s",color.Colors["error"], config.defaultColor))
}

func ParseConfigurations() {
	config = NewConfig()
	path := GetConfigPath("logger", "config.ini")
	if _, err := os.Stat(path); err == nil {
		// readConfiguration(path)
	}

}
