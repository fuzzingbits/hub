package web

import "testing"

func TestGenerateContentSecurityPolicy(t *testing.T) {
	type args struct {
		fileContents []byte
		cspEntries   CSPEntries
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				fileContents: []byte("<style>body { background: block; }</style><script>console.log('foobar');</script>"),
				cspEntries: CSPEntries{
					Default: []string{"'self'"},
					Script:  []string{"'self'"},
					Style:   []string{"'self'"},
					Image:   []string{"'self'"},
				},
			},
			want: "default-src 'self'; script-src 'self' 'sha256-QXZRmRPAsseuAgOGnvjVUJOnlHEzu25Ou1XhFOWnqyI='; style-src 'self' 'sha256-bZB4XqSVB3EpohGYTN6POIJZIjQpgQOsNeoJJkLFkyY='; img-src 'self'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateContentSecurityPolicy(tt.args.fileContents, tt.args.cspEntries)

			if got != tt.want {
				t.Errorf("GenerateContentSecurityPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
