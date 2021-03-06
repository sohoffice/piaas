# The Personal IAAS 

Piaas is a set of tools that help you **develop using multiple machines**.
You can use your favorite IDE on your Macbook, while compile and serve on a more
powerful desktop PC or cloud machine.

We all work on laptop today. However, compiling and serving on laptop
is crazy. The fan keeps working and creating quite a noise. Your battery
drained, cycle count increased. More than ever, peace of mind you lost.

Piaas help your workspace in sync with a remote machine, given it a 
more powerful desktop PC, an bare metal server or a cloud instance. Developing with
your favorite IDE on your laptop, running your api server and database server
from the cloud machine. Your laptop is therefore relaxed to do the thing that 
absolutely required it to do. **Great team work**, isn't it ? 

## Installation

Download the binary from the [github release page](https://github.com/sohoffice/piaas/releases/latest):

- Windows, 64 bit
- Mac OS, 64 bit
- Linux, 64 bit

Extract and place the binary in your path and you're done.

## System requirements

The application support Mac OS, but was only tested on Mojave.

Linux should be fine, but wasn't tested in the real world.

On windows platform, WSL (windows subsystem linux) is required. At the
moment **Windows 10** is the only version with WSL so is the only version
supported.

If you need instructions on how to enable WSL on windows, check [this](https://docs.microsoft.com/en-us/windows/wsl/install-win10).

## Getting started

1. Make sure you have rsync. You can test it via a command line prompt.
2. You need to allow your local user to access target account.

   This is usually done by adding your local public key to remote ~/.ssh/authorized_keys. Once set, use `ssh $user@$target_address` to make sure it actually works.

3. Setup piaasconfig.yml and .piaasignore

Check [sync document](documents/Sync.md) for more details

## Commands

- [sync](documents/Sync.md), Keep your local workspace in sync with a remote directory
- [app](documents/App.md), Run your application in the background and monitor the logs

