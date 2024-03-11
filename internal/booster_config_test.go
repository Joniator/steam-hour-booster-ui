package internal

import (
	"io"
	"os"
	"slices"
	"testing"
)

func TestConfig(t *testing.T) {
	sourceFile, err := os.Open("testdata/config.json")
	if err != nil {
		t.Fatalf("Failed to setup test: %s", err.Error())
	}
	destFile, err := os.Create("testdata/working_config.json")
	if err != nil {
		t.Fatalf("Failed to setup test: %s", err.Error())
	}
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		t.Fatalf("Failed to setup test: %s", err.Error())
	}

	config, err := LoadBoosterConfig("testdata/working_config.json")

	u1, err := config.GetUserConfig("User1")
	if err != nil || !slices.Equal(u1.AppIds, []int{123}) {
		t.Error("Failed to load user1")
	}
	u2, err := config.GetUserConfig("User2")
	if err != nil || !slices.Equal(u2.AppIds, []int{123, 321}) {
		t.Error("Failed to load user2")
	}
	_, err = config.GetUserConfig("User3")
	if err == nil {
		t.Errorf("Not existing didn't error, err=%v", err)
	}

	err = config.AddGame("User1", 123)
	if err != nil || !slices.Equal(u1.AppIds, []int{123}) {
		t.Errorf("Existing game should not duplicate, list=%v, err=%v", u1, err)
	}

	err = config.AddGame("User1", 1337)
	if err != nil || !slices.Equal(u1.AppIds, []int{123, 1337}) {
		t.Errorf("Game should have been added, list=%v, err=%v", u1, err)
	}

	err = config.AddGame("User1", 512)
	if err != nil || !slices.Equal(u1.AppIds, []int{123, 512, 1337}) {
		t.Errorf("List should be sorted, list=%v, err=%v", u1, err)
	}

	err = config.AddGame("User2", 1024)
	if err != nil || !slices.Equal(u2.AppIds, []int{123, 321, 1024}) {
		t.Errorf("List should be sorted, list=%v, err=%v", u2, err)
	}

	err = config.Save()
	if err != nil {
		t.Fatalf("Failed to save config, err=%v", err)
	}

	newConfig, err := LoadBoosterConfig("testdata/working_config.json")
	if err != nil {
		t.Fatalf("Failed to load new config, err=%v", err)
	}

	var u1checked, u2checked bool
	for _, c := range newConfig.UserConfigs {
		switch c.Name {
		case "User1":
			if u1checked || !slices.Equal(c.AppIds, []int{123, 512, 1337}) {
				t.Errorf("User1 not saved correctly, list=%v", c.AppIds)
			}
			u1checked = true
			break
		case "User2":
			if u2checked || !slices.Equal(c.AppIds, []int{123, 321, 1024}) {
				t.Errorf("User2 not saved correctly, list=%v", c.AppIds)
			}
			u2checked = true
			break
		default:
			t.Errorf("Unexpected User %v", c)
			break
		}
	}

}
