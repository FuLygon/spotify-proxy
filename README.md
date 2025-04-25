# Reverse proxy for Spotify API

Initially made this to use with my private fork of [NowPlaying for Spotify](https://github.com/busybox11/NowPlaying-for-Spotify) and my [Homepage](https://github.com/gethomepage/homepage) dashboard.

## Usage and Configuration

### Docker Installation

- Prepare and setting up the `.env` file:

```bash
wget https://raw.githubusercontent.com/fulygon/spotify-proxy/main/.env.example -O .env
```

- Optional, prepare and setting up the `routes.yaml` file:

```bash
wget https://raw.githubusercontent.com/fulygon/spotify-proxy/main/routes.example.yaml -O routes.yaml
```

- Then deploy the service. Compose file example:

```yaml
services:
  spotify-proxy:
    image: ghcr.io/fulygon/spotify-proxy:latest
    container_name: spotify-proxy
    env_file: .env
    ports:
      # Access server port
      - 8000:8000
      # Proxy server port
      - 8001:8001
    # Optional, only if you already configured routes.yaml
    volumes:
      - ./routes.yaml:/app/routes.yaml
```

### Source Installation

Make sure [Go](https://go.dev/doc/install) is installed.

- Clone the repo:

```bash
git clone https://github.com/FuLygon/spotify-proxy.git
cd spotify-proxy
```

- Prepare and setting up the `.env` file:

```bash
cp .env.example .env
```

- Optional, prepare and setting up the `routes.yaml` file:

```bash
cp routes.example.yaml routes.yaml
```

- Then build and run:

```bash
go build -o spotify-proxy ./cmd/main.go
./spotify-proxy
```

### Post Installation

The service runs on two ports by default:

- **8000**: Access port, where you can access the authentication API, should not be exposed to the public.
- **8001**: Reverse proxy port, all your requests to Spotify API will be proxied through this port.

On service startup, you will need to access http://127.0.0.1:8001 to log in, so the service can get and cache your refresh token, if the service is restarted, you will need to log in again.

If you set your refresh token in the `.env`file, then you're no longer need to log in. If you don't know how to get the refresh token, you can set `SPOTIFY_REFRESH_TOKEN_OUTPUT`, then access the login page, after logged in successfully, the refresh token will be printed to the service console. You can then set it in the `.env` file and restart the service.

By default, the proxy server will forward all routes to Spotify API (https://api.spotify.com), you can specify which routes to forward in the `routes.yaml` file.