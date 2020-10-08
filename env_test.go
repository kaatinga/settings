package env_loader

import (
	"fmt"
	"os"
	"testing"
)

// test structure #1
type goodEnvironmentSettings1 struct {
	Port           string `env:"PORT"`
	PathToDatabase string `env:"DB"`
}

// test structure #2
type badEnvironmentSettings1 struct {
	Port           string
	PathToDatabase string `env:"DB"`
}

// test structure #3
type goodEnvironmentSettings2 struct {
	Port           string `env:"DB"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"CACHE"`
}

// test structure #4
type badEnvironmentSettings2 struct {
	Port           string `env:"DB"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE3"`
}

func TestLoadUsingReflect(t *testing.T) {

	err := os.Setenv("BADCACHE3", "-1")
	if err != nil {
		t.Errorf("ENV variables has not been set")
	}

	var goodSettings1 goodEnvironmentSettings1
	var goodSettings2 goodEnvironmentSettings2
	var badSettings1 badEnvironmentSettings1
	var badSettings2 badEnvironmentSettings2

	tests := []struct {
		name     string
		settings interface{}
		wantErr  bool
	}{
		{"ok1", &goodSettings1, false},
		{"ok2", &goodSettings2, false},
		{"!ok1", &badSettings1, true},
		{"!ok2", &badSettings2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err = LoadUsingReflect(tt.settings); (err != nil) != tt.wantErr {
				t.Errorf("LoadUsingReflect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	fmt.Println(err.Error())
	fmt.Println(goodSettings1)
	fmt.Println(goodSettings2)
	fmt.Println(badSettings1)
}
