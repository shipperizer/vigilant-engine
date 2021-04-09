package gdns

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	// IpifyAPI is the endpoint to get the external IP
	IpifyAPI string = "https://api.ipify.org"
)

type IpifyResponse struct {
	IP string `json:"ip"`
}

func FetchExternalIP() (*net.IP, error) {
	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet, IpifyAPI, nil)
	query := req.URL.Query()
	query.Add("format", "json")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)

	ipifyResp := IpifyResponse{}
	err = json.Unmarshal([]byte(data), &ipifyResp)

	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(ipifyResp.IP)
	return &ip, nil
}
