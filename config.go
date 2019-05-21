package piaas

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/sohoffice/piaas/util"
	"io/ioutil"
	"os"
)

// The sync connection target
type SyncConnection struct {
	Host        string
	User        string
	Destination string
}

// Profile controls how the project can sync to remote machine.
type Profile struct {
	// profile name
	Name string
	// The connection information used by sync module
	Connection SyncConnection
	// The name of ignore file, default to .piaasignore
	IgnoreFile string
}

type App struct {
	Name    string
	Command string
}

type Config struct {
	ApiVersion string
	Rsync      string
	Profiles   []Profile
	Apps       []App
}

func ReadConfig(file string) Config {
	var config = Config{
		Rsync: "rsync",
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		println("Error reading config file:", file)
		os.Exit(1)
	}
	err = yaml.Unmarshal(data, &config)
	util.CheckFatal("read yaml", err)
	err = validateProfile(&config)
	util.CheckFatal("validate config yaml", err)

	return config
}

func validateProfile(configPtr *Config) error {
	for i := range configPtr.Profiles {
		prof := &configPtr.Profiles[i]
		if prof.IgnoreFile == "" {
			prof.IgnoreFile = ".piaasignore"
		}
		switch {
		case prof.Connection == SyncConnection{}:
			return fmt.Errorf("profile #%d (0 based) has no connection", i)
		case prof.Name == "":
			return fmt.Errorf("profile #%d (0 based) without name", i)
		}
	}
	return nil
}

func (conf *Config) GetProfile(name string) (Profile, error) {
	for _, p := range conf.Profiles {
		if p.Name == name {
			return p, nil
		}
	}
	return Profile{}, fmt.Errorf("can not find profile named '%s'", name)
}

// Find matching profile and use the information to compose the rsync target string.
// In this format: ${user}@${host}:${target}
func (conf *Config) GetSyncTarget(name string) (string, error) {
	prof, err := conf.GetProfile(name)
	if err != nil {
		return "", err
	}
	c := prof.Connection
	return fmt.Sprintf("%s@%s:%s", c.User, c.Host, c.Destination), nil
}
