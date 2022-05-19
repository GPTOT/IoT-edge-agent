#!/bin/bash

HELPTEXT='redpanda-edge-agent build script
  Usage: ./build [flags]

  Flags:
  --archive  -a            Create a compressed archive file (using tar -czf)
  --build    -b PLATFORM   Build for a specific platform (where PLATFORM is linux/amd64, for example)
  --build-all              Build for the following platforms: darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
  --help     -h            Show this message
  --include-platform       Include platform in resulting filename (always enabled with --build-all)
  --build-version VERSION  Use the given version when naming the archive file
  --verbose  -v            Print task details
'

PLATFORMS=("`go env GOOS`/`go env GOARCH`")
DEFAULT_PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")
FILENAME="redpanda-edge-agent"
INCLUDE_PLATFORM=false
ARCHIVE=false
VERBOSE=false

while [ $# -gt 0 ]; do
  case $1 in
  -h | --help) echo "$HELPTEXT"; exit;;
  -a | --archive) ARCHIVE=true;shift;;
  -b | --build) PLATFORMS=("$2");shift 2;;
  --build-all) PLATFORMS=(${DEFAULT_PLATFORMS[*]});INCLUDE_PLATFORM=true;shift;;
  --build-version) VERSION=$2;shift 2;;
  --include-platform) INCLUDE_PLATFORM=true;shift;;
  -v | --verbose) VERBOSE=true;shift;;
  *) break
  esac
don