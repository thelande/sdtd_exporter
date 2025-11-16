# 7 Days to Die Server Exporter

A [Prometheus](https://prometheus.io) exporter for exposing metrics from the
[7 Days to Die](https://7daystodie.com/) server API.

This README does not cover how to set up a 7 Days to Die server. For a good guide,
please take a look at the [CSMM server cookbook](https://docs.csmm.app/en/7D2D/).

## Configuration

The exporter uses command line flags and environment variables for configuration.
See the [Running](#running) section for details.

## Running

Running the `sdtd_exporter` command without arguments will cause it to
expose the metrics on port `9816` under `/metrics` and attempt to connect to a 7 Days to Die
server API running on the local system on port 8080.

Use the command line flags and environment variables to specify the server URL and API
credentials:

```console
$ ./sdtd_exporter --help
usage: sdtd_exporter [<flags>]


Flags:
  -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
      --server.url="http://127.0.0.1:8080"
                                 The base URL of the 7 Days to Die dedicated server web API (e.g., http://127.0.0.1:8080). ($SDTD_API_URL)
      --server.token-name=SERVER.TOKEN-NAME
                                 The name of the API token to use to authenticate with the web server. ($SDTD_TOKEN_NAME)
      --server.token-secret=SERVER.TOKEN-SECRET
                                 The secret of the API token to use to authenticate with the web server. ($SDTD_TOKEN_SECRET)
      --web.telemetry-path="/metrics"
                                 Path under which to expose metrics.
      --[no-]web.systemd-socket  Use systemd socket activation listeners instead of port listeners (Linux only).
      --web.listen-address=:9816 ...
                                 Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.
      --web.config.file=""       Path to configuration file that can enable TLS or authentication. See:
                                 https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
      --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
      --[no-]version             Show application version.
```

You will likely need to specify the `--server.url` value as well as the `SDTD_TOKEN_NAME` and `SDTD_TOKEN_SECRET`
environment variables for the values unique to your server:

```console
SDTD_TOKEN_NAME="test" SDTD_TOKEN_SECRET="aSecret" ./sdtd_exporter --server.url=http://1.2.3.4:8080
```

## Docker

A Docker image of the exporter is available at
[thelande/sdtd_exporter](https://hub.docker.com/r/thelande/sdtd_exporter) and takes the same
arguments as the command.

You should place your environment variables in an environment file named `.env`:

```shell
SDTD_TOKEN_NAME="<your token name>"
SDTD_TOKEN_SECRET="<your token secret>"
```

```console
docker run -d -p 9816:9816 --env-file=.env \
  thelande/sdtd_exporter --server.url="<your server URL>"
```
