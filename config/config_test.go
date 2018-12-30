package config

import (
	"fmt"
	"net/url"
	"os"
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
	_, err := os.Stat(defaultConfigFile)
	if err == nil {
		tempname := defaultConfigFile + "disabled_for_testing"
		err = os.Rename(defaultConfigFile, tempname)
		if err != nil {
			t.Errorf("Error renaming config file %s to %s: %s", defaultConfigFile, tempname, err)
			return
		}
		defer func() {
			err := os.Rename(tempname, defaultConfigFile)
			if err != nil {
				t.Errorf("Error restoring original config file. Renaming %s to %s failed: %s", tempname, defaultConfigFile, err)
			}
		}()
	}
	envPrefix := "LOXWEBHOOK_"
	testingVersionNumber := "0.0.0"

	configDefaults := Config{
		Version:            testingVersionNumber,
		ConfigFile:         "",
		PublicURI:          "",
		ListenPort:         80,
		MiniserverURL:      new(url.URL),
		MiniserverUser:     "admin",
		MiniserverPassword: "admin",
		MiniserverTimeout:  2 * time.Second,
		LetsEncryptCache:   "./cache/letsencrypt",
		LogFileMain:        "",
		LogFileHTTPError:   "",
		LogFileHTTPAccess:  "",
		ControlsFiles:      "./controls.d",
	}

	configFileExample := Config{
		Version:    testingVersionNumber,
		ConfigFile: "../config.example.toml",
		PublicURI:  "loxwebhook.example.com",
		ListenPort: 4443,
		MiniserverURL: &url.URL{
			Scheme: "http",
			Host:   "192.168.1.1:80",
		},
		MiniserverUser:     "loxwebhook",
		MiniserverPassword: "YourSecretPassword",
		MiniserverTimeout:  2 * time.Second,
		LetsEncryptCache:   "~/.loxwebhook/cache/letsencrypt",
		LogFileMain:        "/var/log/loxwebhook/loxwebhook.log",
		LogFileHTTPError:   "/var/log/loxwebhook/error.log",
		LogFileHTTPAccess:  "/var/log/loxwebhook/access.log",
		ControlsFiles:      "/etc/loxwebhook/controls.d",
	}

	configEnv := Config{
		Version:    testingVersionNumber,
		ConfigFile: configDefaults.ConfigFile,
		PublicURI:  "env.example.com",
		ListenPort: 81,
		MiniserverURL: &url.URL{
			Scheme: "http",
			Host:   "192.168.1.81:80",
		},
		MiniserverUser:     "userEnv",
		MiniserverPassword: "env",
		MiniserverTimeout:  81 * time.Second,
		LetsEncryptCache:   "./cache/letsencrypt/env",
		LogFileMain:        "/var/log/envLogFileMain.log",
		LogFileHTTPError:   "/var/log/envLogFileHTTPError.log",
		LogFileHTTPAccess:  "/var/log/envLogFileHTTPAccess.log",
		ControlsFiles:      "./controls_env.d",
	}

	allEnv := map[string]string{
		"LOGFILEMAIN":        configEnv.LogFileMain,
		"LOGFILEHTTPERROR":   configEnv.LogFileHTTPError,
		"LOGFILEHTTPACCESS":  configEnv.LogFileHTTPAccess,
		"LISTENPORT":         strconv.Itoa(configEnv.ListenPort),
		"PUBLICURI":          configEnv.PublicURI,
		"LETSENCRYPTCACHE":   configEnv.LetsEncryptCache,
		"CONTROLSFILES":      configEnv.ControlsFiles,
		"MINISERVERURL":      configEnv.MiniserverURL.String(),
		"MINISERVERUSER":     configEnv.MiniserverUser,
		"MINISERVERPASSWORD": configEnv.MiniserverPassword,
		"MINISERVERTIMEOUT":  fmt.Sprint(configEnv.MiniserverTimeout.Seconds()),
	}

	configFlag := Config{
		Version:    testingVersionNumber,
		ConfigFile: configDefaults.ConfigFile,
		PublicURI:  "flag.example.com",
		ListenPort: 82,
		MiniserverURL: &url.URL{
			Scheme: "http",
			Host:   "192.168.1.82:80",
		},
		MiniserverUser:     "userFlag",
		MiniserverPassword: "flag",
		MiniserverTimeout:  82 * time.Second,
		LetsEncryptCache:   "./cache/letsencrypt/flag",
		LogFileMain:        "/var/log/flagLogFileMain.log",
		LogFileHTTPError:   "/var/log/flagLogFileHTTPError.log",
		LogFileHTTPAccess:  "/var/log/flagLogFileHTTPAccess.log",
		ControlsFiles:      "./controls_flag.d",
	}

	allFlags := []string{
		os.Args[0],
		"-logfilemain", configFlag.LogFileMain,
		"-logfilehttperror", configFlag.LogFileHTTPError,
		"-logfilehttpaccess", configFlag.LogFileHTTPAccess,
		"-listenport", strconv.Itoa(configFlag.ListenPort),
		"-publicURI", configFlag.PublicURI,
		"-letsencryptCache", configFlag.LetsEncryptCache,
		"-controlsfiles", configFlag.ControlsFiles,
		"-miniserverURL", configFlag.MiniserverURL.String(),
		"-miniserverUser", configFlag.MiniserverUser,
		"-miniserverPassword", configFlag.MiniserverPassword,
		"-miniserverTimeout", fmt.Sprint(configFlag.MiniserverTimeout.Seconds()),
	}
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
			wantCfg: configDefaults,
		},
		{
			name: "allEnv",
			env:  allEnv,
			flags: []string{
				os.Args[0],
			},
			wantCfg: configEnv,
		},
		{
			name:    "allFlags",
			flags:   allFlags,
			wantCfg: configFlag,
		},
		{
			name: "EnvConfigFile",
			flags: []string{
				os.Args[0],
				"-config", configFileExample.ConfigFile,
			},
			wantCfg: configFileExample,
		},
		{
			name: "FlagConfigFile",
			env: map[string]string{
				"CONFIG": configFileExample.ConfigFile,
			},
			flags: []string{
				os.Args[0],
			},
			wantCfg: configFileExample,
		},
		{
			name:    "FlagsOverwriteEnv",
			env:     allEnv,
			flags:   allFlags,
			wantCfg: configFlag,
		},
		{
			name: "EnvAndFlagOverwriteFile",
			env: map[string]string{
				"CONFIG":     configFileExample.ConfigFile,
				"LISTENPORT": strconv.Itoa(configEnv.ListenPort),
			},
			flags: []string{
				os.Args[0],
				"-publicURI", configFlag.PublicURI,
			},
			wantCfg: Config{
				Version:            testingVersionNumber,
				ConfigFile:         configFileExample.ConfigFile,
				ListenPort:         configEnv.ListenPort,
				PublicURI:          configFlag.PublicURI,
				MiniserverURL:      configFileExample.MiniserverURL,
				MiniserverUser:     configFileExample.MiniserverUser,
				MiniserverPassword: configFileExample.MiniserverPassword,
				MiniserverTimeout:  configFileExample.MiniserverTimeout,
				LetsEncryptCache:   configFileExample.LetsEncryptCache,
				LogFileMain:        configFileExample.LogFileMain,
				LogFileHTTPError:   configFileExample.LogFileHTTPError,
				LogFileHTTPAccess:  configFileExample.LogFileHTTPAccess,
				ControlsFiles:      configFileExample.ControlsFiles,
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
