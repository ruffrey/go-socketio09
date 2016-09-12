package socketio09

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type handshakeResponse struct {
	token             string
	heartbeatTimeout  int
	connectionTimeout int
}

func handshake(fullURL string) (hr handshakeResponse, err error) {
	hr = handshakeResponse{}
	timeToken := strconv.Itoa(int(time.Now().Unix()))
	handshakeURL, _ := url.Parse(fullURL)
	handshakeURL.Query().Set("t", timeToken)

	resp, err := http.Get(handshakeURL.String())
	if err != nil {
		return hr, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return hr, err
	}
	handshakeResponse := string(body[:len(body)])

	handshakeParts := strings.Split(handshakeResponse, ":")
	hr.token = handshakeParts[0]
	hr.heartbeatTimeout, _ = strconv.Atoi(handshakeParts[1])
	hr.connectionTimeout, _ = strconv.Atoi(handshakeParts[2])

	return hr, nil
}
