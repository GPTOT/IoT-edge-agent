
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