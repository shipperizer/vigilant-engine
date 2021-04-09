package gdns

import (
	"context"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Cron interface {
	Start()
	Shutdown()
	Refresh(cfg *UpdateConfig) error
}

type DNSRecord struct {
	Name     string
	Username string
	Password string
}

type UpdateConfig struct {
	DNSRecords []DNSRecord
	IP         net.IP
}

type cron struct {
	ticker   *time.Ticker
	shutdown chan bool
	updates  chan *UpdateConfig
	config   *UpdateConfig
	muConfig sync.RWMutex
}

func (c *cron) Start() {
	for {
		select {
		case <-c.shutdown:
			break
		case tick := <-c.ticker.C:
			log.Debugf("Tick at %v", tick)
			c.Refresh(c.config)
		}
	}
}

func (c *cron) Shutdown() {
	c.shutdown <- true
}

func (c *cron) Refresh(cfg *UpdateConfig) error {
	ip, err := FetchExternalIP()
	if err != nil {
		log.Errorf("Failed fetching ip, skipping tick: %s", err)
		return err
	}
	c.muConfig.Lock()
	c.config = cfg
	// update External IP
	c.config.IP = *ip
	c.muConfig.Unlock()
	c.dnsUpdates()
	return nil
}

func (c *cron) dnsUpdates() {
	var wg sync.WaitGroup
	records := make(chan DNSRecord)

	c.muConfig.RLock()
	ip := c.config.IP
	c.muConfig.RUnlock()

	for i := 0; i < 5; i++ {
		go func(r chan DNSRecord) {
			dnsRecord := <-r
			log.Infof("Updating record %s", dnsRecord.Name)
			err := UpdateDNS(
				context.Background(),
				ip,
				dnsRecord.Name,
				&Credentials{Username: dnsRecord.Username, Password: dnsRecord.Password},
			)

			if err != nil {
				log.Error(err)
			}
			wg.Done()
		}(records)
	}

	c.muConfig.RLock()
	for _, dnsRecord := range c.config.DNSRecords {
		records <- dnsRecord
		wg.Add(1)
	}
	c.muConfig.RUnlock()

	wg.Wait()
}

func NewCron(config *UpdateConfig, tickerRate time.Duration) Cron {
	return &cron{
		config:   config,
		updates:  make(chan *UpdateConfig, 1),
		shutdown: make(chan bool),
		ticker:   time.NewTicker(tickerRate),
	}
}
