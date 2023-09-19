# QuakeWorld demobot

> Setup for automated QuakeWorld client streaming demos, accepting commands via Twitch chat.

* :link:**Visit [twitch.tv/QuakeWorldDemos](https://www.twitch.tv/QuakeWorldDemos)** to see it in action.

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
