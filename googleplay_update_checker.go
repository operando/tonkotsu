package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/operando/golack"
)

const (
	GOOGLE_PLAY = "https://play.google.com/store/apps/details?id="
	APP_STORE   = "https://itunes.apple.com/{{country}}/app/{{appId}}"
)

var old_update_date string
var new_update_date string

var oldSoftwareVersion string
var newSoftwareVersion string

func checkUpdate(url string) bool {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		return false
	}
	isUpdate := false
	doc.Find("div[itemprop=\"datePublished\"]").Each(func(_ int, s *goquery.Selection) {
		log.Debug(s.Text())
		if old_update_date == "" {
			old_update_date = s.Text()
			log.Info("Old update date : " + old_update_date)
		} else {
			new_update_date = s.Text()
			if old_update_date != new_update_date {
				log.Info("New update date : " + new_update_date)
				isUpdate = true
			}
		}
	})
	log.Debug(isUpdate)
	return isUpdate
}

func checkUpdateIos(url string) bool {
	isUpdate := false
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		return false
	}
	doc.Find("ul.list span[itemprop=softwareVersion]").Each(func(_ int, s *goquery.Selection) {
		log.Debug(s.Text())
		if oldSoftwareVersion == "" {
			oldSoftwareVersion = s.Text()
			log.Info("Old Software Version : " + oldSoftwareVersion)
		} else {
			newSoftwareVersion = s.Text()
			if oldSoftwareVersion != newSoftwareVersion {
				log.Info("New Software Version : " + newSoftwareVersion)
				isUpdate = true
			}
		}
	})
	log.Debug(isUpdate)
	return isUpdate
}

func createAppStoreURL(ios Ios) string {
	replaceCountryURL := strings.Replace(APP_STORE, "{{country}}", ios.Country, 1)
	appStoreURL := strings.Replace(replaceCountryURL, "{{appId}}", ios.AppId, 1)
	log.Debug(appStoreURL)
	return appStoreURL
}

func createGooglePlayURL(android Android) string {
	googlePlayURL := GOOGLE_PLAY + android.Package
	log.Debug(googlePlayURL)
	return googlePlayURL
}

func main() {
	var configPath = flag.String("c", "", "configuration file path")
	flag.Parse()

	var config Config
	_, err := LoadConfig(*configPath, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	setLogLevel(config.Log)
	sleep := time.Duration(config.SleepTime*60) * time.Second

	var googlePlayURL string
	var appStoreURL string
	payload := golack.Payload{
		config.Slack,
	}

	checkIos := true
	if config.Ios.AppId == "" {
		checkIos = false
		log.Debug("AppId is empty.")
	} else {
		appStoreURL = createAppStoreURL(config.Ios)
		log.Info("Check App Store URL : " + appStoreURL)
	}
	checkAndroid := true
	if config.Android.Package == "" {
		checkAndroid = false
		log.Debug("Package is empty.")
	} else {
		googlePlayURL = createGooglePlayURL(config.Android)
		log.Info("Check Google Play URL : " + googlePlayURL)
	}

	log.Info("Slack Post Message : " + config.Slack.Text)

	for {
		if checkAndroid {
			if checkUpdate(googlePlayURL) {
				golack.Post(payload, config.Webhook)
				log.Info("Update!!!!!!!!!!!")
				break
			} else {
				log.Info("No Update")
			}
		}
		if checkIos {
			if checkUpdateIos(appStoreURL) {
				golack.Post(payload, config.Webhook)
				log.Info("Update!!!!!!!!!!!")
				break
			} else {
				log.Info("No Update")
			}
		}
		time.Sleep(sleep)
	}

	log.Info("Update check end.")
}

func init() {
	log.SetLevel(log.InfoLevel)
}

func setLogLevel(lv string) {
	switch lv {
	case "debug", "d":
		log.SetLevel(log.DebugLevel)
	case "info", "i":
		log.SetLevel(log.InfoLevel)
	case "warn", "w":
		log.SetLevel(log.WarnLevel)
	case "error", "e":
		log.SetLevel(log.ErrorLevel)
	case "fatal", "f":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
