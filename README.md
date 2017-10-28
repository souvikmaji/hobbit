# hobbit

`hobbit` is a wiki software, which uses a git repository to store it's contents. Hobbit is a a single executable which is supereaasy to install.


# Development

## Install build prerequisite

Install [go-bindata](https://github.com/jteeuwen/go-bindata)

```
$go get github.com/jteeuwen/go-bindata
```

## Build from source

```
$ go get github.com/souvikmaji/hobbit
$ cd $GOPATH/src/github.com/souvikmaji/hobbit
$ go generate
$ go build
$ go test
```
