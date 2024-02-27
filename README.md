# Steam Hour Booster UI

This is a web UI for [DrWarpMan's Steam Hour Booster](https://github.com/DrWarpMan/steam-hour-booster)

## Development

Prerequisises:
- go 1.21
- node v20
- [entr](https://github.com/eradman/entr) for live reloading
- Steam Hour Booster config file with a single user

Usage:
- Setup dependencies: `make setup'
- Build the app: `make build`
- Run the app: `make run`
- Run and rebuild for css/template updates: `make watch`
> [!NOTE]
> `make watch` needs entr

