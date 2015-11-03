Tracking
========

This is an API used for all tracking events at Reevoo.

Setup
-----

To begin hacking on this project you need to install a few dependencies first.

Go
---
On the mac install go with `brew install go`

it is important to setup GOPATH and PATH corectly, add this to your .profile

```bash
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$PATH
```

Get the code:
```bash
go get github.com/reevoo/tracker
```

to check everything is working properly run `go test` to run the unit tests.

To start the server, run `go run tracker/server.go`

<!---

Docker
------
Building the production binary and running the acceptance suite requires docker.

The quickest way to get up and running is [boot2docker](https://github.com/boot2docker/boot2docker) a lightweight linux vm with docker ready installed that runs on VirtualBox.

Install with `brew install boot2docker`
Then `boot2docker init` to download the vm image
Then `boot2docker up` to start the vm

In order to get the command line docker client to connect to the docker deamon running inside the boot2docker VM add this to your .profile
```bash
export DOCKER_HOST=tcp://$(boot2docker ip 2>/dev/null):2375
```

-->