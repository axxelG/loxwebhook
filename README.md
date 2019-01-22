# Loxwebhook

Make selected Loxone virtual inputs securely available to any service that can send http requests.

![Show usage in an animated gif](/readme_images/loxwebhook_example.gif)

## Key features

- HTTPS encryption (LetsEncrypt)
- AuthKey authorization

## Use case

Many services that you might want to integrate with your Loxone system are able to send HTTP requests. The Loxone Miniserver is able to accept HTTP requests but you cannot use it for most services because you need to use a complex encryption method that no service I know supports (except services provided by Loxone) or you need to use an unsafe connection method. Loxone is working on implementing a more standard authentication method ([JSON Web Tokens](https://jwt.io/)) but even if JWT are implemented, I doubt that this will be supported by many services.

| Connection method    | Usability  |
|----------------------|------------|
| Unencrypted with [Basic Auth](https://en.wikipedia.org/wiki/Basic_access_authentication) | Because Basic Auth does not provide confidentiality it is not secure to use it on the public internet |
| Encrypted with custom encryption | No service I know supports it. |
| [JSON Web Tokens](https://jwt.io/) | Not yet fully implemented / not supported by many services |

Loxwebhook runs on a separate server (a very low level device like a Raspberry Pi is more than sufficient) to offload the http encryption, protect the Miniserver against DOS attacks (rate limit the requests) and adds an authentication layer based on authKeys to authorize request.

## State of development

Loxwebhook is stable and works for me since end of 2018 without problems.

I use it to connect [IFTTT](https://ifttt.com/) and Loxone. So I am very certain I will actively maintain Loxwebhook as long as Loxone does not offer a build in IFTTT integration.

If you find something that is not working like expected, don't hesitate to create a github issue.

## Data flow example setup

![Overview data flow](/readme_images/DataFlowExampleSetup.svg)

## Overview example setup

![Overview example setup](/readme_images/OverviewExampleSetup.svg)

## Setup

### Prerequisites

A server with [supported operating systems](https://axxelg.github.io/loxwebhook/supported_os.html) on your local network reachable on port 443 from the public internet with working DNS resolution. If you use port forwarding on your router, you can use a different port on the server as long as the public port on the router is 443.

### Install

If you want to install Loxwebhook on a Raspberry Pi you can use [this step by step guide](/docs/walkthrough_raspberry_pi.md) otherwise the following steps will give you an idea what you need to do.

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
- [ ] Support other controls
- [ ] Encrypt Miniserver Communication and get rid of basic auth even for internal traffic