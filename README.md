# Loxwebhook

Make selected Loxone virtual inputs securely available to any service that can send http requests.

## Features

- HTTPS encryption (LetsEncrypt)
- AuthKey authorization

## Use case

The Loxone Miniserver is able to accept http requests and many services are able to send them. This would be a nice, flexible and de facto standard way to connect different services to the Miniserver. Currently the Loxone Miniserver supports two connection methods:

| Connection method    | Usability  |
|----------------------|------------|
| Unencrypted with [Basic Auth](https://en.wikipedia.org/wiki/Basic_access_authentication) | Because Basic Auth does not provide confidentiality it is not secure to use it on the public internet |
| Encrypted with custom encryption | No service I know supports it. |

Loxwebhook runs on a separate server (a very low level device like a Raspberry Pi is more than sufficient) to offload the http encryption, protect the Miniserver against DOS attacks (rate limit the requests) and adds an authentication layer based on authKeys to authorize request.

## State of development

Most probably Loxwebhook is used only by me or maybe a few others. Feel free to try it out.

There are no known bugs. Create a github issue if you find something not working as expected.

## Data flow example setup

![Overview data flow](/readme_images/DataFlowExampleSetup.svg)

## Overview example setup

![Overview example setup](/readme_images/OverviewExampleSetup.svg)

## Setup

### Prerequisites

A server with [supported operating systems](https://axxelg.github.io/loxwebhook/supported_os.html) on your local network reachable on port 443 from the public internet with working DNS resolution. If you use port forwarding on your router, you can use a different port on the server as long as the public port on the router is 443.

### Install

1. Download the installation file suitable for your operating system from the [loxwebhook release page](https://github.com/axxelG/loxwebhook/releases/latest)

1. If you cannot use the .deb-packages unpack the binary file. Most probably you want to start those binaries as a service / daemon or background task.

1. Use the provided `config.example.toml` and `controls.d/example.toml.disabled` to create your config and controls file(s).

See [loxwebhook documentation site](https://axxelg.github.io/loxwebhook/) for more details.

## Targets

- [x] Config from config files, environment variables, and flags
- [x] HTTPS encryption with LetsEncrypt
- [x] AuthKey authorization
- [x] Support virtual inputs
- [x] Rate limiting
- [x] Run as Systemd service
- [x] Provide binaries for Raspberry
- [x] Provide binaries for x86 Linux
- [x] Provide binaries for Windows
- [x] Provide .deb packages for Raspian
- [ ] Run as Windows service
- [ ] Support other inputs
- [ ] Encrypt Miniserver Communication and get rid of basic auth even for internal traffic