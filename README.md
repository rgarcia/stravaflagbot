# stravaflagbot

checks a strava activity and flags it if it is not flagged

# usage

```
go run main.go -h
Usage: main --sessionid SESSIONID --activityid ACTIVITYID --flagcomment FLAGCOMMENT

Options:
  --sessionid SESSIONID [env: STRAVA_SESSION]
  --activityid ACTIVITYID
  --flagcomment FLAGCOMMENT
  --help, -h             display this help and exit
```

To get your session ID, log in to Strava in Chrome. Open Developer Tools, go to Application -> Cookies and look for the `_strava4_session` cookie. Copy the value and assign it to a STRAVA_SESSION environment variable:

```
export STRAVA_SESSION=...
```

Clone this repo and run the program (assumes you have Go installed):

```
go run main.go --activityid <activity id> --flagcomment "This KOM is not legitimate"
```
