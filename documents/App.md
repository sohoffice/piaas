App
====

Piaas app is a wrapper around any arbitrary command. It allows you to 
start your command in the background. Piaas app is slightly better than 
running the command directly for the following reasons.

- It keeps the command running in the background. Normally you would need
  to use screen or nohup to do the same.
- You can manage your app in a unified way:
  - Start with `piaas app run [app name]`.
  - Tail the log with `piaas app log [app name]`.
  - Stop with `piaas app stop [app name]`.
  - Check status with `piaas app status` or simply `piaas app`.

## Config

Create or update the following file in your workspace.

##### piaasconfig.yml

This is the main configuration file, where you describe the connection
info of remote machines. It can also be used to define apps.

When configured with multiple apps, you'll have to specify the app name 
as a target. However, if you only have one app, you may skip it.

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
piaas app <sub command> [app name]
```

The sub commands are:

- run, Run the specified app in the background
- log, Tail the log of the specified app
- stop, Stop the app
- status, Print the status of all apps

The app name is defined in your `piaasconfig.yml` file. You may defined 
multiple apps, but the command work with only one app at a time. With 
only app defined, you may skip the app name, and the sole app will be 
used.

