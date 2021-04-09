package gdns

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	// GoogleDomainsUpdateAPI is the endpoint to update Dynamic DNS records
	GoogleDomainsUpdateAPI string = "https://domains.google.com/nic/update"
)

type Credentials struct {
	Username string
	Password string
}

func UpdateDNS(ctx context.Context, ip net.IP, name string, credentials *Credentials) error {
	if credentials == nil {
		return fmt.Errorf("Missing credentials")
	}

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodPost, GoogleDomainsUpdateAPI, nil)
	req.SetBasicAuth(credentials.Username, credentials.Password)

	query := req.URL.Query()
	query.Add("hostname", name)
	query.Add("myip", ip.String())
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)

	if err != nil {
		if body, errBody := ioutil.ReadAll(resp.Body); errBody != nil {
			log.Errorf("Update unsuccessful %s", body)
		}
		log.Errorf("Error while updating: %s", err)
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Infof("Update successful %s", body)

	return nil
}
