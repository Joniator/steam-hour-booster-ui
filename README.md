# Steam Hour Booster UI

This is a web UI for [DrWarpMan's Steam Hour Booster](https://github.com/DrWarpMan/steam-hour-booster)
It is a UI for the config.json with basic container controls to restart, it does not interact with the booster itself.

![Screenshot of the app](https://raw.githubusercontent.com/Joniator/steam-hour-booster-ui/main/.github/screenshot.png)

## Roadmap

- [x] Add and remove boosted games by AppIDs
- [x] Show the names of the boosted games
- [x] Restart/Manage the container from web UI
- [ ] Add games by name with included search
- [ ] Show human readable logs of the container
- [ ] Disable boosted games without deleting them/remember previous games
- [ ] Track current boosted games, not just games in the config

## Usage

See the docker-compose.yml as a working example.
Configuring the docker container is optional, but recommended to reload the booster if the config changes.
The container does not get restartet on config changes automatically.
The `latest` Tag points to the latest release, edge gets rebuilt on main pushes and might be in a broken state.

To run it standalone: `steam-hour-booster-ui --config ./config.json --container steam_hour_booster`

## Development

Prerequisises:
- go 1.21
- node v20
- [entr](https://github.com/eradman/entr) for live reloading
- Steam Hour Booster config file with a single user

Usage:
- Setup dependencies: `make setup`
- Build the app: `make build`
- Run the app: `make run`
- Run and rebuild for css/template updates: `make watch`
> [!NOTE]
> `make watch` needs entr

