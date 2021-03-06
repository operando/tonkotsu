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
	VERSION = "1.1.1"
)

const (
	GOOGLE_PLAY = "https://play.google.com/store/apps/details?id="
	APP_STORE   = "https://itunes.apple.com/{{country}}/app/{{appId}}"
	KINDLE_STORE = "https://www.amazon.co.jp/gp/product/"
)

var oldUpdateDate string
var newUpdateDate string

var oldSoftwareVersion string
var newSoftwareVersion string

func checkUpdate(url string) (bool, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Error(err)
		return false, err
	}
	isUpdate := false
	doc.Find("div[itemprop=\"datePublished\"]").Each(func(_ int, s *goquery.Selection) {
		log.Debug(s.Text())
		if oldUpdateDate == "" {
			oldUpdateDate = s.Text()
			log.Info("Old update date : " + oldUpdateDate)
		} else {
			newUpdateDate = s.Text()
			if oldUpdateDate != newUpdateDate {
				log.Info("New update date : " + newUpdateDate)
				isUpdate = true
			}
		}
	})
	log.Debug(isUpdate)
	return isUpdate, nil
}

func checkUpdateIos(url string) (bool, error) {
	isUpdate := false
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Error(err)
		return false, err
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
	return isUpdate, nil
}

func checkUpdateKindle(url string) (bool, error) {
	isUpdate := false
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Error(err)
		return false, err
	}
	doc.Find("#productDetailsTable li:nth-child(3)").Each(func(_ int, s *goquery.Selection) {
		log.Debug(strings.TrimSpace(s.Text()))
		if oldSoftwareVersion == "" {
			oldSoftwareVersion = strings.TrimSpace(s.Text())
			log.Info(oldSoftwareVersion)
		} else {
			newSoftwareVersion = strings.TrimSpace(s.Text())
			if oldSoftwareVersion != newSoftwareVersion {
				log.Info(newSoftwareVersion)
				isUpdate = true
			}
		}
	})
	log.Debug(isUpdate)
	return isUpdate, nil
}

func createAppStoreURL(ios Ios) string {
	replaceCountryURL := strings.Replace(APP_STORE, "{{country}}", ios.Country, 1)
	appStoreURL := strings.Replace(replaceCountryURL, "{{appId}}", ios.AppID, 1)
	log.Debug(appStoreURL)
	return appStoreURL
}

func createGooglePlayURL(android Android) string {
	googlePlayURL := GOOGLE_PLAY + android.Package
	log.Debug(googlePlayURL)
	return googlePlayURL
}

func createKindleStoreURL(kindle Kindle) string {
	kindleStoreURL := KINDLE_STORE + kindle.Asin
	log.Debug(kindleStoreURL)
	return kindleStoreURL
}

func main() {
	var configPath = flag.String("c", "", "configuration file path")
	var version = flag.Bool("v",false,"version")
	flag.Parse()

	if *version {
		fmt.Printf(`tonkotus version %s`,VERSION);
		return
	}

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
	var kindleStoreURL string
	uPayload := golack.Payload{
		config.SlackUpdatePost,
	}

	ePayload := golack.Payload{
		config.SlackErrorPost,
	}

	if config.SlackStartPost.Text != "" {
		sPayload := golack.Payload{
			config.SlackStartPost,
		}
		golack.Post(sPayload, config.Webhook)
	}

	checkIos := true
	if config.Ios.AppID == "" {
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
	checkKindle := true
	if config.Kindle.Asin == "" {
		checkKindle = false
		log.Debug("Asin is empty.")
	} else {
		kindleStoreURL = createKindleStoreURL(config.Kindle)
		log.Info("Check Kindle Store URL : " + kindleStoreURL)
	}
	log.Info("Slack Post Message : " + config.SlackUpdatePost.Text)
	log.Info("Slack Errro Message : " + config.SlackErrorPost.Text)

	for {
		if checkAndroid {
			isUpdate, err := checkUpdate(googlePlayURL)
			if err != nil && config.ErrorPost {
				ePayload.Slack.Text = ePayload.Slack.Text + "\n" + err.Error()
				log.Error("Slack Errro Message : " + ePayload.Slack.Text)
				golack.Post(ePayload, config.Webhook)
			}
			if isUpdate {
				golack.Post(uPayload, config.Webhook)
				log.Info("Update!!!!!!!!!!!")
				break
			} else {
				log.Info("No Update")
			}
		}
		if checkIos {
			isUpdate, err := checkUpdateIos(appStoreURL)
			if err != nil && config.ErrorPost {
				golack.Post(ePayload, config.Webhook)
			}
			if isUpdate {
				golack.Post(uPayload, config.Webhook)
				log.Info("Update!!!!!!!!!!!")
				break
			} else {
				log.Info("No Update")
			}
		}
		if checkKindle {
			isUpdate, err := checkUpdateKindle(kindleStoreURL)
			if err != nil && config.ErrorPost {
				golack.Post(ePayload, config.Webhook)
			}
			if isUpdate {
				golack.Post(uPayload, config.Webhook)
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
