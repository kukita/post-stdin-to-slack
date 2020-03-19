//
// Name
//   post-stdin-to-slack.go
//
// Description
//   This program posts Std-In to the specified Slack channel.
//
// Copyright (C) 2020 Keisuke Kukita
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
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

// Config is struct
type Config struct {
	SlackIncomingWebhooksURL string `json:"slack_incoming_webhooks_url"`
	SlackBotName             string `json:"slack_bot_name"`
	SlackBotIcon             string `json:"slack_bot_icon"`
	SlackChannel             string `json:"slack_channel"`
	LogEnabled               bool   `json:"log_enabled"`
	LogLevel                 string `json:"log_level"`
}

// SlackPost is struct
type SlackPost struct {
	SlackBotName         string                `json:"username"`
	SlackBotIcon         string                `json:"icon_emoji"`
	SlackChannel         string                `json:"channel"`
	SlackMarkdownEnabled bool                  `json:"mrkdwn"`
	SlackPostAttachments []SlackPostAttachment `json:"attachments"`
}

// SlackPostAttachment is struct
type SlackPostAttachment struct {
	SlackPostAttachmentColor            string                     `json:"color"`
	SlackPostAttachmentEnabledMarkdowns []string                   `json:"mrkdwn_in"`
	SlackPostAttachmentFields           []SlackPostAttachmentField `json:"fields"`
}

// SlackPostAttachmentField is struct
type SlackPostAttachmentField struct {
	SlackPostTitle string `json:"title"`
	SlackPostValue string `json:"value"`
}

func main() {
	// Getting basename.
	basename := filepath.Base(os.Args[0][:len(os.Args[0])-len(filepath.Ext(os.Args[0]))])

	// Checking if there is a configuration file (JSON file) and create a new one if it does not exist.
	configFilePath := os.Args[0][:len(os.Args[0])-len(filepath.Ext(os.Args[0]))] + ".json"
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		config := Config{
			"YOUR_SLACK_INCOMING_WEBHOOKS_URL",
			"ChatOps Bot",
			":loudspeaker:",
			"#general",
			false,
			"INFO",
		}

		configJSON, err := json.Marshal(&config)
		if err != nil {
			log.Fatal(err)
		}

		configFile, err := os.Create(configFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer configFile.Close()

		configFile.WriteString(string(configJSON))
		log.Print("[INFO] ------------------------------------------------------------")
		log.Print("[INFO] Creating '" + configFilePath + "' has finished successfully. At First, please edit this file.")
		log.Print("[INFO] ------------------------------------------------------------")
		return
	}

	// Loading the configuration file (JSON file).
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	json.Unmarshal(configFile, &config)

	// Setting "hashicorp/logutils".
	if config.LogEnabled {
		logFilePath := os.Args[0][:len(os.Args[0])-len(filepath.Ext(os.Args[0]))] + ".log"
		logfile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer logfile.Close()

		filter := &logutils.LevelFilter{
			Levels:   []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
			MinLevel: logutils.LogLevel(config.LogLevel),
			Writer:   io.MultiWriter(logfile, os.Stdout),
		}
		log.SetOutput(filter)
	} else {
		filter := &logutils.LevelFilter{
			Levels:   []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
			MinLevel: logutils.LogLevel(config.LogLevel),
			Writer:   os.Stdout,
		}
		log.SetOutput(filter)
	}

	log.Print("[INFO] ------------------------------------------------------------")
	log.Print("[INFO] '" + basename + "' is starting.")
	log.Print("[INFO] ------------------------------------------------------------")
	postMessage := flag.String("message", "", "Post Message.")
	postType := flag.String("type", "Info", "Post Type. { Info | Success | Warning | Error }")
	flag.Parse()
	log.Print("[DEBUG] Setting post message: " + *postMessage)
	log.Print("[DEBUG] Setting post type: " + *postType)

	postColorMap := map[string]string{
		"Info":    "#1971FF",
		"Success": "#00B06B",
		"Warning": "#F6AA00",
		"Error":   "#FF4B00",
	}
	log.Print("[DEBUG] Setting post color: " + postColorMap[*postType])

	log.Print("[INFO] The following has been entered as standard input.")
	buffers := make([]byte, 0, 1024)
	stdInScanner := bufio.NewScanner(os.Stdin)
	for stdInScanner.Scan() {
		buffers = append(buffers, stdInScanner.Text()...)
		buffers = append(buffers, "\n"...)
	}
	log.Print("[INFO] " + string(buffers))

	slackPostAttachmentField := SlackPostAttachmentField{
		"[" + *postType + "] " + *postMessage,
		"```" + string(buffers) + "```",
	}

	slackPostAttachment := SlackPostAttachment{
		postColorMap[*postType],
		[]string{"fields"},
		[]SlackPostAttachmentField{slackPostAttachmentField},
	}

	slackPost := SlackPost{
		config.SlackBotName,
		config.SlackBotIcon,
		config.SlackChannel,
		true,
		[]SlackPostAttachment{slackPostAttachment},
	}

	log.Print("[DEBUG] Generating JSON string is starting.")
	slackPostJSON, err := json.Marshal(slackPost)
	if err != nil {
		log.Print("[ERROR] Generating JSON string is failed.")
		log.Fatal(err)
	}
	log.Print("[DEBUG] Generating JSON string has finished.")
	log.Print("[DEBUG] " + string(slackPostJSON))

	log.Print("[INFO] Post form data to the following URL.")
	log.Print("[INFO] " + string(config.SlackIncomingWebhooksURL))
	httpPostFormResponce, err := http.PostForm(
		config.SlackIncomingWebhooksURL,
		url.Values{"payload": {string(slackPostJSON)}},
	)
	if err != nil {
		log.Print("[ERROR] Posting form data is failed.")
		log.Fatal(err)
	}

	log.Print("[INFO] HTTP Responce Status: " + httpPostFormResponce.Status)
	log.Print("[INFO] ------------------------------------------------------------")
	log.Print("[INFO] '" + basename + "' has finished successfully.")
	log.Print("[INFO] ------------------------------------------------------------")
	return
}
