https://community.letsencrypt.org/t/renewal-after-manual-support-of-dns-01-in-automated-plugins/26100

```bash
certbot certonly --manual --manual-auth-hook /path/to/your/script --preferred-challenges dns
```

https://certbot.eff.org/docs/using.html?highlight=auth%20hook#pre-and-post-validation-hooks

```bash
certbot certonly --manual-public-ip-logging-ok --agree-tos -n -m sre@example.com --preferred-challenges=dns --test-cert  --manual --manual-auth-hook /opt/certbot/validation.sh --manual-cleanup-hook /opt/certbot/clean.sh -d t1.example.com -d t2.example.com
```

# wildcard:
--server https://acme-v02.api.letsencrypt.org/directory
