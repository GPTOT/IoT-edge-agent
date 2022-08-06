# Edge Agent Example

This [Docker Compose](.docker-compose-redpanda.yaml) file spins up a local environment for testing the agent. The environment starts two containers:

- `redpanda-source`: simulates an IoT device that runs a single-node Redpanda instance and the agent to store and forward messages
- `redpanda-destination`: simulates a central Redpanda cluster that aggregates messages from all of the IoT devices

## Prerequisites

1. The agent communicates with the source and destinations clusters over TLS enabled interfaces, so before starting the containers, run [generate-certs.sh](./generate-certs.sh) to create the necessary certificates. The resulting `./certs` directory is mounted on the containers in the [compose](./compose.yaml) configuration. *Note: On Linux, it may be necessary to run `sudo chmod -R 777 certs` to ensure this directory can be read by the container user `redpanda`.*
2. Run the [build.sh](./build.sh) script to compile the agent for the Linux-based container

## Start the containers

``