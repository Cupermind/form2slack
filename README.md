# form2slack
`Static site -> form2slack-> slack`

You have a static site and you need to post a form.
This go project receives `HTTP POST` request and translates it to slack.

Config file `config.yml`

```
---
slack:
  enable: yes
  token: "YOUR-TOKEN"
  channel: "#CHANNEL"
  title: "You've got new message"
  from: "Your awesome site"
  color: "#36a64f"
form:
  regexp: '^site\-'
  callback_url_field: "callback"
  replace: yes
endpoint: "/ENDPOINT"
port: SOMEPORT
rpm: 1
```


For example if you have the following config:
```
endpoint: "/slack"
port: 8080
```

In form you should specify action `http://hostname:8080/slack`

Fields:

* `slack` - slack settings
* `form` - form settings
*   `regexp` - post only fields matching it.
*   `callback_url_field` put this to your form `<input * type="hidden" name="{{ callback_url_field }}" value="URL"`, so form2slack will redirect to `URL` after posting.
*   `replace` - if set to true regexp will be replaced with empty string
* `endpoint` - endpoint to hit from from, i.e. `/slack`
* `port` - port to bind
* `rpm` - requests per minute from one IP
