package typescript

import "testing"

func TestInterface_String(t *testing.T) {
	tests := []struct {
		name      string
		Interface Interface
		want      string
	}{
		{
			Interface: Interface{
				Name:   "User",
				Fields: []Field{},
			},
			want: "export interface User {}",
		},
		{
			Interface: Interface{
				Name: "User",
				Fields: []Field{
					{
						Name: "id",
						Type: "number",
					},
				},
			},
			want: "export interface User {\n\tid: number;\n}",
		},
		{
			Interface: Interface{
				Name: "User",
				Fields: []Field{
					{
						Name: "id",
						Type: "number",
					},
					{
						Name: "name",
						Type: "string",
					},
				},
			},
			want: "export interface User {\n\tid: number;\n\tname: string;\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.Interface.String(); got != tt.want {
				t.Errorf("Interface.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
