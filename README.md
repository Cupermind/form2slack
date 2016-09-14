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
callback_url_field: "FIELD"
rpm: 1
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
callback_url_field: "FIELD"
rpm: 1
```

In form you should specify action `http://hostname:8080/slack`

Fields:

* `slack` - slack settings
* `endpoint` - endpoint to hit from from, i.e. `/slack`
* `regexp` - post only fields matching it.
* `port` - port to bind
* `callback_url_field` put this to your form `<input * type="hidden" name="{{ callback_url_field }}" value="URL"`, so form2slack will redirect to `URL` after posting.
* `rpm` - requests per minute from one IP
