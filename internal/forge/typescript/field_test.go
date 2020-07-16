package typescript

import "testing"

func TestField_String(t *testing.T) {
	tests := []struct {
		name  string
		field Field
		want  string
	}{
		{
			name: "standard",
			field: Field{
				Name: "name",
				Type: "string",
			},
			want: "name: string;",
		},
		{
			name: "optional",
			field: Field{
				Name:     "name",
				Optional: true,
				Type:     "string",
			},
			want: "name?: string;",
		},
		{
			name: "optional",
			field: Field{
				Name:  "name",
				Array: true,
				Type:  "string",
			},
			want: "name: string[]|null;",
		},
		{
			name: "null",
			field: Field{
				Name: "name",
				Null: true,
				Type: "string",
			},
			want: "name: string|null;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.field.String(); got != tt.want {
				t.Errorf("Field.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
