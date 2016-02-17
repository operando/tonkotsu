# TONKOTSU

Google Play Application & App Store Application Update Checker in Go

## Run

```bash
go run tonkotsu.go config.go -c config.toml
```

## Config File

```toml
log = "debug"
sleeptime = 1
error_post = true

[slack_update_post]
text = "Update!!!!!"
username = "TONKOTSU bot"
icon_emoji = ":pig:"
channel = "#test"
link_names = true

[slack_error_post]
text = "Error!!!!!"
username = "bot"
icon_emoji = ":ghost:"
channel = "#test"
link_names = true

[slack_start_post]
text = "Running tonkotsu..."
username = "bot"
icon_emoji = ":ghost:"
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

## License

```
Copyright 2016 Shinobu Okano

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
