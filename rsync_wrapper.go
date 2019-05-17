package piaas

import (
	"bytes"
	"fmt"
)

type RsyncWrapper struct {
	// The rsync command to execute, usually just rsync on Mac and Linux.
	// But may involves batch file on windows.
	cmd        string
	basedir    string
	ignoreFile string
	syncTarget string
	sshOptions string

	syncCh chan string
}

// Create a RsyncWrapper using `cmd` as the command, running on `basedir`, syncing to `target`.
// The RsyncWrapper will also be configured to use ssh options 'ConnectTimeout=10'
func NewRsyncWrapper(cmd string, basedir string, target string) RsyncWrapper {
	return RsyncWrapper{
		cmd:        cmd,
		basedir:    basedir,
		syncTarget: target,
		// Use 10 as the default connect timeout. Can be override.
		sshOptions: "ConnectTimeout=10",
	}
}

// Open the channel to start working on sync events
//
// When cmd is received, `process` will be invoked to actually running the command.
func (rs *RsyncWrapper) Start(process func(string)) {
	go func() {
		for {
			cmd := <-rs.syncCh
			process(cmd)
		}
	}()

	// The sync channel of a buffer size of 30.
	rs.syncCh = make(chan string, 30)
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
	var includeBuf bytes.Buffer
	// Build the include file list
	includeBuf.WriteString("--include='*/'")
	for _, f := range files {
		includeBuf.WriteString(fmt.Sprintf(" --include='%s'", f))
	}
	includeBuf.WriteString(" --exclude='*' --delete --copy-links")

	// Ex: rsync -av <ssh options> --exclude-from=...
	//     --include=... . foo@bar.com:~/src
	cmd := fmt.Sprintf("%s -av %s %s %s . %s", rs.cmd, rs.getSshOptionsForRsync(), rs.getExcludeFromForRsync(),
		includeBuf.String(), rs.syncTarget)

	// send the commands to syncCh.
	rs.syncCh <- cmd
}

// Sync all files to remote.
func (rs *RsyncWrapper) SyncAll() {
	cmd := fmt.Sprintf("%s -av %s %s --delete --copy-links . %s", rs.cmd, rs.getSshOptionsForRsync(), rs.getExcludeFromForRsync(), rs.syncTarget)
	rs.syncCh <- cmd
}

func (rs *RsyncWrapper) getSshOptionsForRsync() string {
	var sshOptions string
	if rs.sshOptions != "" {
		sshOptions = fmt.Sprintf("-e 'ssh -o %s'", rs.sshOptions)
	}
	return sshOptions
}

func (rs *RsyncWrapper) getExcludeFromForRsync() string {
	var excludeFrom string
	if rs.ignoreFile != "" {
		excludeFrom = fmt.Sprintf("--exclude-from=%s", rs.ignoreFile)
	}
	return excludeFrom
}
