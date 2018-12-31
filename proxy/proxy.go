package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/axxelG/loxwebhook/config"
	"github.com/axxelG/loxwebhook/controls"
)

var limiter = rate.NewLimiter(1, 3)

type tokenError struct {
	err string
}

func (e *tokenError) Error() string {
	return e.err
}

type commandError struct {
	err string
}

func (e *commandError) Error() string {
	return e.err
}

type controlError struct {
	err string
}

func (e *controlError) Error() string {
	return e.err
}

// func sendErrorPage(w http.ResponseWriter, text string, responseCode int) {
// 	w.WriteHeader(responseCode)
// 	fmt.Fprintf(w, text)
// }

func sendErrorPage(logger *log.Logger, w http.ResponseWriter, err error, responseCode int) {
	logger.Println(err)
	displayText := fmt.Sprintf("%d %s", responseCode, http.StatusText(responseCode))
	w.WriteHeader(responseCode)
	fmt.Fprintln(w, displayText)
	fmt.Fprintf(w, "%s", err)
}

func sendRequest(cfg *config.Config, path string) (*http.Response, error) {
	url := cfg.MiniserverURL
	url.Path = path
	client := http.Client{
		Timeout: cfg.MiniserverTimeout,
	}
	req, err := http.NewRequest("POST", url.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error preparing request")
	}
	req.SetBasicAuth(cfg.MiniserverUser, cfg.MiniserverPassword)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Error sending request to Miniserver")
	}
	return resp, nil
}

func forwardResponse(resp *http.Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func getTokenFromRequest(req *http.Request) (string, error) {
	tokens, ok := req.URL.Query()["t"]
	if !ok {
		return "", &tokenError{
			err: "Request without access token",
		}
	}
	return tokens[0], nil
}

func getValidControls(controls map[string]controls.Control) map[int]bool {
	validControls := map[int]bool{}
	for _, control := range controls {
		validControls[control.ID] = true
	}
	return validControls
}

func getValidTokens(tokens map[string]string, controls map[string]controls.Control, req *http.Request) (map[string]bool, error) {
	parts := strings.Split(strings.ToLower(req.URL.Path), "/")
	category := parts[1]
	ID, err := strconv.Atoi(parts[2])
	if err != nil {
		err = &controlError{
			err: errors.Wrap(err, "Error converting control ID").Error(),
		}
		return nil, err
	}
	command := parts[3]
	validTokenNames := []string{}
	for _, control := range controls {
		if control.Category == category && control.ID == ID {
			for _, allowedCommand := range control.Allowed {
				if allowedCommand == "<all>" || allowedCommand == command {
					validTokenNames = append(validTokenNames, control.Tokens...)
				}
			}
		}
	}
	validTokens := map[string]bool{}
	for _, name := range validTokenNames {
		validTokens[tokens[name]] = true
	}
	return validTokens, nil
}

func autorize(tokens map[string]string, controls map[string]controls.Control, req *http.Request) error {
	validTokens, err := getValidTokens(tokens, controls, req)
	if err != nil {
		return err
	}
	givenToken, err := getTokenFromRequest(req)
	if err != nil {
		return err
	}
	if validTokens[givenToken] {
		return nil
	}
	return &tokenError{
		err: "Token not valid for this request",
	}
}

// StartServer starts the proxy server
func StartServer(listener net.Listener, tlsConfig *tls.Config, cfg *config.Config, loggerErr *log.Logger, loggerAcc *log.Logger, tokens map[string]string, controls map[string]controls.Control) error {
	RootHandler := func(w http.ResponseWriter, req *http.Request) {
		http.NotFound(w, req)
	}

	Limiter := func(nextHandler http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() == false {
				err := errors.New("Request rate limit reached")
				sendErrorPage(loggerErr, w, err, http.StatusTooManyRequests)
				return
			}

			nextHandler.ServeHTTP(w, r)
		})
	}

	LoggingHandler := func(nextHandler http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			loggerAcc.Printf("%s:%s %s", req.RemoteAddr, req.Method, req.URL.Path)
			nextHandler.ServeHTTP(w, req)
		})
	}

	DigitalVirtualInputHandler := func(w http.ResponseWriter, req *http.Request) {
		vi, err := newDigitalVirtualInput(req)
		if err != nil {
			if err, ok := err.(*commandError); ok {
				sendErrorPage(loggerErr, w, err, http.StatusNotFound)
				return
			}
			sendErrorPage(loggerErr, w, err, http.StatusBadRequest)
			return
		}
		validControls := getValidControls(controls)
		if !validControls[vi.ID] {
			err = &controlError{
				err: fmt.Sprintf("Control not found: %d", vi.ID),
			}
			sendErrorPage(loggerErr, w, err, http.StatusNotFound)
			return
		}
		err = autorize(tokens, controls, req)
		if err != nil {
			sendErrorPage(loggerErr, w, err, http.StatusUnauthorized)
			return
		}
		if _, ok := req.URL.Query()["simulate"]; ok {
			fmt.Fprintf(w, "SIMULATE\n")
			fmt.Fprintf(w, "Virtual Input: %d\n", vi.ID)
			fmt.Fprintf(w, "Command:       %s\n", vi.Command)
			fmt.Fprintf(w, "Token:         %s\n", vi.Token)
			fmt.Fprintf(w, "Path:          %s\n", vi.GetPath())
			return
		}
		resp, err := sendRequest(cfg, vi.GetPath())
		if err != nil {
			code := http.StatusBadGateway
			if e, ok := err.(*url.Error); ok {
				if e.Timeout() {
					code = http.StatusGatewayTimeout
				}
			}
			sendErrorPage(loggerErr, w, err, code)
			return
		}
		forwardResponse(resp, w)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", RootHandler)
	for _, control := range controls {
		switch control.Category {
		case "dvi":
			router.HandleFunc("/dvi/{id:[0-9]+}/{command}", LoggingHandler(Limiter(DigitalVirtualInputHandler)))
		}
	}
	s := &http.Server{
		TLSConfig:   tlsConfig,
		Handler:     router,
		ReadTimeout: cfg.MiniserverTimeout,
		ErrorLog:    loggerErr,
	}
	err := s.Serve(listener)
	if err != nil {
		return errors.Wrap(err, "Error starting server")
	}

	return nil
}
