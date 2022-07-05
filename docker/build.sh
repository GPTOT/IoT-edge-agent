#!/bin/bash

HELPTEXT='redpanda-edge-agent docker build script
  Usage: ./build [flags]

  Flags:
  --username  -u USERNAME  The Docker username for pushing to the registry
  --no-push                The script will not push the build to the registry
  --tag       -t TAG       The build tag (default "latest")
'

FILENAM