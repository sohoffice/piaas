package sync

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/sohoffice/piaas"
	"github.com/sohoffice/piaas/stringarrays"
	"github.com/sohoffice/piaas/util"
	"github.com/urfave/cli"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Prepare the sync module.
// This usually involves setup the correct module flags.
func Prepare() cli.Command {
	return cli.Command{
		Name:      "sync",
		Aliases:   []string{"s"},
		Usage:     "Synchronize local files to remote using rsync",
		Flags:     piaas.PrepareCommonFlags(),
		ArgsUsage: "<profile name>",
		Action:    Execute,
	}
}

// Execute the sync command.
//
// It will start monitoring current directory, running rsync if files are changed.
// The method will continue to run, until user press ctrl-c or other way to terminate it.
//
func Execute(c *cli.Context) error {
	if c.NArg() <= 0 {
		return fmt.Errorf("profile name is required")
	}
	err := piaas.HandleDebug(c)
	if err != nil {
		return err
	}
	collects := make(chan []string)
	profileName := c.Args().Get(0)
	config := piaas.ReadConfig(c.String("config"))
	prof, err := config.GetProfile(profileName)
	if err != nil {
		return err
	}
	syncTarget, err := config.GetSyncTarget(profileName)
	if err != nil {
		return err
	}
	basedir, err := filepath.Abs(path.Clean(path.Dir(".")))
	if err != nil {
		return err
	}
	basedir, err = filepath.EvalSymlinks(basedir)
	if err != nil {
		return err
	}
	log.Println("Basedir:", basedir)
	log.Println("Sync to:", syncTarget)
	log.Println("Ignore file: ", prof.IgnoreFile)

	monitor := piaas.NewMonitor(basedir)
	ignore, err := readIgnoreFile(basedir, prof.IgnoreFile)
	if err != nil {
		return err
	}
	rsync := piaas.NewRsyncWrapper(config.Executable, basedir, syncTarget)
	rsync.SetIgnoreFile(prof.IgnoreFile)
	rsync.Start(func(cmd *exec.Cmd) {
		log.Infof("Run: %s\n", cmd.Args)
		err := cmd.Run()
		if err != nil {
			log.Errorf("Error running rsync: %s.\n%s\n", err, cmd.Args)
		} else {
			log.Infof("Done.\n")
		}
	})
	// Trigger a sync all in the beginning.
	rsync.SyncAll()

	// subscribe to file system changes
	monitor.Subscribe(collects)
	// starting watching for file system changes
	monitor.Start(1000)
	for {
		collected := <-collects
		filtered := make(util.StringSet, 0)
		// make sure the files are not excluded by the ignore rules.
		for _, s := range collected {
			rel, err := filepath.Rel(basedir, s)
			if err != nil {
				log.Errorf("error getting relative path: %s.", err)
				rel = s
			}
			log.Debugf("Collected: %s, %s", s, rel)
			if ignore.MatchRelative(rel) == false {
				filtered = *filtered.Add(rel)
			} else {
				log.Debugf("  | Ignored: %s", rel)
			}
		}
		// After filtering, some files should be synced. Trigger rsync.
		if len(filtered) > 0 {
			log.Debugf("Detected file changes:\n%s", stringarrays.ToString(filtered))
			rsync.SyncFiles(filtered)
		}
	}
}

func readIgnoreFile(basedir string, ignorefile string) (piaas.RsyncPatterns, error) {
	f, err := os.Open(ignorefile)
	if err != nil {
		return piaas.RsyncPatterns{}, err
	}
	defer f.Close()

	patterns := make([]piaas.RsyncPattern, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		cleaned := strings.Trim(line, " \t")
		if len(cleaned) > 0 {
			log.Printf("  | %s", cleaned)
			patterns = append(patterns, piaas.NewRsyncPattern(cleaned))
		}
	}

	return piaas.NewRsyncPatterns(basedir, patterns...), nil
}
