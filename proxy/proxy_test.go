package proxy

import (
	"testing"

	"github.com/axxelG/loxwebhook/controls"
)

func Test_authorize(t *testing.T) {
	tokens := map[string]string{
		"test1": "f7932d8a-b37f-46dc-84ee-276c545aec48",
		"test2": "88f3cc74-b741-404e-b6a3-136d76796de8",
		"test3": "d7d47ae7-44d6-4b4b-b65d-06e7f5bf108e",
	}
	ctl := controls.Control{
		Category: "dvi",
		ID:       1,
		Allowed: []string{
			"pulse",
			"impuls",
		},
		Tokens: []string{
			"test1",
			"test2",
		},
	}
	type args struct {
		control    controls.Control
		tokens     map[string]string
		reqToken   string
		reqCommand string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ValidAuth",
			args: args{
				control:    ctl,
				tokens:     tokens,
				reqToken:   "88f3cc74-b741-404e-b6a3-136d76796de8",
				reqCommand: "pulse",
			},
			wantErr: false,
		},
		{
			name: "InvalidCommand",
			args: args{
				control:    ctl,
				tokens:     tokens,
				reqToken:   "88f3cc74-b741-404e-b6a3-136d76796de8",
				reqCommand: "on",
			},
			wantErr: true,
		},
		{
			name: "InvalidExistingToken",
			args: args{
				control:    ctl,
				tokens:     tokens,
				reqToken:   "d7d47ae7-44d6-4b4b-b65d-06e7f5bf108e",
				reqCommand: "pulse",
			},
			wantErr: true,
		},
		{
			name: "InvalidNonExistingToken",
			args: args{
				control:    ctl,
				tokens:     tokens,
				reqToken:   "8c9564a1-6af7-4ed0-8656-add107e882a6",
				reqCommand: "pulse",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := authorize(tt.args.control, tt.args.tokens, tt.args.reqToken, tt.args.reqCommand); (err != nil) != tt.wantErr {
				t.Errorf("authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
