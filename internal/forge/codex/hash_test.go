package codex

import "testing"

func TestHash(t *testing.T) {
	type args struct {
		source string
		salt   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "basic test",
			args: args{
				source: "password123",
				salt:   "hello world!",
			},
			want: "6127e581e3ab2c7a444e9523a98efc5d5bd74b5febbad4627ec6618c7f2fdee6",
		},
		{
			name: "make sure salt changes the result",
			args: args{
				source: "password123",
				salt:   "pepper",
			},
			want: "711bd97d28a1cfa75701ebd2e7314f2ff04e8059296afac90707912f84c67eac",
		},
		{
			name: "test blank",
			args: args{
				source: "",
				salt:   "",
			},
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.source, tt.args.salt); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
