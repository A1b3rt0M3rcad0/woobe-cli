package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Context struct {
	APIURL            string `json:"api_url"`
	Workspace         string `json:"workspace_id,omitempty"`
	Project           string `json:"project_id,omitempty"`
	Credential        string `json:"credential,omitempty"`
	RuntimeCredential string `json:"runtime_credential,omitempty"`
}
type Config struct {
	Version  int                `json:"version"`
	Current  string             `json:"current,omitempty"`
	Contexts map[string]Context `json:"contexts"`
}

func DefaultPath() string {
	if s := os.Getenv("WOOBE_CONFIG"); s != "" {
		return s
	}
	d, err := os.UserConfigDir()
	if err != nil {
		d = "."
	}
	return filepath.Join(d, "woobe", "config.json")
}
func Load(path string) (Config, error) {
	c := Config{Version: 1, Contexts: map[string]Context{}}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return c, nil
	}
	if err != nil {
		return c, err
	}
	if err = json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	if c.Version != 1 {
		return c, fmt.Errorf("unsupported config version %d", c.Version)
	}
	if c.Contexts == nil {
		c.Contexts = map[string]Context{}
	}
	return c, nil
}
func Save(path string, c Config) error {
	b, e := json.MarshalIndent(c, "", "  ")
	if e != nil {
		return e
	}
	return AtomicWrite(path, b, 0600)
}
func AtomicWrite(path string, b []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".woobe-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if err = f.Chmod(mode); err == nil {
		_, err = f.Write(b)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(name, path)
}
