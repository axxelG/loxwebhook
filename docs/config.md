# Config

## Config settings

| Property            | Description | Default |
|---------------------|-------------|---------|
| ConfigFile          | Path and filename to the config file     | OS dependent |
| LogFileMain         | Path and filename to the main log file   | stdout |
| LogFileHTTPError    | Path and filename to the HTTP error log  | stdout |
| LogFileHTTPAccess   | Path and filename to the HTTP access log | stdout |
| ControlsFiles       | Path of the directory containing controls files | OS dependent |
| ListenPort          | Local TCP port where loxwebhook will listen. You can choose any [valid](https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers) and free local port as long as loxwebhook is reachable on port 443 from the public internet.  | 443 |
| PublicURI           | URI (host and domain) where loxwebhook will be reachable on the public internet | none |
| LetsEncryptCache    | Path of the directory where we will store the Let's Encrypt cache. It's important to keep the cache during restarts to avoid hitting Let's Encrypt [rate limits](https://letsencrypt.org/docs/rate-limits/) | `./cache/letsencrypt` |
| MiniserverURL       | URL to reach the Loxone Miniserver including protocol and port `http://192.168.123.1:80` | none |
| MiniserverUser      | Username to access the Loxone Miniserver | `admin` |
| MiniserverPassword  | Password to access the Loxone Miniserver | `admin` |
| MiniserverTimeout   | Timeout (seconds) for requests to Loxone Miniserver | 2 |

## Set config values

You can set config values in a config file, set environment variables or set flags when you start loxwebhook.

Settings in environment variables will overwrite settings in a config file and settings given by flags will overwrite both.

## Config file

You can use `config.example.toml` ([online](https://github.com/axxelG/loxwebhook/blob/master/config.example.toml)) as a starting point

The default location for the config file depends on the operating system.

- Windows: `.\config.toml`
- Other (Linux): `/etc/loxwebhook/config.toml`

You can set a custom file location via environment variable or by flag

- Environment variable: Set `LOXWEBHOOK_CONFIG` to `<path>/<filename>`
- Flag: Call `loxwebhook --config <path>/<filename>`

## Environment variables

- Must be prefixed with `LOXWEBHOOK_`
- Must be all uppercase

Examples:

- Windows:
  - cmd `set LOXWEBHOOK_LISTENPORT=1234`
  - powershell `$env:LOXWEBHOOK_LISTENPORT = 1234`
- Linux: `export LOXWEBHOOK_LISTENPORT=1234`

## Flags

Use `loxwebhook -h` to get a list with all possible flags