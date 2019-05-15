package piaas

// The sync connection target
type SyncConnection struct {
	Host        string
	User        string
	Destination string
}

// Profile controls how the project can work with remote server.
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
	Profiles []Profile
	Apps     []App
}
