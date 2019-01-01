package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Config holds the configuration values
type Config struct {
	Version            string
	ConfigFile         string
	LogFileMain        string
	LogFileHTTPError   string
	LogFileHTTPAccess  string
	ListenPort         int
	PublicURI          string
	LetsEncryptCache   string
	ControlsFiles      string
	MiniserverURL      *url.URL
	MiniserverUser     string
	MiniserverPassword string
	MiniserverTimeout  time.Duration
}

// String returns a multiline String to print Config.
func (c *Config) String() string {
	return fmt.Sprintf(
		"Config:\n"+
			"Version:              %s\n"+
			"Config file:          %s\n"+
			"Log file main:        %s\n"+
			"Log file http errors: %s\n"+
			"Log file http access: %s\n"+
			"Listen Port:          %d\n"+
			"Public URI:           %s\n"+
			"LetsEncrypt Cache:    %s\n"+
			"Configs Directory:    %s\n"+
			"Miniserver URL:       %s\n"+
			"Miniserver User:      %s\n"+
			"Miniserver Timeout:   %d seconds\n",
		c.Version,
		c.ConfigFile,
		c.LogFileMain,
		c.LogFileHTTPError,
		c.LogFileHTTPAccess,
		c.ListenPort,
		c.PublicURI,
		c.LetsEncryptCache,
		c.ControlsFiles,
		c.MiniserverURL,
		c.MiniserverUser,
		int64(c.MiniserverTimeout.Seconds()),
	)
}

func (c *Config) checkFile(fn, description string) error {
	if fn == "" {
		return nil
	}
	_, err := os.Stat(fn)
	if err != nil {
		return errors.Wrap(err, "Error testing access to "+description)
	}
	return nil
}

