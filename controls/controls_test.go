package controls

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name         string
		args         args
		wantAuthKeys map[string]string
		wantControls map[string]Control
		wantErr      bool
	}{
		{
			name: "OneFile",
			args: args{
				dir: filepath.Join("testdata", "OneFile"),
			},
			wantAuthKeys: map[string]string{
				"testOne":   "43b2c690-f281-42bb-af2d-979f5dbe9517",
				"testTwo":   "69b9a1ad-1224-4c93-8411-e88e65ebe582",
				"testThree": "84627dbd-bd68-476f-9e53-35522285783b",
			},
			wantControls: map[string]Control{
				"test1": Control{
					Category: "dvi",
					ID:       1,
					Allowed: []string{
						"<all>",
					},
					AuthKeys: []string{
						"testOne",
					},
				},
				"test2": Control{
					Category: "dvi",
					ID:       2,
					Allowed: []string{
						"on",
					},
					AuthKeys: []string{
						"testTwo",
					},
				},
				"test3": Control{
					Category: "dvi",
					ID:       3,
					Allowed: []string{
						"on",
						"off",
					},
					AuthKeys: []string{
						"testOne",
						"testThree",
					},
				},
			},
		},
		{
			name: "ThreeFiles",
			args: args{
				dir: filepath.Join("testdata", "ThreeFiles"),
			},
			wantAuthKeys: map[string]string{
				"testOne":   "43b2c690-f281-42bb-af2d-979f5dbe9517",
				"testTwo":   "69b9a1ad-1224-4c93-8411-e88e65ebe582",
				"testThree": "84627dbd-bd68-476f-9e53-35522285783b",
			},
			wantControls: map[string]Control{
				"test1": Control{
					Category: "dvi",
					ID:       1,
					Allowed: []string{
						"<all>",
					},
					AuthKeys: []string{
						"testOne",
					},
				},
				"test2": Control{
					Category: "dvi",
					ID:       2,
					Allowed: []string{
						"on",
					},
					AuthKeys: []string{
						"testTwo",
					},
				},
				"test3": Control{
					Category: "dvi",
					ID:       3,
					Allowed: []string{
						"on",
						"off",
					},
					AuthKeys: []string{
						"testOne",
						"testThree",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAuthKeys, gotControls, err := Read(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAuthKeys, tt.wantAuthKeys) {
				t.Errorf("Read() gotAuthKeys = %v, want %v", gotAuthKeys, tt.wantAuthKeys)
			}
			if !reflect.DeepEqual(gotControls, tt.wantControls) {
				t.Errorf("Read() gotControls = %v, want %v", gotControls, tt.wantControls)
			}
		})
	}
}

func TestControl_Validate(t *testing.T) {
	tests := []struct {
		name string
		c    *Control
		want ControlError
	}{
		{
			name: "validDvi",
			c: &Control{
				Category: "dvi",
				ID:       1,
				Allowed: []string{
					"<all>",
				},
				AuthKeys: []string{
					"f6694286-66e6-4b79-8936-9e45284eba60",
				},
			},
			want: nil,
		},
		{
			name: "invalidCategory",
			c: &Control{
				Category: "NonExistendCategory",
				ID:       1,
				Allowed: []string{
					"<all>",
				},
				AuthKeys: []string{
					"f6694286-66e6-4b79-8936-9e45284eba60",
				},
			},
			want: newInvalidCategoryError("DummyCategory"),
		},
		{
			name: "invalidAllowedCommand",
			c: &Control{
				Category: "dvi",
				ID:       1,
				Allowed: []string{
					"NotAllowedCommand",
				},
				AuthKeys: []string{
					"f6694286-66e6-4b79-8936-9e45284eba60",
				},
			},
			want: newInvalidCommandError("DummyCategory", "DummyCommand"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if (tt.want == nil) && (err != nil) {
				t.Errorf("Control.Validate() error = %v", err)
			}
			if (tt.want != nil) && (err == nil) {
				t.Errorf("Got no error but error was expected")
			}
			if (tt.want != nil) && (err.GetType() != tt.want.GetType()) {
				t.Errorf("Got wrong error type. Got: %s Want: %s", err.GetType(), tt.want.GetType())
			}
		})
	}
}

func Test_controlImport_Validate(t *testing.T) {
	tests := []struct {
		name string
		ci   controlImport
		want ControlError
	}{
		{
			name: "ValidName",
			ci: controlImport{
				AuthKeys: map[string]string{
					"ValidAuthKey": "325ce159-0ddf-433a-966f-a94b313a7eb5",
				},
				Controls: map[string]Control{
					"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-": {
						Category: "dvi",
						ID:       1,
						Allowed: []string{
							"<all>",
						},
						AuthKeys: []string{
							"ValidAuthKey",
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "InvalidNameSpace",
			ci: controlImport{
				AuthKeys: map[string]string{
					"ValidAuthKey": "325ce159-0ddf-433a-966f-a94b313a7eb5",
				},
				Controls: map[string]Control{
					"No spaces please": {
						Category: "dvi",
						ID:       1,
						Allowed: []string{
							"<all>",
						},
						AuthKeys: []string{
							"ValidAuthKey",
						},
					},
				},
			},
			want: newInvalidControlNameError("No spaces please"),
		},
		{
			name: "InvalidNamePlus",
			ci: controlImport{
				AuthKeys: map[string]string{
					"ValidAuthKey": "325ce159-0ddf-433a-966f-a94b313a7eb5",
				},
				Controls: map[string]Control{
					"No+please": {
						Category: "dvi",
						ID:       1,
						Allowed: []string{
							"<all>",
						},
						AuthKeys: []string{
							"ValidAuthKey",
						},
					},
				},
			},
			want: newInvalidControlNameError("No+please"),
		},
		{
			name: "InvalidNameColon",
			ci: controlImport{
				AuthKeys: map[string]string{
					"ValidAuthKey": "325ce159-0ddf-433a-966f-a94b313a7eb5",
				},
				Controls: map[string]Control{
					"No:please": {
						Category: "dvi",
						ID:       1,
						Allowed: []string{
							"<all>",
						},
						AuthKeys: []string{
							"ValidAuthKey",
						},
					},
				},
			},
			want: newInvalidControlNameError("No:please"),
		},
		{
			name: "ValidAuthKey",
			ci: controlImport{
				AuthKeys: map[string]string{
					"ValidAuthKey": "325ce159-0ddf-433a-966f-a94b313a7eb5",
				},
				Controls: map[string]Control{
					"ControlName": {
						Category: "dvi",
						ID:       1,
						Allowed: []string{
							"<all>",
						},
						AuthKeys: []string{
							"ValidAuthKey",
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "InvalidAuthKey",
			ci: controlImport{
				AuthKeys: map[string]string{
					"ValidAuthKey": "325ce159-0ddf-433a-966f-a94b313a7eb5",
				},
				Controls: map[string]Control{
					"ControlName": {
						Category: "dvi",
						ID:       1,
						Allowed: []string{
							"<all>",
						},
						AuthKeys: []string{
							"InvalidAuthKey",
						},
					},
				},
			},
			want: newInvalidAuthKeyError("DummyName"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ci.Validate()
			if tt.want == nil && got != nil {
				t.Errorf("controlImport.Validate() = %v, want %v", got, tt.want)
			}
			if (tt.want != nil) && (got.GetType() != tt.want.GetType()) {
				t.Errorf("Got wrong error type. Got: %s Want: %s", got.GetType(), tt.want.GetType())
			}
		})
	}
}
