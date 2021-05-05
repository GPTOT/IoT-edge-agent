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

func (d Direction) String() string {
	switch d {
	case Push:
		return "push"
	case Pull:
		return "pull"
	default:
		return fmt.Sprintf("%d", int(d))
	}
}

type Topic struct {
	sourceName      string
	destinationName string
	direction       Direction
}

func (t Topic) String() string {
	if t.direction == Push {
		return fmt.Sprintf("%s > %s",
			t.sourceName, t.destinationName)
	} else {
		return fmt.Sprintf("%s < %s",
			t.sourceName, t.destinationName)
	}
}

// Returns the name of the topic to consume from.
// If the topic direction is `Push` then consume from the source topic.
// If the topic direction is `Pull` then consume from the destination topic.
func (t Topic) consumeFrom() string {
	if t.direction == Push {
		return t.sourceName
	} else {
		return t.destinationName
	}
}

// Returns the name of the topic to produce to.
// If the topic direction is `Push` then produce to the destination topic.
// If the topic direction