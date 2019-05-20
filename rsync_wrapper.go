package piaas

import (
	"fmt"
	"os/exec"
)

type RsyncWrapper struct {
	// The rsync command to execute, usually just 'rsync' on Mac and Linux.
	// But may involves batch file on windows.
	rsyncCmd   string
	basedir    string
	ignoreFile string
	syncTarget string
	sshOptions string

	syncCh chan *exec.Cmd
}

// Create a RsyncWrapper using `cmd` as the command, running on `basedir`, syncing to `target`.
// The RsyncWrapper will also be configured to use ssh options 'ConnectTimeout=10'
func NewRsyncWrapper(rsyncCmd string, basedir string, target string) RsyncWrapper {
	return RsyncWrapper{
		rsyncCmd:   rsyncCmd,
		basedir:    basedir,
		syncTarget: target,
		// Use 10 as the default connect timeout. Can be override.
		sshOptions: "ConnectTimeout=10",
	}
}

// Open the channel to start working on sync events
//
// When cmd is received, `process` will be invoked to actually running the command.
func (rs *RsyncWrapper) Start(process func(cmd *exec.Cmd)) {
	go func() {
		for {
			cmd := <-rs.syncCh
			process(cmd)
		}
	}()

	// The sync channel of a buffer size of 30.
	rs.syncCh = make(chan *exec.Cmd, 30)
}

func (rs *RsyncWrapper) SetIgnoreFile(ignore string) {
	rs.ignoreFile = ignore
}

func (rs *RsyncWrapper) SetSshOptions(options string) {
	rs.sshOptions = options
}

// Sync only the specified files
// If the files list is empty, do nothing
func (rs *RsyncWrapper) SyncFiles(files []string) {
	if len(files) <= 0 {
		return
	}
	var arguments = []string{"-av"}
	arguments = append(arguments, rs.getSshOptionsForRsync()...)
	arguments = append(arguments, rs.getExcludeFromForRsync()...)

	// Build the include file list
	arguments = append(arguments, []string{"--include='*/'"}...)
	for _, f := range files {
		arguments = append(arguments, fmt.Sprintf("--include='%s'", f))
	}
	arguments = append(arguments, []string{"--exclude='*'", "--delete", "--copy-links"}...)
	arguments = append(arguments, []string{".", rs.syncTarget}...)

	cmd := exec.Command(rs.rsyncCmd, arguments...)

	// Build command similar to this
	// Ex: rsync -av <ssh options> --exclude-from=...
	//     --include=... . foo@bar.com:~/src

	// send the commands to syncCh.
	rs.syncCh <- cmd
}

// Sync all files to remote.
func (rs *RsyncWrapper) SyncAll() {
	var arguments = []string{"-av"}
	arguments = append(arguments, rs.getSshOptionsForRsync()...)
	arguments = append(arguments, rs.getExcludeFromForRsync()...)
	arguments = append(arguments, []string{"--delete", "--copy-links", ".", rs.syncTarget}...)

	cmd := exec.Command(rs.rsyncCmd, arguments...)

	rs.syncCh <- cmd
}

func (rs *RsyncWrapper) getSshOptionsForRsync() []string {
	var sshOptions []string
	if rs.sshOptions != "" {
		sshOptions = []string{"-e", fmt.Sprintf("ssh -o %s", rs.sshOptions)}
	}
	return sshOptions
}

func (rs *RsyncWrapper) getExcludeFromForRsync() []string {
	var excludeFrom []string
	if rs.ignoreFile != "" {
		excludeFrom = []string{fmt.Sprintf("--exclude-from=%s", rs.ignoreFile)}
	}
	return excludeFrom
}
