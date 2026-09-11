package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIURL     string `json:"api_url" yaml:"api_url"`
	UserAPIKey string `json:"user_api_key" yaml:"user_api_key"`
}

type configFile struct {
	Version       int               `yaml:"version"`
	ActiveProfile string            `yaml:"active_profile"`
	Profiles      map[string]Config `yaml:"profiles"`
}

const (
	envProfile = "UPROC_PROFILE"
)

var cliConfigPath string

func SetConfigPath(path string) {
	cliConfigPath = strings.TrimSpace(path)
}

func legacyConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "bizzmod-cli", "config.json"), nil
}

func resolveReadPath() (string, error) {
	if configured := strings.TrimSpace(cliConfigPath); configured != "" {
		absPath, err := filepath.Abs(configured)
		if err == nil {
			return absPath, nil
		}
		return configured, nil
	}
	absPath, err := filepath.Abs("config.yml")
	if err == nil {
		return absPath, nil
	}
	return "config.yml", nil
}

func ResolvedConfigPath() (string, error) {
	return resolveReadPath()
}

func resolveWritePath() (string, error) {
	if configured := strings.TrimSpace(cliConfigPath); configured != "" {
		absPath, err := filepath.Abs(configured)
		if err == nil {
			return absPath, nil
		}
		return configured, nil
	}
	absPath, err := filepath.Abs("config.yml")
	if err == nil {
		return absPath, nil
	}

	return "config.yml", nil
}

func ensureMigrated(writePath string) error {
	if _, err := os.Stat(writePath); err == nil {
		return nil
	}

	legacyPath, err := legacyConfigPath()
	if err != nil {
		return nil
	}

	b, err := os.ReadFile(legacyPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return nil
	}

	var legacy Config
	if err := json.Unmarshal(b, &legacy); err != nil {
		return nil
	}

	legacy = normalize(legacy)

	if legacy.APIURL == "" && legacy.UserAPIKey == "" {
		return nil
	}

	file := configFile{
		Version:       1,
		ActiveProfile: "default",
		Profiles: map[string]Config{
			"default": legacy,
		},
	}

	return writeConfigFile(writePath, file)
}

func loadConfigFile(path string) (configFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return configFile{}, err
	}

	var file configFile
	if err := yaml.Unmarshal(b, &file); err != nil {
		return configFile{}, err
	}

	if file.Version == 0 {
		file.Version = 1
	}
	if file.Profiles == nil {
		file.Profiles = map[string]Config{}
	}

	for name, cfg := range file.Profiles {
		file.Profiles[name] = normalize(cfg)
	}

	return file, nil
}

func writeConfigFile(path string, file configFile) error {
	if file.Version == 0 {
		file.Version = 1
	}
	if file.Profiles == nil {
		file.Profiles = map[string]Config{}
	}

	for name, cfg := range file.Profiles {
		file.Profiles[name] = normalize(cfg)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	b, err := yaml.Marshal(file)
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0o600)
}

func effectiveProfileName(file configFile) string {
	if name := strings.TrimSpace(os.Getenv(envProfile)); name != "" {
		return name
	}

	if name := strings.TrimSpace(file.ActiveProfile); name != "" {
		return name
	}

	return "default"
}

func normalize(cfg Config) Config {
	cfg.APIURL = strings.TrimRight(strings.TrimSpace(cfg.APIURL), "/")
	cfg.UserAPIKey = strings.TrimSpace(cfg.UserAPIKey)
	return cfg
}

func Load() (Config, error) {
	path, err := resolveReadPath()
	if err != nil {
		return Config{}, err
	}

	writePath, err := resolveWritePath()
	if err == nil {
		_ = ensureMigrated(writePath)
	}

	file, err := loadConfigFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("missing config file at %s (run `uproc processes login --profile default --use`)", path)
		}
		return Config{}, err
	}

	selected := effectiveProfileName(file)
	saved, ok := file.Profiles[selected]
	if !ok {
		return Config{}, fmt.Errorf("profile %q not found in %s", selected, path)
	}

	saved = normalize(saved)
	if err := Validate(saved); err != nil {
		return Config{}, fmt.Errorf("profile %q is invalid: %w", selected, err)
	}

	return saved, nil
}

func GetActiveProfileName() (string, error) {
	path, err := resolveReadPath()
	if err != nil {
		return "default", err
	}

	file, err := loadConfigFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "default", nil
		}
		return "default", err
	}

	return effectiveProfileName(file), nil
}

func LoadProfile(name string) (Config, error) {
	path, err := resolveReadPath()
	if err != nil {
		return Config{}, err
	}

	file, err := loadConfigFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
		}
		return Config{}, err
	}

	if strings.TrimSpace(name) == "" {
		name = effectiveProfileName(file)
	}

	cfg, ok := file.Profiles[name]
	if !ok {
		return Config{}, nil
	}

	return normalize(cfg), nil
}

func ListProfiles() ([]string, string, error) {
	path, err := resolveReadPath()
	if err != nil {
		return nil, "", err
	}

	file, err := loadConfigFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, "", nil
		}
		return nil, "", err
	}

	names := make([]string, 0, len(file.Profiles))
	for name := range file.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	return names, effectiveProfileName(file), nil
}

func SaveProfile(name string, cfg Config, setActive bool) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("profile name is required")
	}

	if strings.TrimSpace(cfg.APIURL) == "" || strings.TrimSpace(cfg.UserAPIKey) == "" {
		return fmt.Errorf("all login fields are required")
	}

	cfg = normalize(cfg)

	path, err := resolveWritePath()
	if err != nil {
		return err
	}

	file := configFile{Version: 1, Profiles: map[string]Config{}}
	loaded, err := loadConfigFile(path)
	if err == nil {
		file = loaded
	} else if !os.IsNotExist(err) {
		return err
	}

	if file.Profiles == nil {
		file.Profiles = map[string]Config{}
	}

	file.Profiles[name] = cfg
	if setActive || strings.TrimSpace(file.ActiveProfile) == "" {
		file.ActiveProfile = name
	}

	return writeConfigFile(path, file)
}

func SetActiveProfile(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("profile name is required")
	}

	path, err := resolveWritePath()
	if err != nil {
		return err
	}

	file, err := loadConfigFile(path)
	if err != nil {
		return err
	}

	if _, ok := file.Profiles[name]; !ok {
		return fmt.Errorf("profile %q not found", name)
	}

	file.ActiveProfile = name
	return writeConfigFile(path, file)
}

func Validate(cfg Config) error {
	if strings.TrimSpace(cfg.APIURL) == "" {
		return fmt.Errorf("missing UPROC_PROCESSES_API_URL or login config")
	}
	if strings.TrimSpace(cfg.UserAPIKey) == "" {
		return fmt.Errorf("missing USER_API_KEY or login config")
	}
	return nil
}
