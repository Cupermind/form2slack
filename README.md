# form2slack
`Static site -> form2slack-> slack`

You have a static site and you need to post a form.
This go project receives `HTTP POST` request and translates it to slack.

Config file `config.yml`

```
---
slack:
  token: "YOUR-TOKEN"
  channel: "#CHANNEL"
  text: "you've got new message"
endpoint: "/ENDPOINT"
regexp: ".*"
port: SOMEPORT
```

For example if you have the following config:
```
---
slack:
  token: "YOUR-TOKEN"
  channel: "#CHANNEL"
  text: "you've got new message"
endpoint: "/slack"
regexp: ".*"
port: 8080
```

In form you should specify action `http://hostname:8080/slack`

`regexp` - post only fields matching it.
