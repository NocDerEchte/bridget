<h1 align="center"> Bridget </h1>

<h4 align="center">A simple application that "bridges" data from subscribed mqtt-topics to a database.</h4>

<p align="center">
    <a href="https://github.com/NocDerEchte/bridget/actions/workflows/run-tests.yml"></a>
        <img src="https://github.com/NocDerEchte/bridget/actions/workflows/run-tests.yml/badge.svg">
    <a href="https://goreportcard.com/badge/github.com/nocderechte/bridget"></a>
        <img src="https://goreportcard.com/badge/github.com/nocderechte/bridget">
</p>

## Installation

```bash
# Clone repository
git clone https://github.com/NocDerEchte/bridget.git

# Build local docker image
make docker

# Or build binary in current directory
make build

# Cleanup binary
make clean
```

## Configuration Overview

For bridget to work it needs a YAML configuration file. By default bridget will look for `/etc/bridget/config.yml` but the config file location can by changed by setting the `BRIDGET_CONFIG` environment variable.
An example config can be found here:

```yaml
---
mqtt:
  host: "127.0.0.1"
  port: 1883
  timeoutMilliseconds: 100
  topics:
    apartment/living_room: 1
    apartment/kitchen: 0
database:
  influx:
    host: "127.0.0.1"
    port: 8181
    org: my-personal-org
    database: home-db
    authScheme: Token
    timeoutSeconds: 5
```

### Configuration keys

| Key                              |  Type  | Description                                                              | Required |
| -------------------------------- | :----: | ------------------------------------------------------------------------ | :------: |
| `mqtt.host`                      | string | IP or DNS name of the MQTT broker                                        |    ✔️     |
| `mqtt.port`                      | string | Port the MQTT broker is listening on                                     |    ✔️     |
| `mqtt.timeoutMilliseconds`       |  int   | Duration in milliseconds to wait when closing the connection             |    ✔️     |
| `mqtt.topics`                    |  map   | Map of MQTT topics to QoS level. Format: `table/measurement: qos_level`  |    ✔️     |
|                                  |        |                                                                          |          |
| `database.influx`                |  map   | InfluxDB configuration                                                   |    ✔️     |
| `database.influx.host`           | string | IP or DNS name of the InfluxDB server                                    |    ✔️     |
| `database.influx.port`           | string | Port InfluxDB is listening on (default: `8181`)                          |    ❌     |
| `database.influx.org`            | string | Name of the InfluxDB organization                                        |    ❌     |
| `database.influx.database`       | string | Name of the target database                                              |    ✔️     |
| `database.influx.authScheme`     | string | Authentication method — currently only `Token` is supported              |    ✔️     |
| `database.influx.timeoutSeconds` |  int   | Duration in seconds to wait till timeout when connecting to the database |    ✔️     |

### Environment variables

| Name             | Description                                            | Default                   |
| ---------------- | ------------------------------------------------------ | ------------------------- |
| `BRIDGET_CONFIG` | Path to the main YAML configuration file               | `/etc/bridget/config.yml` |
| `LOG_FORMAT`     | Format for log output — options: `json`, `text`        | `json`                    |
| `LOG_LEVEL`      | Minimum log level — options: debug, info, warn, error  | `info`                    |
| `INFLUXDB_TOKEN` | Authentication token used for writing data to InfluxDB | `none`                    |

## Running bridget

### Local binary
  
To run bridget directly on your machine:

- Make sure the configuration file is correctly set up.
- Run the binary with no flags:

```bash
./bridget
```

- You should see a startup success message if LOG_LEVEL si set to `debug` or `info`.
- Once initialized, bridget will begin bridging data from MQTT to your configured database.

### Docker (recommended)

To run bridget inside a container:

- Bind your local configuration directory to `/etc/bridget` inside the container.
- A sample compose.yml is available in the repository under `deployments/docker/compose.yml`

## Roadmap

- [ ] Add support for different databases
- [ ] configure max concurrent database writes

## Your Feedback is welcomed

As this is my first go project, I'd greatly appreciate any constructive feedback you have!
Feel free to open issues if you encounter a problem.
