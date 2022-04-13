#!/bin/bash

HELPTEXT='redpanda-edge-agent build script
  Usage: ./build [flags]

  Flags:
  --archive  -a            Create a compressed archive file (using tar -czf)
  --build    -b PLATFORM   Build for a specific platform (where PLATFORM is linux/amd64, for example)
  --build-all              Build for the following platforms: darwin/amd64, d