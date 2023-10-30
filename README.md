
<div align="center">
  <img width="15%" src="./redpanda_lab1.png" />
  <br />
  This is a <a href="https://github.com/GPTOT/IoT-edge-agent">IoT Edge Agent</a> project
</div>

# IoT Edge Agent

An efficient IoT agent that collaborates with edge devices to forward events to a central Kafka API compatible cluster. This agent is written in Go and employs the [franz-go](https://github.com/twmb/franz-go) Kafka client library.

<p align="center">
  <img width="50%" src="./redpanda_iot.png" />
</p>

# Build the agent

Build the agent for any Go-supported target platform by setting the `GOOS` and `GOARCH` variables. Refer to the `go tool dist list` command to get a full list of supported architectures.

```shell
go clean

# MacOS (Intel)
env GOOS=darwin GOARCH=amd64 go build -a -v -o IoT-edge-agent ./agent

# MacOS (M1)
env GOOS=darwin GOARCH=arm64 go build -a -v -o IoT-edge-agent ./agent

# Linux (x86_64)
env GOOS=linux GOARCH=amd64 go build -a -v -o IoT-edge-agent ./agent

# Linux (Arm)
env GOOS=linux GOARCH=arm64 go build -a -v -o IoT-edge-agent ./agent
```

Please note that the [build.sh](./build.sh) script builds the agent and adds the resulting executable to a tarball for the platforms listed above.

# Running the agent

```shell
Usage of ./IoT-edge-agent:
  -config string
    path to agent config file (default "agent.yaml")
  -loglevel string
    logging level (default "info")
```

# Configuration

Example `agent.yaml` can be found below:

```yaml
# The unique identifier for the agent. If not specified the id defaults to the
# hostname reported by the kernel. When forwarding a record, if the record's
# key is empty, then it is set to the agent id to ensure that all records sent
# by the same agent are routed to the same destination topic partition.
id: