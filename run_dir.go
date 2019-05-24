package piaas

import (
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas/util"
	"os"
	"path"
	"path/filepath"
)

var PidfileName = "pid"
var LogfileName = "output.log"

// RunDir is the directory where application log and pid file are kept.
// Usually .piaas.d
type RunDir string

// Prepare run dir for an app.
func (dir RunDir) prepare(app string) {
	p := path.Join(string(dir), app)
	err := os.MkdirAll(string(p), 0755)
	util.CheckFatal("prepare run dir", err)
}

func (dir RunDir) Pidfile(app string) string {
	dir.prepare(app)
	return path.Join(string(dir), app, PidfileName)
}

func (dir RunDir) Logfile(app string) string {
	dir.prepare(app)
	return path.Join(string(dir), app, LogfileName)
}

func NewRunDir(dir string) RunDir {
	cleaned := path.Clean(dir)
	err := os.MkdirAll(cleaned, 0755)
	util.CheckFatal("new run dir", err)
	evaluated, err := filepath.EvalSymlinks(cleaned)
	if err != nil {
		log.Errorf("Error evaluating symlink %s: %s.", cleaned, err)
		return RunDir(cleaned)
	} else {
		return RunDir(evaluated)
	}
}
