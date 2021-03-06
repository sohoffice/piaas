Sync
====

Piaas sync is a wrapper around rsync. It will detect workspace changes and sync
necessary files to remote. Piaas sync is slightly more optimized than running
rsync with a file watcher for the following reasons:

1. Piaas detects changed files and update only the changed files.
2. Piaas checks through the ignore file, only necessary files are synced.

Piaas is also more suitable than using modern file syncing services like Dropbox.

1. Dropbox will sync all files, including the compiled files. 

   As you know, the compilation can happen on both ends of the sync.
   This will eventually create a lot of conflicts and headaches.
2. It's not easy to exclude files, if possible. 

   With Dropbox you can exclude folders, but not a single file. 
   The exclude is lost if the folder is deleted. Think about a `make clean`.
3. Dropbox is slightly slower.

   Everything will be synced. 
   If your machines are not in the same local network, they will be synced through server. (Through VPN for example)
   This is slightly slower than simply using piaas.

## Config

Create the following 2 files in your workspace.

##### piaasconfig.yml

This is the main configuration file, where you describe the connection
info of remote machines.

You may configure more than one machines, find samples of different platform
in the [document/samples](documents/samples) folder.

An example of piaasconfig.yml

```
apiVersion: 1
profiles:
  - name: dev
    connection:
      host: remote
      user: foo
      destination: ~/src
```

###### Windows

The piaas configuration on windows platform is trickier than the other.
Find the example in the below and note the executable section. This is
necessary to execute rsync through WSL. Using other rsync binary may work
but was never tested.

Windows example:

```
apiVersion: 1
executable:
  cmd: wsl
  params: ["rsync"]
profiles:
  - name: dev
    connection:
      host: remote.host.name
      user: foo
      destination: ~/src

```

##### .piaasignore

The ignore file uses a very similar format as rsync exclude files. The
following notations are supported.

- Exact, path element without special character. Ex: `foo`
- Anchored, path start with "/". Ex: `/foo`
- Wildcard, path element that matches any character. Ex: `foo*`
- Double wildcards, path elements that matches any character including 
  sub directories. Ex: `foo/**`
- Multi segments, path elements that span over one level in the hierarchy. 
  Ex: `foo/bar`

An example .piaasignore:

```
/.idea
/.git
*___jb_tmp___
*___jb_old___
*.class
```

## Prepare

To get piaas sync running, we need to make sure rsync is running. Rsync
uses ssh under the hood. The following are required.

- SSH key
- Remote machine was added in known_hosts
- Local public key is authorized in remote machine

To make sure the above work, follow the below steps

1. Make sure you have ssh key

   ```
   ls -al ~/.ssh
   # You should have id_rsa and id_rsa.pub files.
   ```
   
   If you do not have `id_rsa` and `id_rsa.pub`, generate the ssh key
    
   ```
   ssh-keygen
   ```
   
2. Connect to remote to make sure the connection works

   ```
   # user and remote_host are the connection info from piaasconfig.yml
   ssh user@remote_host
   ```
   
   You may be prompted with the below
   
   ```
   The authenticity of host 'remote_host (10.1.1.1)' can't be established.
   ECDSA key fingerprint is SHA256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.
   Are you sure you want to continue connecting (yes/no)?
   ```
   
   Say `yes` to remember the fingerprint.

3. Authorize yourself with remote machine

   ```
   ssh-copy-id user@remote_host
   ```
   
   The command will ask you to input password. Upon successful authenticated,
   you'll be authorized to login to remote host without having to input
   password each time.
   
   ###### No ssh-copy-id ?
   
   If your system do not have ssh-copy-id, use your package manager to install
   it.
   
   ```
   Mac OS: brew install ssh-copy-id
   ``` 

##### Windows

Piaas sync is using WSL (windows subsystem linux). You need to make sure
WSL is enabled and execute the above command in WSL bash. 

## Running

Use the `sync` command of piaas.

```
piaas sync <profile name>
```

The profile name is defined in your `piaasconfig.yml` file. You may defined 
multiple profiles, but sync can only work with one profile at a time.

Your workspace will keep sync with the remote server. You may now go to
the server to compile and serve from there. 

##### Windows

Piaas sync will invoke rsync through WSL, but it's recommended to run
it through normal windows command prompt. If you run it through WSL bash, the command may not be terminated by 
ctrl-c. You'll therefore have to close the entire command box.
