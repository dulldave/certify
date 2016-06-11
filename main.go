package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

var host = flag.String("host", "google.com", "The host to check against")
var every = flag.Int("every", 86400, "When to check the certificates in seconds, defaults to every day.")
var warnDay = flag.Int("day", 0, "How many days before the certifcate expires")
var warnMonth = flag.Int("month", 0, "How many month(s) before the certifcate expires")
var warnYear = flag.Int("year", 0, "How many year(s) before the certifcate expires")

var slackAPIKey = flag.String("slack-api-key", "abcd-1234-5678", "Your slack API key")
var slackChan = flag.String("slack-channel", "example", "The channel for the slack bot to post to")

func main() {
	flag.Parse()

	// Instantiate the Slack SDK
	api := slack.New(*slackAPIKey)
	params := slack.PostMessageParameters{Username: "Certify Bot"}

	t := time.NewTicker(time.Duration(*every) * time.Second)

	for range t.C {
		conn, err := tls.Dial("tcp", *host+":443", nil)
		if err != nil {
			logrus.WithFields(logrus.Fields{"Error": "Failed to connect to host"}).Error(err)
			continue
		}

		// Iterate over the chains and gather the certifcate data.
		// We use this to compare whether the certificate is getting closer
		// to the expiration date.
		for _, chain := range conn.ConnectionState().VerifiedChains {
			for _, cert := range chain {
				if time.Now().AddDate(*warnYear, *warnMonth, *warnDay).After(cert.NotAfter) {
					message := fmt.Sprintf("SSL/TLS Certificate Warning: %s. Expiration: %s.", cert.Subject.CommonName, cert.NotAfter)

					// Log the message to stdout
					logrus.WithFields(logrus.Fields{
						"Cert Name":  cert.Subject.CommonName,
						"Expiration": cert.NotAfter,
					}).Info(message)

					if _, _, err := api.PostMessage(*slackChan, message, params); err != nil {
						logrus.WithFields(logrus.Fields{
							"Error":   "Failed to post message to slack channel",
							"Message": message,
						}).Error(err)
					}
				}
			}
		}

		conn.Close()
	}
}
