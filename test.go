package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"gopkg.in/yaml.v2"
)

const configFile string = "config.yml"

//Config of the app
type Config struct {
	Slack struct {
		Token   string `yaml:"token"`
		Channel string `yaml:"channel"`
		Text    string `yaml:"text"`
	}
	Regexp   string `yaml:"regexp"`
	EndPoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
}

//App config
var config Config

//Read config file
func readConfig() {
	var filename string

	if len(os.Args) < 2 {
		filename = configFile
	} else {
		filename = os.Args[1]
	}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
}

func main() {
	readConfig()

	http.HandleFunc("/", Index)
	http.HandleFunc(fmt.Sprintf("%s", config.EndPoint), Slack)

	fmt.Printf("Slack token: %s", config.Slack.Token)

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Index function
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is Form2Slack")
}

//Slack function
func Slack(w http.ResponseWriter, r *http.Request) {
	api := slack.New(config.Slack.Token)
	fmt.Printf("Slack token: %s", config.Slack.Token)
	re, err := regexp.Compile(config.Regexp)
	if err != nil {
		fmt.Printf("Regexp error")
		return
	}
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text:   config.Slack.Text,
		Fields: []slack.AttachmentField{},
	}
	r.ParseForm()
	for k, v := range r.Form {
		if re.MatchString(k) {
			attachment.Fields = append(attachment.Fields,
				slack.AttachmentField{
					Title: k,
					Value: strings.Join(v, ""),
				})
		}
	}
	if len(attachment.Fields) > 0 {
		params.Attachments = []slack.Attachment{attachment}
		channelID, timestamp, err := api.PostMessage(config.Slack.Channel, "Form2Slack", params)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	}
}
