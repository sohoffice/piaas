App
====

Piaas app is a wrapper around any arbitrary command. It allows you to 
start your commands in the background. Piaas app is slightly better than 
running the command directly on your own.

- It keeps the command running in the background. Normally you would need
  to use screen or nohup to do the same.
- You can manage your app in a unified set of commands:
  - Start with `piaas app run [app name]`.
  - Tail the log with `piaas app log [app name]`.
  - Stop with `piaas app stop [app name]`.

## Config

Create or update the following file in your workspace.

##### piaasconfig.yml

This is the main configuration file, where you describe the connection
info of remote machines. It can also be used to define apps.

You can configure multiple apps. When configured with multiple apps, you'll 
have to specify the app name as a target. If only one app is configured,
you may skip that part in command line.

An example of piaasconfig.yml

```
apiVersion: 1
profiles:
  - name: dev
    connection:
      ...
apps:
  - name: web
    cmd: "ng"
    params: ["serve", "--ssl", "--host", "0.0.0.0", "--disable-host-check"]
```

## Running

Use the `app` command of piaas.

```
piaas app <subcommand> [app name]
```

The sub commands are:

- run, Start the specified app in the background
- log, Tail the log of the specified app
- stop, Stop the app
- status, Print the status of all apps

The app name is defined in your `piaasconfig.yml` file. You may defined 
multiple apps, but only one app at a time. If you have only app defined,
you do not have to supply the app name. The sole app will be used.

