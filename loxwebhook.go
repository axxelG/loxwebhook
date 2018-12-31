package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/axxelG/crypto/acme/autocert"
	"github.com/coreos/go-systemd/daemon"
	"github.com/pkg/errors"

	"github.com/axxelG/loxwebhook/config"
	"github.com/axxelG/loxwebhook/controls"
	"github.com/axxelG/loxwebhook/proxy"
)

var version string // Will be set on compile time

func initLogging(FileName string) (*log.Logger, *os.File, error) {
	logFormat := log.Ldate | log.Ltime | log.Lshortfile
	if FileName == "" {
		logger := log.New(os.Stderr, "", logFormat)
		return logger, nil, nil
	}
	logFile, err := os.OpenFile(FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error opening logfile")
	}
	logger := log.New(logFile, "", logFormat)
	return logger, logFile, nil
}

func startLetsEncryptListener(cfg *config.Config) (net.Listener, *tls.Config) {
	m := &autocert.Manager{
		Cache:      autocert.DirCache(cfg.LetsEncryptCache),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfg.PublicURI),
	}
	return m.ListenerCustomAddress(cfg.GetListenPort()), m.TLSConfig()
}

func main() {
	cfg, err := config.NewConfig(version)
	if err != nil {
		log.Print(errors.Wrap(err, "Cannot read/load config"))
		os.Exit(1)
	}
	loggerMain, logFileMain, err := initLogging(cfg.LogFileMain)
	if err != nil {
		log.Print(errors.Wrap(err, "Cannot write logfile"))
		os.Exit(1)
	}
	defer logFileMain.Close()

	loggerMain.Println("Starting loxwebhook")
	loggerMain.Print(cfg.String())

	tokens, controls, err := controls.Read(cfg.ControlsFiles)
	if err != nil {
		log.Print(errors.Wrap(err, "Error importing controls"))
		os.Exit(1)
	}

	LoggerHTTPErrors, logFileHTTPErrors, err := initLogging(cfg.LogFileHTTPError)
	if err != nil {
		log.Print(errors.Wrap(err, "Cannot write logfile http errors"))
		os.Exit(1)
	}
	defer logFileHTTPErrors.Close()

	LoggerHTTPAccess, LogFileHTTPAccess, err := initLogging(cfg.LogFileHTTPAccess)
	if err != nil {
		log.Print(errors.Wrap(err, "Cannot write logfile http access"))
		os.Exit(1)
	}
	defer LogFileHTTPAccess.Close()

	listener, tlsConfig := startLetsEncryptListener(cfg)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	err = proxy.StartServer(listener, tlsConfig, cfg, LoggerHTTPErrors, LoggerHTTPAccess, tokens, controls)
	if err != nil {
		log.Print(errors.Wrap(err, "Error starting server"))
		os.Exit(1)
	}
}
