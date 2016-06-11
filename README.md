Certify
=======

#### What does it do?

Certify is a bog-standard, simple, SSL/TLS certificate checker. It sprung up after someone let the expiration of our certificates slip by causing a not-so-fun day for integration partners.

It simply pings the given host, pulls down the certificate and checks how far away from the expiration date the certificate is.

If the expiration is within the provided parameters a notification will be sent via Slack.

#### How to run:g
```
certify -host google.com -day 30 -slack-api-key your-api-key-here -slack-channel slack-channel-name
```
