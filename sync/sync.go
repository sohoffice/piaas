package sync

import (
	"fmt"
	"github.com/sohoffice/piaas"
	"github.com/sohoffice/piaas/util"
	"github.com/urfave/cli"
	"log"
	"path"
)

// Prepare the sync module.
// This usually involves setup the correct module flags.
func Prepare() cli.Command {
	return cli.Command{
		Name:    "sync",
		Aliases: []string{"s"},
		Usage:   "Synchronize local files to remote using rsync",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "config",
				Usage: "Specify the piaas config file name and path. Default to piaasconfig.yml",
				Value: "piaasconfig.yml",
			},
		},
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
	changes := make(chan []string)
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
	basedir := path.Clean(path.Dir("."))
	log.Println("Basedir:", basedir)
	log.Println("Rsync to:", syncTarget)
	log.Println("Piaas ignore file:", prof.IgnoreFile)

	monitor := piaas.NewRecursiveMonitor(basedir)
	ignore := piaas.NewRsyncPatterns(basedir)
	rsync := piaas.NewRsyncWrapper("rsync", basedir, syncTarget)
	rsync.Start(func(s string) {
		log.Println("Run:", s)
	})

	// subscribe to file system changes
	monitor.SubscribeToCollects(changes)
	// starting watching for file system changes
	monitor.Watch(1000)
	for {
		changed := <-changes
		// make sure the files are not excluded by the ignore rules.
		filtered := make(util.StringSet, 0)
		for _, s := range changed {
			if ignore.Match(s) == false {
				filtered = *filtered.Add(s)
			}
		}
		// After filtering, some files should be synced. Trigger rsync.
		if len(filtered) > 0 {
			rsync.SyncFiles(filtered)
		}
	}
}
