chmod +x /usr/local/bin/loxwebhook/loxwebhook
systemctl daemon-reload
chown loxwebhook /etc/loxwebhook/config.toml
chmod 0660 /etc/loxwebhook/config.toml
