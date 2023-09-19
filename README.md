# QuakeWorld demobot

> Setup for automated QuakeWorld client streaming demos, accepting commands via Twitch chat.

* **Visit [twitch.tv/QuakeWorldDemos](https://www.twitch.tv/QuakeWorldDemos)** to see it in action.

## How does it work? (TLDR version)

ezQuake reads from a pipe located at `/tmp/ezquake_[username]` on posix systems, where `username` is the username of the
user who started the ezQuake process.

So basically all you have to do is to write commands to `/tmp/ezquake_[username]`.

## Stack

* Written in [Go (Golang)](https://github.com/golang/go)
* [ZeroMQ](https://zeromq.org/) - Communication/messages (single proxy and multiple subscribers/publishers)

## Overview

![image](https://github.com/vikpe/qw-demobot/assets/1616817/5010507a-c773-4d26-a57b-92a015613fba)

* **Message Proxy**: Central point for communication.
* **Quake Manager**: Interaction with ezQuake
    * Log monitor (thread): Read in-game events (demo started, demo stopped, etc)
    * Process monitor (thread): ezQuake events (started, stopped)
* **Twitch Manager**: Interaction with Twitch channel (e.g. set title).
* **Twitch Bot**: Interaction with Twitch chat.

## Development

### Directory structure

Uses the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

```bash
cmd/       # Main applications
internal/  # Private application and library code
scripts/   # Various build, install operations
```

### Build

**Build specific app**

Example: build proxy

```shell
cd cmd/proxy
go build
```

**Build all apps**

```shell
./scripts/build.sh
```

### Run

**Single app**

Example: start the proxy.

```shell
./cmd/proxy/proxy 
```

**App controller scripts**

Runs app forever (restarts on error/sigint with short timeout in between).

```shell
bash scripts/controllers/proxy.sh
bash scripts/controllers/quake_manager.sh
bash scripts/controllers/twitch_manager.sh
bash scripts/controllers/twitch_chatbot.sh
bash scripts/controllers/ezquake.sh
```

### Test

```shell
go test ./... --cover
```

## Production

Build all apps and run all app controller scripts.

```shell
./scripts/build.sh && ./scripts/start.sh
```
