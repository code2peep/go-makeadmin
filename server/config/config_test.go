package config

import "testing"

func TestConfigDefaultTablePrefixIsMakeAdmin(t *testing.T) {
	if Config.DbTablePrefix != "ma_" {
		t.Fatalf("Config.DbTablePrefix = %q, want ma_", Config.DbTablePrefix)
	}
}

func TestConfigPathFromArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "short flag", args: []string{"-c", "/tmp/app.env"}, want: "/tmp/app.env"},
		{name: "short flag equals", args: []string{"-c=/tmp/app.env"}, want: "/tmp/app.env"},
		{name: "long flag", args: []string{"--c", "/tmp/app.env"}, want: "/tmp/app.env"},
		{name: "long flag equals", args: []string{"--c=/tmp/app.env"}, want: "/tmp/app.env"},
		{name: "ignore go test flag", args: []string{"-test.testlogfile=/tmp/test.log"}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := configPathFromArgs(tt.args); got != tt.want {
				t.Fatalf("configPathFromArgs(%v) = %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
