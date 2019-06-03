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
	Name   string
	Cmd    string
	Params []string
}

type Executable struct {
	Cmd    string
	Params []string
}

type Config struct {
	ApiVersion string
	Executable
	Profiles []Profile
	Apps     []App
}

func ReadConfig(file string) Config {
	var config = Config{
		Executable: Executable{
			Cmd:    "rsync",
			Params: []string{},
		},
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

// Get app by name.
// Returned the only app, if name is empty string.
//
// Return error, if the app of the name can not be found. Or name is empty string but there're more than one app.
//
func (conf *Config) GetApp(name string) (App, error) {
	if name == "" && len(conf.Apps) == 1 {
		return conf.Apps[0], nil
	} else {
		for _, app := range conf.Apps {
			if app.Name == name {
				return app, nil
			}
		}
	}

	return App{}, fmt.Errorf("can not find app %s", name)
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
