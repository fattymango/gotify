package server

import (
	"context"
	"fmt"
	"gotify/config"
	"gotify/pkg/auth"
	httppkg "gotify/pkg/http"
	"gotify/pkg/logger"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

type Server struct {
	config    *config.Config
	http      *http.Server
	logger    *logger.Logger
	Auth      *auth.Auth
	ch        chan<- *oauth2.Token
	token     *oauth2.Token
	isRunning bool
}

// takes a config, logger, and a channel to send the token to, after auth is complete
func NewServer(cfg *config.Config, l *logger.Logger, a *auth.Auth) *Server {
	return &Server{
		config: cfg,
		http:   httppkg.NewHttp(cfg),
		logger: l,
		Auth:   a,
	}
}

// will start the server and listen for requests, will block the main thread until the server is stopped
func (s *Server) Start(ch chan<- *oauth2.Token) error {
	s.ch = ch

	// if s.Auth == nil {
	// 	s.Auth = auth.NewAuth(s.config)
	// }

	s.RegisterRoutes()

	go func() {
		s.logger.Logger.Infof("Server started on %s", s.config.Server.Address)
		s.isRunning = true

		if err := s.http.ListenAndServe(); err != nil {
			// return fmt.Errorf("failed to start server: %w", err)
			s.isRunning = false
			s.logger.Logger.Errorf("Server stopped: %s", err)
			return
		}
	}()

	url := s.getAuthUrl()
	fmt.Println("Opening browser to authenticate with Spotify...")
	fmt.Println("If the browser does not open, please visit the following URL:\n", url)
	openBrowser(url)

	return nil
}

func (s *Server) Stop() {
	if !s.isRunning {
		return
	}

	s.isRunning = false
	close(s.ch)

	if err := s.http.Shutdown(context.TODO()); err != nil {
		log.Fatalf("Failed to shutdown server: %s", err)
	}

	s.logger.Logger.Info("Server stopped")
}

func (s *Server) RegisterRoutes() {
	// Register routes here
	http.HandleFunc("GET /callback", s.completeAuth)
	http.HandleFunc("GET /complete", s.authComplete)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.logger.Logger.Infof("Got request for: %s", r.URL.String())
	})
}

// show the user that the auth is complete in a nice way
func (s *Server) authComplete(w http.ResponseWriter, r *http.Request) {
	NotAuthenticatedMsg := "You are not authenticated with Spotify. Please authenticate by visiting " + s.getAuthUrl()
	AuthenticatedMsg := "You are authenticated with Spotify. You can close this window now."

	if s.token == nil {
		w.Write([]byte(NotAuthenticatedMsg))
	} else {
		w.Write([]byte(AuthenticatedMsg))
	}

}

func (s *Server) completeAuth(w http.ResponseWriter, r *http.Request) {
	s.logger.Logger.Info("Completing Auth...")
	state := s.config.Spotify.State
	tok, err := s.Auth.SpotifyAuth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// send the token to the channel
	s.logger.Logger.Info("Login Completed!")
	s.token = tok
	s.ch <- tok

	go func() {
		time.Sleep(5 * time.Second)
		s.Stop()
	}()

	http.Redirect(w, r, "/complete", http.StatusSeeOther)
}
