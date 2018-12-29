package config

import (
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func removeEnvVars(prefix string) map[string]string {
	oldEnv := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		varName := pair[0]
		varValue := pair[1]
		if strings.HasPrefix(pair[0], prefix) {
			oldEnv[varName] = varValue
			os.Unsetenv(varName)
		}
	}
	return oldEnv
}

func TestNewConfig(t *testing.T) {
	//TODO: This is a huge mess :-(
	testingVersionNumber := "0.0.0"
	envPrefix := "LOXWEBHOOK_"
	defaultConfigFile := ""
	exampleConfigFile := "../config.example.toml"
	envConfigFile := filepath.Join("testdata", "config_for_tests_env.toml")
	envConfigFileLogFileMain := "/var/log/configmain_env.log"
	envConfigFileLogFileHTTPError := "/var/log/configerr_env.log"
	envConfigFileLogFileHTTPAccess := "/var/log/configacc_env.log"
	envConfigFileListenPort := 811
	envConfigFilePublicURI := "config.env.example.com"
	envConfigFileLetsencryptCache := "./cache/letsencrypt/config/env"
	envConfigControlsFiles := "./controls.d/env"
	envConfigFileMiniserverURL := &url.URL{
		Scheme: "http",
		Host:   "192.168.1.71:80",
	}
	envConfigFileMiniserverUser := "envConfigFileMiniserverUser"
	envConfigFileMiniserverPassword := "envConfigFileMiniserverPassword"
	envConfigFileMiniserverTimeout := time.Duration(811) * time.Second
	flagConfigFile := filepath.Join("testdata", "config_for_tests_flag.toml")
	flagConfigFileLogFileMain := "/var/log/configmain_flag.log"
	flagConfigFileLogFileHTTPError := "/var/log/configerr_flag.log"
	flagConfigFileLogFileHTTPAccess := "/var/log/configacc_flag.log"
	flagConfigFileListenPort := 822
	flagConfigPublicURI := "config.flag.example.com"
	flagConfigLetsencryptCache := "./cache/letsencrypt/config/flag"
	flagConfigControlsFiles := "./controls.d/flag"
	flagConfigFileMiniserverURL := &url.URL{
		Scheme: "http",
		Host:   "192.168.1.72:80",
	}
	flagConfigFileMiniserverUser := "flagConfigFileMiniserverUser"
	flagConfigFileMiniserverPassword := "flagConfigFileMiniserverPassword"
	flagConfigFileMiniserverTimeout := time.Duration(822) * time.Second
	defaultLogFileMain := ""
	exampleLogFileMain := ""
	defaultLogFileHTTPError := ""
	exampleLogFileHTTPError := ""
	defaultLogFileHTTPAccess := ""
	exampleLogFileHTTPAccess := ""
	envLogFileMain := filepath.Join("var", "log", "envLogFileMain.log")
	envLogFileHTTPError := filepath.Join("var", "log", "envLogFileHTTPError.log")
	envLogFileHTTPAccess := filepath.Join("var", "log", "envLogFileHTTPAccess.log")
	flagLogFileMain := filepath.Join("var", "log", "flagLogFileMain.log")
	flagLogFileHTTPError := filepath.Join("var", "log", "flagLogFileHTTPError.log")
	flagLogFileHTTPAccess := filepath.Join("var", "log", "flagLogFileHTTPAccess.log")
	defaultListenPort := 80
	exampleListenPort := 80
	envListenPort := 81
	flagListenPort := 82
	defaultPublicURI := ""
	examplePublicURI := "example.example.com"
	envPublicURI := "env.example.com"
	flagPublicURI := "flag.example.com"
	defaultLetsencryptCache := "./cache/letsencrypt"
	exampleLetsencryptCache := "./cache/letsencrypt/example"
	envLetsencryptCache := "./cache/letsencrypt/env"
	flagLetsencryptCache := "./cache/letsencrypt/flag"
	defaultControlsFiles := "./controls.d"
	exampleControlsFiles := "./controls_example.d"
	envControlsFiles := "./controls_env.d"
	flagControlsFiles := "./controls_flag.d"
	defaultMiniserverURL := new(url.URL)
	exampleMiniserverURL := &url.URL{
		Scheme: "http",
		Host:   "192.168.1.2:80",
	}
	envMiniserverURLstr := "http://192.168.1.81:80"
	envMiniserverURL := &url.URL{
		Scheme: "http",
		Host:   "192.168.1.81:80",
	}
	flagMiniserverURLstr := "http://192.168.1.82:80"
	flagMiniserverURL := &url.URL{
		Scheme: "http",
		Host:   "192.168.1.82:80",
	}
	defaultMiniserverUser := "admin"
	exampleMiniserverUser := "user"
	envMiniserverUser := "userEnv"
	flagMiniserverUser := "userEnv"
	defaultMiniserverPassword := "admin"
	exampleMiniserverPassword := "SecretPassword"
	envMiniserverPassword := "env"
	flagMiniserverPassword := "flag"
	defaultMiniserverTimeout := 2 * time.Second
	exampleMiniserverTimeout := 2 * time.Second
	envMiniserverTimeoutStr := "81"
	envMiniserverTimeout := 81 * time.Second
	flagMiniserverTimeoutStr := "82"
	flagMiniserverTimeout := 82 * time.Second
	type args struct {
		configFile *string
	}
	tests := []struct {
		name    string
		env     map[string]string
		flags   []string
		args    args
		wantCfg Config
		wantErr bool
	}{
		{
			name: "allDefaults",
			flags: []string{
				os.Args[0],
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         defaultConfigFile,
				LogFileMain:        defaultLogFileMain,
				LogFileHTTPError:   defaultLogFileHTTPError,
				LogFileHTTPAccess:  defaultLogFileHTTPAccess,
				ListenPort:         defaultListenPort,
				PublicURI:          defaultPublicURI,
				LetsEncryptCache:   defaultLetsencryptCache,
				ControlsFiles:      defaultControlsFiles,
				MiniserverURL:      defaultMiniserverURL,
				MiniserverUser:     defaultMiniserverUser,
				MiniserverPassword: defaultMiniserverPassword,
				MiniserverTimeout:  defaultMiniserverTimeout,
			},
		},
		{
			name: "exampleConfig",
			flags: []string{
				os.Args[0],
				"-config", exampleConfigFile,
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         exampleConfigFile,
				LogFileMain:        exampleLogFileMain,
				LogFileHTTPError:   exampleLogFileHTTPError,
				LogFileHTTPAccess:  exampleLogFileHTTPAccess,
				ListenPort:         exampleListenPort,
				PublicURI:          examplePublicURI,
				LetsEncryptCache:   exampleLetsencryptCache,
				ControlsFiles:      exampleControlsFiles,
				MiniserverURL:      exampleMiniserverURL,
				MiniserverUser:     exampleMiniserverUser,
				MiniserverPassword: exampleMiniserverPassword,
				MiniserverTimeout:  exampleMiniserverTimeout,
			},
		},
		{
			name: "allEnv",
			env: map[string]string{
				"LOGFILEMAIN":        envLogFileMain,
				"LOGFILEHTTPERROR":   envLogFileHTTPError,
				"LOGFILEHTTPACCESS":  envLogFileHTTPAccess,
				"LISTENPORT":         strconv.Itoa(envListenPort),
				"PUBLICURI":          envPublicURI,
				"LETSENCRYPTCACHE":   envLetsencryptCache,
				"CONTROLSFILES":      envControlsFiles,
				"MINISERVERURL":      envMiniserverURLstr,
				"MINISERVERUSER":     envMiniserverUser,
				"MINISERVERPASSWORD": envMiniserverPassword,
				"MINISERVERTIMEOUT":  envMiniserverTimeoutStr,
			},
			flags: []string{
				os.Args[0],
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         "",
				LogFileMain:        envLogFileMain,
				LogFileHTTPError:   envLogFileHTTPError,
				LogFileHTTPAccess:  envLogFileHTTPAccess,
				ListenPort:         envListenPort,
				PublicURI:          envPublicURI,
				LetsEncryptCache:   envLetsencryptCache,
				ControlsFiles:      envControlsFiles,
				MiniserverURL:      envMiniserverURL,
				MiniserverUser:     envMiniserverUser,
				MiniserverPassword: envMiniserverPassword,
				MiniserverTimeout:  envMiniserverTimeout,
			},
		},
		{
			name: "envConfigFile",
			env: map[string]string{
				"CONFIG": envConfigFile,
			},
			flags: []string{
				os.Args[0],
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         envConfigFile,
				LogFileMain:        envConfigFileLogFileMain,
				LogFileHTTPError:   envConfigFileLogFileHTTPError,
				LogFileHTTPAccess:  envConfigFileLogFileHTTPAccess,
				ListenPort:         envConfigFileListenPort,
				PublicURI:          envConfigFilePublicURI,
				LetsEncryptCache:   envConfigFileLetsencryptCache,
				ControlsFiles:      envConfigControlsFiles,
				MiniserverURL:      envConfigFileMiniserverURL,
				MiniserverUser:     envConfigFileMiniserverUser,
				MiniserverPassword: envConfigFileMiniserverPassword,
				MiniserverTimeout:  envConfigFileMiniserverTimeout,
			},
		},
		{
			name: "allFlags",
			flags: []string{
				os.Args[0],
				"-config", flagConfigFile,
				"-logfilemain", flagLogFileMain,
				"-logfilehttperror", flagLogFileHTTPError,
				"-logfilehttpaccess", flagLogFileHTTPAccess,
				"-listenport", strconv.Itoa(flagListenPort),
				"-publicURI", flagPublicURI,
				"-letsencryptCache", flagLetsencryptCache,
				"-controlsfiles", flagControlsFiles,
				"-miniserverURL", flagMiniserverURLstr,
				"-miniserverUser", flagMiniserverUser,
				"-miniserverPassword", flagMiniserverPassword,
				"-miniserverTimeout", flagMiniserverTimeoutStr,
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         flagConfigFile,
				LogFileMain:        flagLogFileMain,
				LogFileHTTPError:   flagLogFileHTTPError,
				LogFileHTTPAccess:  flagLogFileHTTPAccess,
				ListenPort:         flagListenPort,
				PublicURI:          flagPublicURI,
				LetsEncryptCache:   flagLetsencryptCache,
				ControlsFiles:      flagControlsFiles,
				MiniserverURL:      flagMiniserverURL,
				MiniserverUser:     flagMiniserverUser,
				MiniserverPassword: flagMiniserverPassword,
				MiniserverTimeout:  flagMiniserverTimeout,
			},
		},
		{
			name: "flagConfigFile",
			flags: []string{
				os.Args[0],
				"-config", flagConfigFile,
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         flagConfigFile,
				LogFileMain:        flagConfigFileLogFileMain,
				LogFileHTTPError:   flagConfigFileLogFileHTTPError,
				LogFileHTTPAccess:  flagConfigFileLogFileHTTPAccess,
				ListenPort:         flagConfigFileListenPort,
				PublicURI:          flagConfigPublicURI,
				LetsEncryptCache:   flagConfigLetsencryptCache,
				ControlsFiles:      flagConfigControlsFiles,
				MiniserverURL:      flagConfigFileMiniserverURL,
				MiniserverUser:     flagConfigFileMiniserverUser,
				MiniserverPassword: flagConfigFileMiniserverPassword,
				MiniserverTimeout:  flagConfigFileMiniserverTimeout,
			},
		},
		{
			name: "flagAndDefault",
			flags: []string{
				os.Args[0],
				"-listenport", strconv.Itoa(flagListenPort),
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         defaultConfigFile,
				LogFileMain:        defaultLogFileMain,
				LogFileHTTPError:   defaultLogFileHTTPError,
				LogFileHTTPAccess:  defaultLogFileHTTPAccess,
				ListenPort:         flagListenPort,
				PublicURI:          defaultPublicURI,
				LetsEncryptCache:   defaultLetsencryptCache,
				ControlsFiles:      defaultControlsFiles,
				MiniserverURL:      defaultMiniserverURL,
				MiniserverUser:     defaultMiniserverUser,
				MiniserverPassword: defaultMiniserverPassword,
				MiniserverTimeout:  defaultMiniserverTimeout,
			},
		},
	}
	for _, tt := range tests {
		oldEnv := removeEnvVars(envPrefix)
		for name, value := range tt.env {
			err := os.Setenv(envPrefix+name, value)
			if err != nil {
				t.Error(err)
			}
		}
		defer func() {
			for name := range tt.env {
				os.Unsetenv(envPrefix + name)
			}
			for varName, varValue := range oldEnv {
				os.Setenv(varName, varValue)
			}
		}()
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = tt.flags
		t.Run(tt.name, func(t *testing.T) {
			gotCfg, err := NewConfig("0.0.0")
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*gotCfg, tt.wantCfg) {
				t.Errorf("NewConfig() = %v, want %v", *gotCfg, tt.wantCfg)
			}
		})
	}
}
