package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"

	"github.com/axxelG/loxwebhook/config"
	"github.com/axxelG/loxwebhook/controls"
	"github.com/axxelG/loxwebhook/helpers"
)

var limiter = rate.NewLimiter(1, 3)

type authKeyError struct {
	err string
}

func (e *authKeyError) Error() string {
	return e.err
}

type commandError struct {
	err string
}

type controlError struct {
	err string
}

func (e *controlError) Error() string {
	return e.err
}

func sendErrorPage(logger *log.Logger, w http.ResponseWriter, err error, responseCode int) {
	logger.Println(err)
	displayText := fmt.Sprintf("%d %s", responseCode, http.StatusText(responseCode))
	w.WriteHeader(responseCode)
	fmt.Fprintln(w, displayText)
	fmt.Fprintf(w, "%s", err)
}

func sendRequest(cfg *config.Config, path string, logger *log.Logger) (*http.Response, error) {
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
	logger.Printf("%s:%s %s", req.RemoteAddr, req.Method, req.URL.Path)
	return resp, nil
}

func forwardResponse(resp *http.Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
	io.Copy(w, resp.Body)
	resp.Body.Close()
}

func getAuthKeyFromRequest(req *http.Request) (string, error) {
	authKeys, ok := req.URL.Query()["t"]
	if !ok {
		return "", &authKeyError{
			err: "Request without access authKey",
		}
	}
	return authKeys[0], nil
}

func authorize(control controls.Control, authKeys map[string]string, reqAuthKey, reqCommand string) error {
	reqAuthKeyKey, ok := helpers.GetMapStringKeyFromStringValue(reqAuthKey, authKeys)
	if !ok {
		return fmt.Errorf("Unknown authKey: %s", reqAuthKey)
	}
	if !helpers.IsStringInSlice(reqAuthKeyKey, control.AuthKeys) {
		return fmt.Errorf("AuthKey %s is not valid for this control", reqAuthKeyKey)
	}
	if !helpers.IsStringInSlice(reqCommand, control.Allowed) {
		return fmt.Errorf("Command %s is not allowed on this control", reqCommand)
	}
	return nil
}

func getControlID(controls map[string]controls.Control, control string) (int, error) {
	ctl, ok := controls[control]
	if !ok {
		return 0, fmt.Errorf("Unknown control %s", control)
	}
	return ctl.ID, nil
}

// StartServer starts the proxy server
func StartServer(
	listener net.Listener,
	tlsConfig *tls.Config,
	cfg *config.Config,
	loggerErr *log.Logger,
	loggerAcc *log.Logger,
	authKeys map[string]string,
	controls map[string]controls.Control,
) error {

	notFoundHandler := func(w http.ResponseWriter, req *http.Request) {
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
		controlName, command, authKey := parseRequestDigitalVirtualInput(req)
		ctl, ok := controls[controlName]
		if !ok {
			err := fmt.Errorf("Unknown control %s", controlName)
			sendErrorPage(loggerErr, w, err, http.StatusNotFound)
			return
		}
		err := authorize(ctl, authKeys, authKey, command)
		if err != nil {
			sendErrorPage(loggerErr, w, err, http.StatusUnauthorized)
			return
		}
		vi, err := newDigitalVirtualInput(ctl.ID, command, authKey)
		if err != nil {
			sendErrorPage(loggerErr, w, err, http.StatusNotFound)
			return
		}
		if err != nil {
			sendErrorPage(loggerErr, w, err, http.StatusUnauthorized)
			return
		}
		if _, ok := req.URL.Query()["simulate"]; ok {
			fmt.Fprintf(w, "SIMULATE\n")
			fmt.Fprintf(w, "Virtual Input: %d\n", vi.ID)
			fmt.Fprintf(w, "Command:       %s\n", vi.Command)
			fmt.Fprintf(w, "AuthKey:         %s\n", vi.AuthKey)
			fmt.Fprintf(w, "Path:          %s\n", vi.GetPath())
			return
		}
		resp, err := sendRequest(cfg, vi.GetPath(), loggerAcc)
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
	router.HandleFunc("/", notFoundHandler)
	for _, control := range controls {
		switch control.Category {
		case "dvi":
			router.HandleFunc("/dvi/{control}/{command}", LoggingHandler(Limiter(DigitalVirtualInputHandler)))
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
