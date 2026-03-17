package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	dir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}
	if dir == "" {
		t.Error("getConfigDir() should not return empty string")
	}
	if filepath.Base(dir) != "worldclock" {
		t.Errorf("getConfigDir() = %q, should end with worldclock", dir)
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("getConfigPath() returned error: %v", err)
	}
	if path == "" {
		t.Error("getConfigPath() should not return empty string")
	}
	if filepath.Ext(path) != ".json" {
		t.Errorf("getConfigPath() = %q, should end with .json", path)
	}
}

func TestLoadDefaultConfig(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpDir, err := os.MkdirTemp("", "worldclock-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}
	if cfg.Timezones == nil {
		t.Error("Load() should return empty Timezones slice, not nil")
	}
	if len(cfg.Timezones) != 0 {
		t.Errorf("Load() with no config should return empty Timezones, got %d", len(cfg.Timezones))
	}
	if cfg.Theme != "dracula" {
		t.Errorf("Load() default theme should be dracula, got %q", cfg.Theme)
	}
}

func TestSaveAndLoad(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpDir, err := os.MkdirTemp("", "worldclock-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)

	cfg := &Config{
		Timezones: []string{"America/New_York", "Europe/London", "Asia/Tokyo"},
		Theme:     "nord",
	}

	err = Save(cfg)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if len(loadedCfg.Timezones) != len(cfg.Timezones) {
		t.Errorf("Loaded Timezones length = %d, want %d", len(loadedCfg.Timezones), len(cfg.Timezones))
	}

	for i, tz := range cfg.Timezones {
		if loadedCfg.Timezones[i] != tz {
			t.Errorf("Timezones[%d] = %q, want %q", i, loadedCfg.Timezones[i], tz)
		}
	}

	if loadedCfg.Theme != cfg.Theme {
		t.Errorf("Theme = %q, want %q", loadedCfg.Theme, cfg.Theme)
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpDir, err := os.MkdirTemp("", "worldclock-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)

	cfg := &Config{
		Timezones: []string{"UTC"},
		Theme:     "dracula",
	}

	err = Save(cfg)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	configPath, _ := getConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Save() should create config file")
	}
}

func TestLoadWithExistingConfig(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpDir, err := os.MkdirTemp("", "worldclock-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".config", "worldclock")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.json")

	configContent := `{"timezones":["America/Chicago","Europe/Paris"],"theme":"gruvbox"}`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if len(cfg.Timezones) != 2 {
		t.Errorf("Expected 2 timezones, got %d", len(cfg.Timezones))
	}

	if cfg.Timezones[0] != "America/Chicago" {
		t.Errorf("First timezone = %q, want %q", cfg.Timezones[0], "America/Chicago")
	}

	if cfg.Theme != "gruvbox" {
		t.Errorf("Theme = %q, want %q", cfg.Theme, "gruvbox")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	origHome := os.Getenv("HOME")
	defer os.Setenv("HOME", origHome)

	tmpDir, err := os.MkdirTemp("", "worldclock-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".config", "worldclock")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.json")

	configContent := `{invalid json`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err = Load()
	if err == nil {
		t.Error("Load() with invalid JSON should return error")
	}
}
