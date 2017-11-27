# Go-fswatcher

Util to monitor file changes and run commands on event written in Go

## Usage

Just run 

```
$ ./go-fswatcher config.example.json
```

## Configure

Config file cat look like this:

```json
[
  {
    "Path": "/tmp/foo",
    "Commands": ["ping -c 4 ya.ru", "traceroute ya.ru"]
  }
]
```

You can add multiple entires, with path to file and commands list
