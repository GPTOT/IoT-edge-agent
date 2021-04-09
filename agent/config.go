package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kversion"
	"github.com/twmb/franz-go/pkg/sasl/aws"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/tlscfg"

	log "github.com/sirupsen/logrus"
)

const schemaTopic = "_schemas"

var (
	lock   = &sync.Mutex{}
	config = koanf.New(".")
)

// Configuration prefix
type Prefix string

const (
	Source      Prefix = "source"
	Destination Prefix = "destination"
)

type Direction int8

const (
	Push Direction = iota // Push from source topic to destination topic
	Pull                  // Pull from destination topic to source topic
)

fu