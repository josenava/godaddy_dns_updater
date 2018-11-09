# Godaddy dns updater
This is a beginner golang script to learn how to read/write files, make api calls
and play with concurrency.

## WIP currently only getting the current ip address and writing it to a file is working :)

To use it just export `ip_file_path`, `ip_finder_url`,

`godaddy_api_key`, `godaddy_api_secret`, `godaddy_url`, `domain` and `name` env variables
then `go run godaddy_dns_updater.go`

example:

```
export ip_file_path="/tmp/ip.json"
export ip_finder_url="https://api.ipify.org?format=json"
export godaddy_api_url="example.com"
export godaddy_api_key="thisismysecretapikey"
export godaddy_api_secret="thisismysecretapisecret"
export domain="example.com"


go run godaddy_dns_updater.go
```

