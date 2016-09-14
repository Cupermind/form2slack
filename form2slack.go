package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/nlopes/slack"
	"gopkg.in/yaml.v2"
)

const configFile string = "config.yml"

//Config of the app
type Config struct {
	Slack struct {
		Enable  bool   `yaml:"enable"`
		Token   string `yaml:"token"`
		Channel string `yaml:"channel"`
		Title   string `yaml:"title"`
		From    string `yaml:"from"`
		Color   string `yaml:"color"`
	}
	Form struct {
		Regexp           string `yaml:"regexp"`
		CallBackURLField string `yaml:"callback_url_field"`
		Replace          bool   `yaml:"replace"`
	}
	EndPoint string `yaml:"endpoint"`
	Port     int    `yaml:"port"`
	RPM      int64  `yaml:"rpm"`
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
	http.Handle(fmt.Sprintf("%s", config.EndPoint),
		tollbooth.LimitFuncHandler(tollbooth.NewLimiter(config.RPM, time.Minute), Slack))

	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//Index function
func Index(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s to Index endpoint", r.RemoteAddr)
	fmt.Fprintf(w, "This is Form2Slack")
}

//Slack function
func Slack(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s to Slack endpoint", r.RemoteAddr)
	callbackURL := ""
	field := ""
	api := slack.New(config.Slack.Token)
	re, err := regexp.Compile(config.Form.Regexp)
	if err != nil {
		log.Fatal("There is an error in regexp: ", err)
		fmt.Fprintf(w, "You have error in regexp config")
		return
	}
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Title:      config.Slack.Title,
		Color:      config.Slack.Color,
		Fields:     []slack.AttachmentField{},
		MarkdownIn: []string{"fields"},
	}
	r.ParseForm()
	for k, v := range r.Form {
		field = k
		if k == config.Form.CallBackURLField {
			callbackURL = strings.Join(v, "")
		} else if re.MatchString(k) && strings.Join(v, "") != "" {

			if config.Form.Replace {
				field = re.ReplaceAllString(k, "")
			}
			attachment.Fields = append(attachment.Fields,
				slack.AttachmentField{
					Title: field,
					Value: fmt.Sprintf("`%s`", strings.Join(v, "")),
				})
		}
	}
	if config.Slack.Enable && len(attachment.Fields) > 0 {
		params.Attachments = []slack.Attachment{attachment}
		channelID, timestamp, err := api.PostMessage(config.Slack.Channel, config.Slack.From, params)
		if err != nil {
			log.Fatal("Error sending to Slack: ", err)
			return
		}
		log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	} else if config.Slack.Enable {
		log.Printf("Warning, there were no fields in the form mathing request")
	} else {
		log.Printf("Sending to Slack disabled")
	}

	if callbackURL != "" {
		http.Redirect(w, r, callbackURL, 301)
		log.Printf("Redirecting user to %s", callbackURL)
	} else {
		log.Printf("Warning, there was no field in the form matching callback_url_field (%s)",
			config.Form.CallBackURLField)
	}
}
