
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"

	log "github.com/sirupsen/logrus"
)

type Redpanda struct {
	name   string
	prefix Prefix
	topics []Topic
	client *kgo.Client
	adm    *kadm.Client
}

var (
	source          Redpanda
	sourceOnce      sync.Once
	destination     Redpanda
	destinationOnce sync.Once
	wg              sync.WaitGroup
)

// Closes the source and destination client connections
func shutdown() {
	log.Infoln("Closing client connections")
	source.adm.Close()
	source.client.Close()
	destination.adm.Close()
	destination.client.Close()
}

// Creates new Kafka and Admin clients to communicate with a cluster.
//
// The `prefix` must be set to either `Source` or `Destination` as it
// determines what settings are read from the configuration.
//
// The topics listed in `source.topics` are the topics that will be pushed by
// the agent from the source cluster to the destination cluster.
//
// The topics listed in `destination.topics` are the topics that will be pulled
// by the agent from the destination cluster to the source cluster.
func initClient(rp *Redpanda, mutex *sync.Once, prefix Prefix) {
	mutex.Do(func() {
		var err error
		name := config.String(
			fmt.Sprintf("%s.name", prefix))
		servers := config.String(
			fmt.Sprintf("%s.bootstrap_servers", prefix))

		topics := GetTopics(prefix)
		var consumeTopics []string
		for _, t := range topics {
			consumeTopics = append(consumeTopics, t.consumeFrom())
			log.Infof("Added %s topic: %s", t.direction.String(), t.String())
		}

		group := config.String(
			fmt.Sprintf("%s.consumer_group_id", prefix))

		opts := []kgo.Opt{}
		opts = append(opts,
			kgo.SeedBrokers(strings.Split(servers, ",")...),
			// https://github.com/redpanda-data/redpanda/issues/8546
			kgo.ProducerBatchCompression(kgo.NoCompression()),
		)
		if len(topics) > 0 {
			opts = append(opts,
				kgo.ConsumeTopics(consumeTopics...),
				kgo.ConsumerGroup(group),
				kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
				kgo.SessionTimeout(60000*time.Millisecond),
				kgo.DisableAutoCommit(),
				kgo.BlockRebalanceOnPoll())
		}
		maxVersionPath := fmt.Sprintf("%s.max_version", prefix)
		if config.Exists(maxVersionPath) {
			opts = MaxVersionOpt(config.String(maxVersionPath), opts)
		}
		tlsPath := fmt.Sprintf("%s.tls", prefix)
		if config.Exists(tlsPath) {
			tlsConfig := TLSConfig{}
			config.Unmarshal(tlsPath, &tlsConfig)
			opts = TLSOpt(&tlsConfig, opts)
		}
		saslPath := fmt.Sprintf("%s.sasl", prefix)
		if config.Exists(saslPath) {
			saslConfig := SASLConfig{}
			config.Unmarshal(saslPath, &saslConfig)
			opts = SASLOpt(&saslConfig, opts)
		}

		rp.name = name
		rp.prefix = prefix
		rp.topics = topics
		rp.client, err = kgo.NewClient(opts...)
		if err != nil {
			log.Fatalf("Unable to load client: %v", err)
		}
		// Check connectivity to cluster
		if err = rp.client.Ping(context.Background()); err != nil {
			log.Errorf("Unable to ping %s cluster: %s",
				prefix, err.Error())
		}

		rp.adm = kadm.NewClient(rp.client)
		brokers, err := rp.adm.ListBrokers(context.Background())
		if err != nil {
			log.Errorf("Unable to list brokers: %v", err)
		}
		log.Infof("Created %s client", name)
		for _, broker := range brokers {
			brokerJson, _ := json.Marshal(broker)
			log.Debugf("%s broker: %s", prefix, string(brokerJson))
		}
	})
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}