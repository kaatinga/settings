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
	CacheSize      byte   `env:"BADCACHE1"`
}

// test structure #5
type badEnvironmentSettings3 struct {
	Port           string `env:"DB"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE2"`
}

// test structure #6
type badEnvironmentSettings4 struct {
	Port           string `env:"DB"`
	PathToDatabase string `env:"DB"`
	CacheSize      byte   `env:"BADCACHE3"`
}

func TestLoadUsingReflect(t *testing.T) {

	// ENV settings PORT=80;DB=db/file;CACHE=5;BADCACHE1=i;BADCACHE2=300
	var err error
	err = os.Setenv("PORT", "80")
	err = os.Setenv("DB", "db/file")
	err = os.Setenv("CACHE", "5")
	err = os.Setenv("BADCACHE1", "i")
	err = os.Setenv("BADCACHE2", "300")
	err = os.Setenv("BADCACHE3", "-1")
	if err != nil {
		t.Errorf("ENV variables has not been set")
	}

	var goodSettings1 goodEnvironmentSettings1
	var goodSettings2 goodEnvironmentSettings2
	var badSettings1 badEnvironmentSettings1
	var badSettings2 badEnvironmentSettings2
	var badSettings3 badEnvironmentSettings3
	var badSettings4 badEnvironmentSettings4

	tests := []struct {
		name     string
		settings interface{}
		wantErr  bool
	}{
		{"ok1", &goodSettings1, false},
		{"ok2", &goodSettings2, false},
		{"!ok1", &badSettings1, true},
		{"!ok2", &badSettings2, true},
		{"!ok3", &badSettings3, true},
		{"!ok4", &badSettings4, true},
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
