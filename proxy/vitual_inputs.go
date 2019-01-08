package proxy

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var virtualInputBasePath = "/dev/sps/io"

type digitalVirtualInput struct {
	ID      int
	Command string
	Token   string
}

func (vi *digitalVirtualInput) setCommand(command string) (err error) {
	validCommands := map[string]string{
		"ein":    "On",
		"on":     "On",
		"aus":    "Off",
		"off":    "Off",
		"impuls": "Pulse",
		"pulse":  "Pulse",
	}
	cmd, ok := validCommands[strings.ToLower(command)]
	if !ok {
		return fmt.Errorf("Unknown command for digital virtual input: %s", command)
	}
	vi.Command = cmd
	return nil
}

func (vi *digitalVirtualInput) getControl() string {
	return "VI" + strconv.Itoa(vi.ID)
}

// GetPath returns a path that sends a command to the Miniserver
func (vi *digitalVirtualInput) GetPath() string {
	ep := fmt.Sprintf("%s/%s/%s",
		virtualInputBasePath,
		vi.getControl(),
		vi.Command,
	)
	return ep
}

func newDigitalVirtualInput(id int, command, token string) (*digitalVirtualInput, error) {
	vi := new(digitalVirtualInput)
	err := vi.setCommand(command)
	if err != nil {
		return vi, err
	}
	vi.Token = token
	return vi, nil
}

// newDigitalVirtualInput returns a DigitalVitualEndpoint with data parsed from req
func parseRequestDigitalVirtualInput(req *http.Request) (ID int, command, token string, err error) {
	ID, err = strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		err = errors.Wrap(err, "Error casting virtual input ID to integer")
		return
	}
	command = mux.Vars(req)["command"]
	token = req.URL.Query().Get("t")
	return
}
