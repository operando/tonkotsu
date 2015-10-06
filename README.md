# TONKOTSU

Google Play Application & App Store Application Update Checker in Go

## Run

```bash
go run googleplay_update_checker.go config.go -c config.toml
```

## Config File

```toml
log = "debug"
sleeptime = 1

[slack]
text = "TONKOTSU TEST"
username = "TONKOTSU bot"
icon_emoji = ":pig:"
channel = "#test"
link_names = true

[webhook]
url = "webhook_url" # your Incoming WebHooks URL for Slack

[ios]
app_id = "id667861049" # your iOS application app id
country = "jp"

[android]
package = "com.mercariapp.mercari" # your Android application package name

```
