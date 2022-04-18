# Panorama

Isometric world mapper for Minetest

![screenshot](https://user-images.githubusercontent.com/4698994/163820087-6473cbc4-b790-4e6d-9130-aedb5bf1eddf.png)

## Installation guide

*Note: Panorama started as our in-house mapper, and installation is
non-trivial as a result. If you're not comfortable with complicated
setups, check out [mapserver] instead!*

### Prerequisites

- PostgreSQL backend for your world
- Several gigabytes of disk space for tiles
- A decent CPU and about a gigabyte of RAM, depending on workload
- Recent [Go compiler][go]

### Step 1: Install Panorama server

Currently, there are no pre-built binaries, so you'll have to build
it yourself:

```sh
git clone https://github.com/lord-server/panorama
cd panorama
go build
```

### Step 2: Extract game data

Install [`panorama_api` mod][panorama_api] and enable it.
This mod dumps all info (besides game assets) required to render the
map to your world directory.

### Step 3: Configure Panorama

Copy `config.example.toml` to `config.toml`:
```sh
cp config.example.toml config.toml
```

Edit `config.toml` and configure paths to your world and game:

```toml
game_path = "/path/to/games/your_game"
world_path = "/path/to/worlds/your_world"
world_dsn = 'host=localhost port=5432 user=postgres password=pass dbname=world'
```

### Step 4: Do a full render

This command will perform an initial render. It'll take a lot of
time to finish, especially if your map is big.

```sh
./panorama --fullrender
```

### Step 5: Run Panorama in server mode

Now you can run this command to serve tiles from address specified in config
(`localhost:33333` by default) 

```sh
./panorama --serve
```

## License

MIT

[mapserver]: https://github.com/minetest-mapserver/mapserver
[go]: https://go.dev/
[panorama_api]: https://github.com/lord-server/panorama_api