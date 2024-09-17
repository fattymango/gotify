package server

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start() // For Linux
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start() // For Windows
	case "darwin":
		err = exec.Command("open", url).Start() // For macOS
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal("Failed to open browser:", err)
	}
}

func (s *Server) getAuthUrl() string {
	return s.Auth.SpotifyAuth.AuthURL(s.config.Spotify.State)
}
