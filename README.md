# dinner-elo

[![coverage](https://raw.githubusercontent.com/USA-RedDragon/dinner-elo/main/.github/badges/coverage.svg)](https://github.com/USA-RedDragon/dinner-elo/actions)

Rank dinners with Elo ratings.

## Configuration

Configuration is read from `config.yaml` in the working directory (or the file
given with `--config`), environment variables, and command-line flags. See
[`config.example.yaml`](config.example.yaml) for an example file.

<!-- configulator:begin -->

| Key                    | Type           | Default                                       | Environment            | Flag                     | Description                                                           |
|------------------------|----------------|-----------------------------------------------|------------------------|--------------------------|-----------------------------------------------------------------------|
| `log-level`            | string         | `info`                                        | `LOG_LEVEL`            | `--log-level`            | Logging level for the application. One of debug, info, warn, or error |
| `storage.type`         | string         | `sqlite`                                      | `STORAGE_TYPE`         | `--storage.type`         | Storage type. One of mysql, postgres, sqlite                          |
| `storage.dsn`          | string         | `file:./.dinner-elo.db?cache=shared&mode=rwc` | `STORAGE_DSN`          | `--storage.dsn`          | Data source name for the storage                                      |
| `http.url`             | string         |                                               | `HTTP_URL`             | `--http.url`             | URL where the HTTP server is deployed, used for redirects. Required   |
| `http.address`         | string         |                                               | `HTTP_ADDRESS`         | `--http.address`         | IP address to listen on. Empty listens on all interfaces              |
| `http.port`            | integer        | `8080`                                        | `HTTP_PORT`            | `--http.port`            | Port to listen on                                                     |
| `http.trusted-proxies` | list of string |                                               | `HTTP_TRUSTED_PROXIES` | `--http.trusted-proxies` | Trusted proxies for the HTTP server                                   |
| `metrics.enabled`      | boolean        |                                               | `METRICS_ENABLED`      | `--metrics.enabled`      | Enable metrics server                                                 |
| `metrics.address`      | string         |                                               | `METRICS_ADDRESS`      | `--metrics.address`      | Address to listen on                                                  |
| `metrics.port`         | integer        | `9000`                                        | `METRICS_PORT`         | `--metrics.port`         | Port to listen on                                                     |
| `pprof.enabled`        | boolean        |                                               | `PPROF_ENABLED`        | `--pprof.enabled`        | Enable pprof server                                                   |
| `pprof.address`        | string         |                                               | `PPROF_ADDRESS`        | `--pprof.address`        | Address to listen on                                                  |
| `pprof.port`           | integer        | `9999`                                        | `PPROF_PORT`           | `--pprof.port`           | Port to listen on                                                     |
| `auth.jwt-secret`      | string         |                                               | `AUTH_JWT_SECRET`      | `--auth.jwt-secret`      | JWT secret for signing tokens. Required (secret)                      |
| `auth.client-id`       | string         |                                               | `AUTH_CLIENT_ID`       | `--auth.client-id`       | OAuth2 client ID                                                      |
| `auth.client-secret`   | string         |                                               | `AUTH_CLIENT_SECRET`   | `--auth.client-secret`   | OAuth2 client secret (secret)                                         |
| `auth.token-url`       | string         |                                               | `AUTH_TOKEN_URL`       | `--auth.token-url`       | OAuth2 token URL, e.g. https://example.com/oauth2/token               |
| `auth.user-url`        | string         |                                               | `AUTH_USER_URL`        | `--auth.user-url`        | OAuth2 user information URL, e.g. https://example.com/oauth2/userinfo |

<!-- configulator:end -->
