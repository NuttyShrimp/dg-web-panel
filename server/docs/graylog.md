# Docs for graylog API

Their docs are fucking garbage so here are the small bits we should need for the panel.

To must usefull endpoint to get messages is the: 

All `GET` request send data via query parameters. This is probably going to be the same case for POST request but not sure at this time


GET: `/search/universal/relative`

To request all the messages with all their fields we use the following query. All these options a required for the request to actually respond
```json
{
    "query": "*",
    "fields": "*",
    "range": 0
}
```

We can extend our request with following options:
```json
{
    "query": "logtype:join",
    "limit": 10,
    "offset": 10
    "filter": "streams:6291f3b466fb921361cd8b72",
    "fields": "message,timestamp",
    "sort": "timestamp:desc",
    "range": 300
}
```

| key |required| Description |
|-----|--------|-------------|
| query | x | Their query you would put in the web panel, can be set to `"*"` to show all messages |
| fields | x | Which field of the message should be displayed, can be set to `"*"` to show all fields|
| range | x | the time from now in seconds, if set to `0` it will fetch all the messages from the the time of the first message send to graylog |
| limit |  | the maximum amount of the messages that should be returnend (Does not working when testing on low values) |
| filter |  | filter based on streams and maybe other things, probably only gonna use it to fetch from a specific filter |
| sort |  | sort the messages based on the field can be `:asc` or `:desc`
| offset |  | Amount of messages to skip from start of set, (handy when displaying in eg. a paginated table) |
