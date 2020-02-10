# Demo

This example shows how to use the library to implement a simple pub/sub application.

## Running the example

The example requires a working Go development environment.

### Server

    $ go run main.go

Next command line options are allowed:

Option | Default | Description
--- | --- | ---
addr | localhost:8080 | Server host and port
path | / | Server URL
publish | 1s | Messages publishing interval
channels | general,public,private | Publishing channels

### WEB

1. Open `index.html` in your browser
2. Click **Connect** button
2. Click **Subscribe** button

Now, every published message which is matched your subscription will be shown in the **Messages** field.