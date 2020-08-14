package typescript

import (
	"reflect"
	"testing"
	"time"
)

type TestUser struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Age       int
	Score     int    `json:"score,string"`
	Admin     bool   `json:"admin,omitempty"`
	Password  string `json:"-"`
	Roles     []string
	CreatedAt time.Time
	DeletedAt *time.Time
	secret    string
	Data      interface{}
}

func TestGenerate(t *testing.T) {
	type args struct {
		goStructs []interface{}
	}
	tests := []struct {
		name string
		args args
		want []Interface
	}{
		{
			name: "basic test",
			args: args{
				goStructs: []interface{}{
					TestUser{},
				},
			},
			want: []Interface{
				{
					Name: "TestUser",
					Fields: []Field{
						{
							Name: "id",
							Type: "number",
						},
						{
							Name: "name",
							Type: "string",
						},
						{
							Name: "Age",
							Type: "number",
						},
						{
							Name: "score",
							Type: "string",
						},
						{
							Name:     "admin",
							Type:     "boolean",
							Optional: true,
						},
						{
							Name: "Roles",
							Type: "string[]",
						},
						{
							Name: "CreatedAt",
							Type: "string",
						},
						{
							Name: "DeletedAt",
							Type: "string",
							Null: true,
						},
						{
							Name: "Data",
							Type: "any",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Generate(tt.args.goStructs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
