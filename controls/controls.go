package controls

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type controlImport struct {
	Tokens   map[string]string
	Controls map[string]Control
}

func (ci controlImport) Validate() ControlError {
	validName := regexp.MustCompile(`^[0-9a-zA-z_-]+$`)
	for name, c := range ci.Controls {
		if !validName.MatchString(name) {
			return newInvalidControlNameError(name)
		}
		err := c.Validate()
		if err != nil {
			return err
		}
		// Check if token configured in this control exists
		for _, c := range ci.Controls {
			for _, t := range c.Tokens {
				if _, ok := ci.Tokens[t]; !ok {
					return newInvalidTokenError(t)
				}
			}
		}
	}
	return nil
}

// Control holds the config for one Miniserver control
type Control struct {
	Category string
	ID       int
	Allowed  []string
	Tokens   []string
}

func (c *Control) validateAllowedCommandsDvi() ControlError {
	// Loxone documentation for allowed commands: https://www.loxone.com/enen/kb/web-services/
	for _, command := range c.Allowed {
		switch strings.ToLower(command) {
		case
			"<all>",
			"0",
			"1",
			"on",
			"off",
			"impuls",
			"pulse",
			"impulsplus",
			"impulsminus",
			"pulseup",
			"pulsedown",
			"impulsauf",
			"impulsab",
			"pulseopen",
			"pulseclose",
			"plusein",
			"plusaus",
			"upon",
			"upoff",
			"aufein",
			"aufaus",
			"openon",
			"openoff",
			"minusein",
			"minusaus",
			"downon",
			"downoff",
			"abein",
			"abaus",
			"closeon",
			"closeoff":
			// Nothing to do
		default:
			return newInvalidCommandError("dvi", command)
		}
	}
	return nil
}

// Validate returns an error if a control contains invalid data
func (c *Control) Validate() ControlError {
	switch c.Category {
	case
		"dvi":
		if err := c.validateAllowedCommandsDvi(); err != nil {
			return err
		}
	default:
		return newInvalidCategoryError(c.Category)
	}

	return nil
}

// Read imports all *.toml files from dir (including subdirectories) and returns
// tokens and controls
func Read(dir string) (map[string]string, map[string]Control, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".toml" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	tokens := make(map[string]string)
	controls := make(map[string]Control)
	impCtl := new(controlImport)
	for _, fn := range files {
		err = importFile(impCtl, fn)
		if err != nil {
			return nil, nil, errors.Wrap(err, "Error importing control definitions from file")
		}
	}
	err = impCtl.Validate()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error validating controls")
	}
	for k, v := range impCtl.Tokens {
		tokens[k] = v
	}
	for k, v := range impCtl.Controls {
		controls[k] = v
	}
	return tokens, controls, nil
}

func importFile(impCtl *controlImport, fn string) error {
	var err error
	fc, err := ioutil.ReadFile(fn)
	if err != nil {
		return errors.Wrap(err, "Error opening file with control definitions")
	}
	err = toml.Unmarshal(fc, impCtl)
	if err != nil {
		err = errors.Wrap(err, "Error unmarschaling toml data")
		return err
	}
	return err
}