func (c *Config) tomlCount(dir string) error {
	counter := 0
	deadline := time.Now().Add(2 * time.Second)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".toml" {
			counter++
		}
		if time.Now().After(deadline) {
			return errors.New("Searching for control files took too long. Maybe you provided a huge directory structure as control files path")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if counter == 0 {
		return errors.New("No toml file found in controls dir")
	}
	return nil
}

func (c *Config) checkHostname(h string) error {
	// Regex copied from https://socketloop.com/tutorials/golang-validate-hostname
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if !re.MatchString(h) {
		if strings.HasPrefix(strings.ToLower(h), "http://") || strings.HasPrefix(strings.ToLower(h), "https://") {
			return errors.New("Public URI must not start with http:// or https://")
		}
		return errors.New("Invalid public URI " + h)
	}
	return nil
}

func (c *Config) reachMiniserver(address *url.URL) error {
	testEndpoint := address
	testEndpoint.Path = "/jdev/cfg/api"
	client := http.Client{
		Timeout: c.MiniserverTimeout,
	}
	resp, err := client.Get(testEndpoint.String())
	if err != nil {
		return errors.Wrap(err, "Cannot reach miniserver")
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Miniserver responded with status code %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}

// Validate returns an error if the validation of config values failed
func (c *Config) Validate() error {
	if err := c.checkFile(c.ConfigFile, "config file"); err != nil {
		return err
	}
	if c.ListenPort < 1 {
		return errors.New("ListenPort must be >= 1")
	}
	if c.ListenPort >= 65535 {
		// We are using port 65535 as default value for the listenport flag.
		return errors.New("ListenPort must be < 65535")
	}
	if err := c.checkHostname(c.PublicURI); err != nil {
		return err
	}
	if _, err := net.LookupIP(c.PublicURI); err != nil {
		return errors.Wrap(err, "Error looking up public URI")
	}
	if err := c.checkFile(c.ControlsFiles, "control files dir"); err != nil {
		return err
	}
	if err := c.tomlCount(c.ControlsFiles); err != nil {
		return err
	}
	if err := c.reachMiniserver(c.MiniserverURL); err != nil {
		return err
	}
	//TODO: Validate username and password
	return nil
}

// GetListenPort return a string usable by http.ListenAndServe
func (c Config) GetListenPort() string {
	return ":" + strconv.Itoa(c.ListenPort)
}

// basicTypeConfig holds the config internally used by the
// config package. We are not using types like *url.URL or
// time.Duration because they are not supported by flags,
// environment variables or toml.
type basicTypeConfig struct {
	ConfigFile         string
	LogFileMain        string
	LogFileHTTPError   string
	LogFileHTTPAccess  string
	ListenPort         int
	PublicURI          string
	LetsEncryptCache   string
	ControlsFiles      string
	MiniserverURL      string
	MiniserverUser     string
	MiniserverPassword string
	MiniserverTimeout  int // Seconds
}

func (btc *basicTypeConfig) getConfig() (*Config, error) {
	var err error
	cfg := new(Config)
	cfg.ConfigFile = btc.ConfigFile
	cfg.LogFileMain = btc.LogFileMain
	cfg.LogFileHTTPError = btc.LogFileHTTPError
	cfg.LogFileHTTPAccess = btc.LogFileHTTPAccess
	cfg.ControlsFiles = btc.ControlsFiles
	cfg.ListenPort = btc.ListenPort
	cfg.PublicURI = btc.PublicURI
	cfg.LetsEncryptCache = btc.LetsEncryptCache
	cfg.ControlsFiles = btc.ControlsFiles
	cfg.MiniserverURL, err = url.Parse(btc.MiniserverURL)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing MiniserverURL")
	}
	cfg.MiniserverUser = btc.MiniserverUser
	cfg.MiniserverPassword = btc.MiniserverPassword
	cfg.MiniserverTimeout = time.Duration(btc.MiniserverTimeout) * time.Second
	return cfg, nil
}

// newDefaultConfig returns importConfig struct with default values
func newDefaultConfig() *basicTypeConfig {
	cfg := new(basicTypeConfig)
	cfg.ConfigFile = ""
	cfg.LogFileMain = ""
	cfg.LogFileHTTPError = ""
	cfg.LogFileHTTPAccess = ""
	cfg.ControlsFiles = "./controls.d"
	cfg.ListenPort = 80
	cfg.PublicURI = ""
	cfg.LetsEncryptCache = "./cache/letsencrypt"
	cfg.MiniserverURL = ""
	cfg.MiniserverUser = "admin"
	cfg.MiniserverPassword = "admin"
	cfg.MiniserverTimeout = 2 // Seconds
	return cfg
}

func newEnvConfig() (*basicTypeConfig, error) {
	pref := "LOXWEBHOOK_"
	cfg := newDefaultConfig()
	if val, ok := os.LookupEnv(pref + "CONFIG"); ok {
		cfg.ConfigFile = val
	}
	if val, ok := os.LookupEnv(pref + "LOGFILEMAIN"); ok {
		cfg.LogFileMain = val
	}
	if val, ok := os.LookupEnv(pref + "LOGFILEHTTPERROR"); ok {
		cfg.LogFileHTTPError = val
	}
	if val, ok := os.LookupEnv(pref + "LOGFILEHTTPACCESS"); ok {
		cfg.LogFileHTTPAccess = val
	}
	if val, ok := os.LookupEnv(pref + "CONTROLSFILES"); ok {
		cfg.ControlsFiles = val
	}
	if val, ok := os.LookupEnv(pref + "LISTENPORT"); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting LISTENPORT from env")
		}
		cfg.ListenPort = v
	}
	if val, ok := os.LookupEnv(pref + "PUBLICURI"); ok {
		cfg.PublicURI = val
	}
	if val, ok := os.LookupEnv(pref + "LETSENCRYPTCACHE"); ok {
		cfg.LetsEncryptCache = val
	}
	if val, ok := os.LookupEnv(pref + "MINISERVERURL"); ok {
		cfg.MiniserverURL = val
	}
	if val, ok := os.LookupEnv(pref + "MINISERVERUSER"); ok {
		cfg.MiniserverUser = val
	}
	if val, ok := os.LookupEnv(pref + "MINISERVERPASSWORD"); ok {
		cfg.MiniserverPassword = val
	}
	if val, ok := os.LookupEnv(pref + "MINISERVERTIMEOUT"); ok {
		v, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.Wrap(err, "Error converting MINISERVERTIMEOUT from env")
		}
		cfg.MiniserverTimeout = v
	}
	return cfg, nil
}

func newFlagConfig(versionStr string) *basicTypeConfig {
	// We are using a separate FlagSet to be able reassign values
	// to os.Args in tests without getting errors.
	//TODO: Better flag handling
	// - make flags case insensitive
	// - enable short flags (-v)
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	versionFlag := flags.Bool("version", false, "Show program version")
	config := flags.String("config", "", "Config file")
	logFileMain := flags.String("logfilemain", "", "Log file")
	logFileHTTPError := flags.String("logfilehttperror", "", "Log file")
	logFileHTTPAccess := flags.String("logfilehttpaccess", "", "Log file")
	controlsFiles := flags.String("controlsfiles", "", "Directory containing controls files")
	listenPort := flags.Int("listenport", 65535, "Port to listen on")
	publicURI := flags.String("publicURI", "", "URI where this service is reachable like myhome.example.com")
	letsencryptCache := flags.String("letsencryptCache", "", "Folder where letsencrypt can store cached data")
	miniserverURL := flags.String("miniserverURL", "", "Miniserver URL like http://192.168.1.2:80")
	miniserverUser := flags.String("miniserverUser", "", "Miniserver user")
	miniserverPassword := flags.String("miniserverPassword", "", "Miniserver password")
	miniserverTimeout := flags.Int("miniserverTimeout", 0, "Timeout for requests to the Miniserver")
	flags.Parse(os.Args[1:])
	if *versionFlag {
		fmt.Printf("Version  : %s\n", versionStr)
		fmt.Printf("Build for: %s\n", buildForOS)
		os.Exit(0)
	}
	cfg := newDefaultConfig()
	if *config != "" {
		cfg.ConfigFile = *config
	}
	if *logFileMain != "" {
		cfg.LogFileMain = *logFileMain
	}
	if *logFileHTTPError != "" {
		cfg.LogFileHTTPError = *logFileHTTPError
	}
	if *logFileHTTPAccess != "" {
		cfg.LogFileHTTPAccess = *logFileHTTPAccess
	}
	if *controlsFiles != "" {
		cfg.ControlsFiles = *controlsFiles
	}
	if *listenPort != 65535 {
		cfg.ListenPort = *listenPort
	}
	if *publicURI != "" {
		cfg.PublicURI = *publicURI
	}
	if *letsencryptCache != "" {
		cfg.LetsEncryptCache = *letsencryptCache
	}
	if *miniserverURL != "" {
		cfg.MiniserverURL = *miniserverURL
	}
	if *miniserverUser != "" {
		cfg.MiniserverUser = *miniserverUser
	}
	if *miniserverPassword != "" {
		cfg.MiniserverPassword = *miniserverPassword
	}
	if *miniserverTimeout != 0 {
		cfg.MiniserverTimeout = *miniserverTimeout
	}
	return cfg
}

