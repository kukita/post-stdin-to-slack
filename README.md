# Slack Poster

Slack Posterは、標準入力から受け取ったテキストをSlackチャンネルに投稿するGo言語で書かれたコマンドラインツールです<br>
(Slack Poster is a command-line tool written in Go that posts text received from standard input to a Slack channel.)

## 機能 (Features)

- 標準入力からテキストを読み取り、指定されたSlackチャンネルに投稿します<br>
(I read text from standard input and post it to the specified Slack channel.)
- メッセージタイプ（Info、Success、Warning、Error）に基づいて色分けします<br>
(I color-code messages based on message types (Info, Success, Warning, Error).)
- JSON設定ファイルによる簡単な設定が可能です<br>
(I allow easy configuration using a JSON configuration file.)
- 詳細なログ記録オプションを提供します<br>
(I provide detailed logging options.)

## インストール (Installation)

```bash
go get github.com/yourusername/slack-poster
```

## 設定 (Configuration)

初回実行時に、実行ファイルと同じディレクトリに`slack-poster.json`という設定ファイルを作成します<br>
(I create a configuration file named `slack-poster.json` in the same directory as the executable file on first run.)

このファイルを編集して、SlackのWebhook URLやその他の設定を行ってください<br>
(Please edit this file to configure the Slack Webhook URL and other settings.)

```json
{
  "slack_incoming_webhooks_url": "YOUR_SLACK_INCOMING_WEBHOOKS_URL",
  "slack_bot_name": "ChatOps Bot",
  "slack_bot_icon": ":loudspeaker:",
  "slack_channel": "#general",
  "log_enabled": false,
  "log_level": "INFO"
}
```

## 使用方法 (Usage)

```bash
echo "Hello, Slack!" | slack-poster -message "Greeting" -type Info
```

### オプション (Options)

- `-message`: 投稿するメッセージのタイトルです（必須）<br>
(The title of the message to post (required).)
- `-type`: メッセージのタイプです（Info、Success、Warning、Error）（デフォルト: Info）<br>
(The type of the message (Info, Success, Warning, Error) (default: Info).)

## ログ (Logging)

ログを有効にするには、設定ファイルの`log_enabled`を`true`に設定してください<br>
(To enable logging, set `log_enabled` to `true` in the configuration file.)

ログファイルは実行ファイルと同じディレクトリに`slack-poster.log`という名前で作成します<br>
(I create a log file named `slack-poster.log` in the same directory as the executable file.)

## ライセンス (License)

このプロジェクトはApache License 2.0のもとで公開しています<br>
(I release this project under the Apache License 2.0.)

詳細は[LICENSE](LICENSE)ファイルをご覧ください<br>
(Please see the [LICENSE](LICENSE) file for details.)

## 貢献 (Contributing)

バグ報告や機能リクエストは、GitHubのIssueでお願いします<br>
(Please submit bug reports and feature requests through GitHub Issues.)

Pull Requestも歓迎します<br>
(I also welcome Pull Requests.)

## 作者 (Author)

Keisuke Kukita
