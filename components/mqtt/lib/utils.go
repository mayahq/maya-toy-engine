package lib

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

//
// ParseOptionsURI parses URI in the following format:
// tcp://username:password@127.0.0.1:1883/topic?clientID=...&clean=true&qos=0
// Keep in mind # should be expressed as %23
//
// Returns client options, default topic, qos, error
//
func ParseOptionsURI(uri string) (*mqtt.ClientOptions, string, mqtt.QoS, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, "", mqtt.QOS_ZERO, err
	}

	opts := mqtt.NewClientOptions()
	if url.User != nil {
		if url.User.Username() != "" {
			opts.SetUsername(url.User.Username())
		}
		if p, ok := url.User.Password(); ok {
			opts.SetPassword(p)
		}
	}

	opts.AddBroker(fmt.Sprintf("%s://%s", url.Scheme, url.Host))

	clientID := url.Query().Get("clientID")
	if clientID == "" {
		hn, _ := os.Hostname()
		clientID = "mqtt-pub_" + strings.Split(hn, ".")[0] + strconv.Itoa(time.Now().Nanosecond())
	}
	opts.SetClientId(clientID)

	cleanSession := true
	if url.Query().Get("clean") != "" && url.Query().Get("clean") != "true" {
		cleanSession = false
	}
	opts.SetCleanSession(cleanSession)

	qos := mqtt.QOS_ZERO
	i, err := strconv.Atoi(url.Query().Get("qos"))
	if i == 1 {
		qos = mqtt.QOS_ONE
	} else if i == 2 {
		qos = mqtt.QOS_TWO
	}

	return opts, url.Path, qos, nil
}
