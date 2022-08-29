# Edge Agent Example

This [Docker Compose](.docker-compose-redpanda.yaml) file spins up a local environment for testing the agent. The environment starts two containers:

- `redpanda-source`: simulates an IoT device that runs a single-node Redpanda instance and the agent to store and forward messages
- `redpanda-destination`: simulates a central Redpanda cluster that aggregates messages from all of the IoT devices

## Prerequisites

1. The agent communicates with the source and destinations clusters over TLS enabled interfaces, so before starting the containers, run [generate-certs.sh](./generate-certs.sh) to create the necessary certificates. The resulting `./certs` directory is mounted on the containers in the [compose](./compose.yaml) configuration. *Note: On Linux, it may be necessary to run `sudo chmod -R 777 certs` to ensure this directory can be read by the container user `redpanda`.*
2. Run the [build.sh](./build.sh) script to compile the agent for the Linux-based container

## Start the containers

```bash
cd example
docker-compose -f docker-compose-redpanda.yaml up -d
[+] Running 3/3
 ⠿ Network example_redpanda_network  Created
 ⠿ Container redpanda_source         Started
 ⠿ Container redpanda_destination    Started
```

## Test the agent

Open a new terminal and produce some messages to the source's `telemetryB` topic (note that the example [agent](./agent.yaml) is configured to create the topics on startup):

```bash
export REDPANDA_BROKERS=localhost:19092
for i in {1..60}; do echo $(cat /dev/urandom | head -c10 | base64) | rpk topic produce telemetryB; sleep 1; done
```

The agent will forward the messages to a topic with the same name on the destination. Open a second terminal and consume the messages:

```bash
export REDPANDA_BROKERS=localhost:29092
rpk topic consume telemetryC
{
  "topic": "telemetryC",
  "key": "a0f1fd421b85",
  "value": "aZ7NEkkd977GXQ==",
  "timestamp": 1674753398569,
  "partition": 0,
  "offset": 0
}
{
  "topic": "telemetryC",
  "key": "a0f1fd421b85",
  "value"