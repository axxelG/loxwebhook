# Loxwebhook

Make selected Loxone virtual inputs securely available to any service that can send http requests.

## State of development

Loxwebhook is in early development but it works stable for me without major issues like crashing service or messed up request. There are still some essential features missing so expect bigger changes that might change the request path, authKey handling or config values. All changes will be documented in the change log.

## Features

- HTTPS encryption (LetsEncrypt)
- AuthKey authorization

## Overview example setup

![Overview example setup](/readme_images/OverviewExampleSetup.svg)

## Background

The Loxone Miniserver is able to accept http requests and many services are able to send them. This would be a nice, flexible and de facto standard way to connect different services to the Miniserver. Unfortunately, due to limited hardware ressources, the Miniserver only supports http basic auth over unencrypted connections and a quite complicated "hand crafted" encryption that no service I know supports.

Loxwebhook runs on a seperate server (a very low level device like a Raspberry Pi is more than sufficient) to offload the http encryption, protect the Miniserver against DOS attacks (Rate limit the requests) and adds an authentication layer based on authKeys to authorize request.

## Targets

- [x] Config from config files, environment variables, and flags
- [x] HTTPS encryption with LetsEncrypt
- [x] AuthKey authorization
- [x] Support vitual inputs
- [ ] Support other inputs
- [x] Rate limiting
- [x] Run as systemd service
- [ ] Run as Windows service
- [x] Provide binaries for Raspberry
- [x] Provide binaries for x86 Linux
- [x] Provide binaries for Windows
- [x] Provide .deb packages for Raspian
- [ ] Encrypt Miniserver Communication and get rid of basic auth