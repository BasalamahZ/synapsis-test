## Documentation API

This is the documentation API for this project [Postman](https://documenter.getpostman.com/view/19464042/2sA2r81PZJ)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

#### Golang

You need to have [Go v1.19.4](https://golang.org/dl/) installed on your machine. Follow the [official installation guide](https://golang.org/doc/install) to install Go. Or, follow [managing installations guide](https://go.dev/doc/manage-install) to have multiple Go versions on your machine.

#### PostgreSQL

This service has dependency with PostgreSQL. For development environment, you need to have a PostgreSQL server running on your machine.

### Building

1. Once you have all the prerequisites, you can start by cloning this repository into your machine.

```sh
$ mkdir -p $GOPATH/src/github.com/synapsis-test/
$ cd $GOPATH/src/github.com/synapsis-test
$ git clone https://github.com/BasalamahZ/synapsis-test.git
$ cd synapsis-test
```

> The rest of this instructions assumes that your current working directory is on `$GOPATH/src/github.com/synapsis-test/`

2. Build binaries using the `go build` command.

```sh
$ go build ./cmd/synapsistest-api-http
```

### Running

1. If needed, you can modify the app config for development environment through .env file.

2. Execute the binary to start the service

```sh
$ ./synapsistest-api-http
```

## Directory Structure

This repository is organized with the following structure

```
synapsis-test
|-- cmd                                 # Contains executables codes
|   |-- synapsistest-api-http           # HTTP server
|-- global                              # Contains helper files
|   |-- helper                 
|-- internal                            # Application service packages
```

## Contributing

### Code

Application service packages should be developed in the `internal` directory, as those logic should not be used/imported by external repositories.

Application service packages are made using the domain-driven design concept. Some articles to read:

- https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- https://medium.easyread.co/golang-clean-archithecture-efd6d7c43047

Application service package's naming should be self-explanatory about its purpose, so that other developers would not misinterpret the package.
