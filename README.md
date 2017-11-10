# websocketemail-go

This repository is the official go library for [websocket.email](https://websocket.email).
This code lets you subscribe to email addresses at the websocket.email domain from the from go code.

### Installing and using

To build and install the development command line client into $GOPATH/bin run:

```
go get github.com/websocket-email/websocketemail-go/cmd/wsemail
```


To get a prebuilt version of the cli client download one from the [releases page](https://github.com/websocket-email/websocketemail-go/releases).

To use the go library, follow the example provided in [this file](https://github.com/websocket-email/websocketemail-go/blob/master/cmd/wsemail/main.go)

Library documentation can be found [here](https://godoc.org/github.com/websocket-email/websocketemail-go).

## Running the tests

Get a valid API token from [websocket.email](https://websocket.email), change to the project directoryin your go path and run:

```
export WEBSOCKETEMAIL_TOKEN="$YOUR_TOKEN_HERE"
go test
```

## Versioning

We use [SemVer](http://semver.org/) for versioning.

## License

See [LICENSE.md](LICENSE.md) file for details

