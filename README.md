# Godaddy dns updater
This is a beginner golang script to learn how to read/write files, make api calls
and play with concurrency.

example:

```
export ip_file_path="/tmp/ip.json"
export ip_finder_url="https://api.ipify.org?format=json"
export godaddy_api_url="example.com"
export godaddy_api_key="thisismysecretapikey"
export godaddy_api_secret="thisismysecretapisecret"
export domain_url="example.com"


go run godaddy_dns_updater.go
```

