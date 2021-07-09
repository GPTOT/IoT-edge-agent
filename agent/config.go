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
// If the topic direction is `Pull` then produce to the source topic.
func (t Topic) produceTo() string {
	if t.direction == Push {
		return t.destinationName
	} else {
		return t.sourceName
	}
}

type SASLConfig struct {
	SaslMethod   string `koanf:"sasl_method"`
	SaslUsername string `koanf:"sasl_username"`
	SaslPassword string `koanf:"sasl_password"`
}

type TLSConfig struct {
	Enabled        bool   `koanf:"enabled"`
	ClientKeyFile  string `koanf:"client_key"`
	ClientCertFile string `koanf:"client_cert"`
	CaFile         string `koanf:"ca_cert"`
}

var defaultConfig = confmap.Provider(map[string]interface{}{
	"id":                            defaultID,
	"create_topics":                 false,
	"max_poll_records":              1000,
	"max_backoff_secs":              600, // ten minutes
	"source.name":                   "source",
	"source.bootstrap_servers":      "127.0.0.1:19092",
	"source.consumer_group_id":      defaultID,
	"destination.name":              "destination",
	"destination.bootstrap_servers": "127.0.0.1:29092",
	"destination.consumer_group_id": defaultID,
}, ".")

// Returns the hostname reported by the kernel to use as the default ID for the
// agent and consumer group IDs
var defaultID = func() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Unable to get hostname from kernel. Set Id in config")
	}
	log.Debugf("Hostname: %s", hostname)
	return hostname
}()

// Initialize the agent configuration from the provided .yaml file
func InitConfig(path *string) {
	lock.Lock()
	defer lock.Unlock()

	config.Load(defaultConfig, nil)
	log.Infof("Init config from file: %s", *path)
	if err := config.Load(file.Provider(*path), yaml.Parser()); err != nil {
		log.Errorf("Error loading config: %v", err)
	}
	validate()
	log.Debugf(config.Sprint())
}

// Parse topic configuration
func parseTopics(topics []string, direction Direction) []Topic {
	var all []Topic
	for _, t := range topics {
		s := strings.Split(t, ":")
		if len(s) == 1 {
			all = append(all, Topic{
				sourceName:      strings.TrimSpace(s[0]),
				destinationName: strings.TrimSpace(s[0]),
				direction:       direction,
			})
		} else if len(s) == 2 {
			// Push from source topic to destination topic
			var src = strings.TrimSpace(s[0])
			var dst = strings.TrimSpace(s[1])
			if direction == Pull {
				// Pull from destination topic to source topic
				src = strings.TrimSpace(s[1])
