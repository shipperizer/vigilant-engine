package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"

	"github.com/shipperizer/vigilant-engine/gdns"
)

// TODO @shipperizer use env vars and limit to 1 records
// move to vault/localstack eventually
type EnvSpec struct {
	DNSRecord   string `envconfig:"dns_record"`
	Username    string `envconfig:"username"`
	Password    string `envconfig:"password"`
	RefreshRate int    `envconfig:"refresh_rate"`
}

func main() {
	var specs EnvSpec
	err := envconfig.Process("", &specs)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.SetLevel(log.DebugLevel)

	cron := gdns.NewCron(
		&gdns.UpdateConfig{
			DNSRecords: []gdns.DNSRecord{
				gdns.DNSRecord{
					Name:     specs.DNSRecord,
					Username: specs.Username,
					Password: specs.Password,
				},
			},
		},
		time.Duration(specs.RefreshRate)*time.Hour,
	)

	go cron.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	// Block until we receive our signal.
	<-c

	cron.Shutdown()

	log.Info("Shutting down")
	os.Exit(0)
}
