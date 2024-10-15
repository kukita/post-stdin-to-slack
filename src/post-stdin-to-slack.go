package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/hashicorp/logutils"
)

type Config struct {
	SlackIncomingWebhooksURL string `json:"slack_incoming_webhooks_url"`
	SlackBotName             string `json:"slack_bot_name"`
	SlackBotIcon             string `json:"slack_bot_icon"`
	SlackChannel             string `json:"slack_channel"`
	LogEnabled               bool   `json:"log_enabled"`
	LogLevel                 string `json:"log_level"`
}

type SlackPost struct {
	Username    string                `json:"username"`
	IconEmoji   string                `json:"icon_emoji"`
	Channel     string                `json:"channel"`
	Mrkdwn      bool                  `json:"mrkdwn"`
	Attachments []SlackPostAttachment `json:"attachments"`
}

type SlackPostAttachment struct {
	Color    string                     `json:"color"`
	MrkdwnIn []string                   `json:"mrkdwn_in"`
	Fields   []SlackPostAttachmentField `json:"fields"`
}

type SlackPostAttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
}

const (
	colorInfo    = "#1971FF"
	colorSuccess = "#00B06B"
	colorWarning = "#F6AA00"
	colorError   = "#FF4B00"
)

func main() {
	config := loadConfig()
	setupLogging(config)

	message, postType := parseFlags()
	inputText := readStdin()

	slackPost := createSlackPost(config, message, postType, inputText)
	postToSlack(config.SlackIncomingWebhooksURL, slackPost)

	log.Print("[INFO] post completed successfully")
}

func loadConfig() Config {
	configPath := getConfigPath()
	createConfigIfNotExists(configPath)

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}
	return config
}

func getConfigPath() string {
	executable, _ := os.Executable()
	return executable[:len(executable)-len(filepath.Ext(executable))] + ".json"
}

func createConfigIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		config := Config{
			SlackIncomingWebhooksURL: "YOUR_SLACK_INCOMING_WEBHOOKS_URL",
			SlackBotName:             "ChatOps Bot",
			SlackBotIcon:             ":loudspeaker:",
			SlackChannel:             "#general",
			LogEnabled:               false,
			LogLevel:                 "INFO",
		}

		configJSON, _ := json.Marshal(&config)
		if err := ioutil.WriteFile(path, configJSON, 0644); err != nil {
			log.Fatalf("failed to create config file: %v", err)
		}
		log.Printf("[INFO] created new config file: %s", path)
	}
}

func setupLogging(config Config) {
	var writer io.Writer
	if config.LogEnabled {
		logPath := getLogPath()
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		writer = io.MultiWriter(logFile, os.Stdout)
	} else {
		writer = os.Stdout
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
		MinLevel: logutils.LogLevel(config.LogLevel),
		Writer:   writer,
	}
	log.SetOutput(filter)
}

func getLogPath() string {
	executable, _ := os.Executable()
	return executable[:len(executable)-len(filepath.Ext(executable))] + ".log"
}

func parseFlags() (string, string) {
	message := flag.String("message", "", "Post Message")
	postType := flag.String("type", "Info", "Post Type { Info | Success | Warning | Error }")
	flag.Parse()
	return *message, *postType
}

func readStdin() string {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for scanner.Scan() {
		input += scanner.Text() + "\n"
	}
	return input
}

func createSlackPost(config Config, message, postType, input string) SlackPost {
	var color string
	switch postType {
	case "Info":
		color = colorInfo
	case "Success":
		color = colorSuccess
	case "Warning":
		color = colorWarning
	case "Error":
		color = colorError
	default:
		color = colorInfo
	}

	attachment := SlackPostAttachment{
		Color:    color,
		MrkdwnIn: []string{"fields"},
		Fields: []SlackPostAttachmentField{
			{
				Title: "[" + postType + "] " + message,
				Value: "```" + input + "```",
			},
		},
	}

	return SlackPost{
		Username:    config.SlackBotName,
		IconEmoji:   config.SlackBotIcon,
		Channel:     config.SlackChannel,
		Mrkdwn:      true,
		Attachments: []SlackPostAttachment{attachment},
	}
}

func postToSlack(webhookURL string, post SlackPost) {
	postJSON, err := json.Marshal(post)
	if err != nil {
		log.Fatalf("failed to marshal slack post: %v", err)
	}

	resp, err := http.PostForm(webhookURL, url.Values{"payload": {string(postJSON)}})
	if err != nil {
		log.Fatalf("failed to post to slack: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("[INFO] slack api response: %s", resp.Status)
}
