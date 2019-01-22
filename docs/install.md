# Installation

There are multiple ways to install loxwebhook.

| Install type        | Supported OS                          | Pros                       | Cons |
|---------------------|---------------------------------------|----------------------------|------|
| Debian repository   | Debian                                | Easy setup / easy updates  | - |
| Debian package      | Debian                                | Easy setup                 | update = install |
| Archive             | Linux (amd64 / armhf) / Windows (x64) | flexible setup | Requires some OS knowledge |
| Compile from source | Nearly all | Works on all platforms go supports | Require a go environment and knowledge + setup like from an archive |

## Debian repository

### Install

```bash
# Download public repository key
wget https://axxelg.github.io/loxwebhook/files/deb_repo_pub.key
# Add key to apt key ring (You will trust all repositories signed by this key)
sudo apt-key add deb_repo_pub.key
# Remove download file
rm deb_repo_pub.key
# Install loxwebhook
sudo apt-get update
sudo apt-get install loxwebhook
```

### Update

```bash
sudo apt-get update
sudo apt-get upgrade
```