func readConfigFile(filename string) (name string, f []byte, err error) {
	if filename == "" {
		f, err = ioutil.ReadFile(defaultConfigFile)
		if err == nil {
			name = defaultConfigFile
		}
		// Ignore errors because we might get all nedded parameters
		// from flags and/or environment variables
		err = nil
		return
	}
	f, err = ioutil.ReadFile(filename)
	if err == nil {
		name = filename
	}
	return
}

func newFileConfig(filename string) (*basicTypeConfig, string, error) {
	cfg := newDefaultConfig()
	fn, f, err := readConfigFile(filename)
	if err != nil {
		return cfg, "", err
	}
	if len(f) == 0 {
		// No config file found
		return cfg, "", nil
	}
	err = toml.Unmarshal(f, cfg)
	if err != nil {
		return cfg, "", errors.Wrap(err, "Error unmarschaling toml data from "+filename)
	}

	return cfg, fn, nil
}

func mergeConfig(cfg *basicTypeConfig, c *basicTypeConfig) {
	defCfg := newDefaultConfig()
	if c.ConfigFile != defCfg.ConfigFile {
		cfg.ConfigFile = c.ConfigFile
	}
	if c.LogFileMain != defCfg.LogFileMain {
		cfg.LogFileMain = c.LogFileMain
	}
	if c.LogFileHTTPError != defCfg.LogFileHTTPError {
		cfg.LogFileHTTPError = c.LogFileHTTPError
	}
	if c.LogFileHTTPAccess != defCfg.LogFileHTTPAccess {
		cfg.LogFileHTTPAccess = c.LogFileHTTPAccess
	}
	if c.ControlsFiles != defCfg.ControlsFiles {
		cfg.ControlsFiles = c.ControlsFiles
	}
	if c.ListenPort != defCfg.ListenPort {
		cfg.ListenPort = c.ListenPort
	}
	if c.PublicURI != defCfg.PublicURI {
		cfg.PublicURI = c.PublicURI
	}
	if c.LetsEncryptCache != defCfg.LetsEncryptCache {
		cfg.LetsEncryptCache = c.LetsEncryptCache
	}
	if c.ControlsFiles != defCfg.ControlsFiles {
		cfg.ControlsFiles = c.ControlsFiles
	}
	if c.MiniserverURL != defCfg.MiniserverURL {
		cfg.MiniserverURL = c.MiniserverURL
	}
	if c.MiniserverUser != defCfg.MiniserverUser {
		cfg.MiniserverUser = c.MiniserverUser
	}
	if c.MiniserverPassword != defCfg.MiniserverPassword {
		cfg.MiniserverPassword = c.MiniserverPassword
	}
	if c.MiniserverTimeout != defCfg.MiniserverTimeout {
		cfg.MiniserverTimeout = c.MiniserverTimeout
	}
	return
}

// NewConfig return an initialized Config struct
func NewConfig(version string) (*Config, error) {
	var err error
	envCfg, err := newEnvConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Error reading config from environment")
	}
	ConfigFile := envCfg.ConfigFile
	flagCfg := newFlagConfig(version)
	if flagCfg.ConfigFile != "" {
		ConfigFile = flagCfg.ConfigFile
	}
	fileCfg, usedConfigFile, err := newFileConfig(ConfigFile)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading config from file")
	}
	impCfg := newDefaultConfig()
	impCfg.ConfigFile = usedConfigFile
	mergeConfig(impCfg, fileCfg)
	mergeConfig(impCfg, envCfg)
	mergeConfig(impCfg, flagCfg)
	cfg, err := impCfg.getConfig()
	if err != nil {
		return cfg, err
	}
	cfg.Version = version
	err = cfg.Validate()
	if err != nil {
		return cfg, errors.Wrap(err, "Error validating config")
	}
	return cfg, nil
}
