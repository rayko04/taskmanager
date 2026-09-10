package main

import (
	"taskmanager/model"
	"testing"
	"io"
	"strings"
)

func TestValidReq(t *testing.T) {

	tests := []struct {
		name 	string
		req		model.TaskRequest
		wantErr bool
	}{
		{
			name: "valid title",
			req: model.TaskRequest{
				Title: "Get Job",
			},
			wantErr: true,
		},
		{
			name: "empty title",
			req: model.TaskRequest{
				Title: "",
			},
			wantErr: false,
		},
		{
			name: "whitespace title",
			req: model.TaskRequest{
				Title: "    ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validReq(&tt.req)

			if got != tt.wantErr {
				t.Errorf("validReq() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestValidUpdateReq(t *testing.T) {

	tests := []struct {
		name	string
		req		model.TaskUpdateRequest
		want	bool
	}{
		{
			name: "No Title",
			req: model.TaskUpdateRequest{
				Title: nil,
				Description: new("test"),
				Completed: new(false),
			},
			want: false,
		},
		{
			name: "No Description",
			req: model.TaskUpdateRequest{
				Title: new("test"),
				Description: nil,
				Completed: new(false),
			},
			want: false,
		},
		{
			name: "No Completed",
			req: model.TaskUpdateRequest{
				Title: new("test"),
				Description: new("test"),
				Completed: nil,
			},
			want: false,
		},
		{
			name: "None",
			req: model.TaskUpdateRequest{
				Title: nil,
				Description: nil,
				Completed: nil,
			},
			want: false,
		},
		{
			name: "Valid Req",
			req: model.TaskUpdateRequest{
				Title: new("test"),
				Description: new("test"),
				Completed: new(false),
			},
			want: true,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			got := validUpdateReq(&tt.req)

			if got != tt.want{
				t.Errorf("validUpdateReq() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestValidPatchReq(t *testing.T) {

	tests := []struct {
		name	string
		req		model.TaskPatchRequest
		want	bool
	}{
		{
			name: "No Title",
			req: model.TaskPatchRequest{
				Title: nil,
				Description: new("test"),
				Completed: new(false),
			},
			want: true,
		},
		{
			name: "No Description",
			req: model.TaskPatchRequest{
				Title: new("test"),
				Description: nil,
				Completed: new(false),
			},
			want: true,
		},
		{
			name: "No Completed",
			req: model.TaskPatchRequest{
				Title: new("test"),
				Description: new("test"),
				Completed: nil,
			},
			want: true,
		},
		{
			name: "Invalid Req",
			req: model.TaskPatchRequest{
				Title: nil,
				Description: nil,
				Completed: nil,
			},
			want: false,
		},
		{
			name: "All Fields",
			req: model.TaskPatchRequest{
				Title: new("test"),
				Description: new("test"),
				Completed: new(false),
			},
			want: true,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			got := validPatchReq(&tt.req)

			if got != tt.want{
				t.Errorf("validPatchReq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCollectionPath(t *testing.T) {

	tests := []struct {
		name string
		path []string
		want bool
	}{
		{
			name: "valid path",
			path: []string{"/", "tasks"},
			want: true,
		},
		{
			name: "invalid path",
			path: []string{"api", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"api", "", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"/", "task"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"/", "tasks", "1"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{},
			want: false,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			got := isCollectionPath(tt.path)

			if got != tt.want{
				t.Errorf("isCollectionPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIndividualPath(t *testing.T) {

	tests := []struct {
		name string
		path []string
		want bool
	}{
		{
			name: "valid path",
			path: []string{"/", "tasks", "5"},
			want: true,
		},
		{
			name: "invalid path",
			path: []string{"api", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"api", "", "tasks"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{"/", "tasks", "5", "8"},
			want: false,
		},
		{
			name: "invalid path",
			path: []string{},
			want: false,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			got := isIndividualPath(tt.path)

			if got != tt.want{
				t.Errorf("isIndividualPath() = %v, want %v", got, tt.want)
			}
		})
	}
}


func TestDecodeJSON(t *testing.T) {

	tests := []struct {
		name string
		read io.ReadCloser
		want bool
	}{
		{
			name: "valid json",
			read : io.NopCloser(strings.NewReader(`{"title":"Buy milk","description":"From the store"}`)),
			want: false,
		},
		{
			name: "empty",
			read : io.NopCloser(strings.NewReader(`{}`)),
			want: false,
		},
		{
			name: "invalid json",
			read : io.NopCloser(strings.NewReader(`{"title":"Buy milk"`)),
			want: true,
		},
		{
			name: "unknown field",
			read : io.NopCloser(strings.NewReader(`{"title":"Buy milk","foo":"From the store"}`)),
			want: true,
		},
		{
			name: "multiple json",
			read : io.NopCloser(strings.NewReader(`{"title":"Buy milk"}{"description":"From the store"}`)),
			want: true,
		},
		{
			name: "trailing garbage",
			read : io.NopCloser(strings.NewReader(`{"title":"Buy milk","description":"From the store"} garbage`)),
			want: true,
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T) {
			got := false
			err := decodeJSON(tt.read, &model.TaskRequest{})

			if err != nil {
				got = true
			}

			if got != tt.want{
				t.Errorf("decodeJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
