[![Go Report Card](https://goreportcard.com/badge/github.com/khihadysucahyo/go-clean-arch-boilerplate)](https://goreportcard.com/report/github.com/khihadysucahyo/go-clean-arch-boilerplate)
[![Maintainability](https://api.codeclimate.com/v1/badges/25184a167ac260f033e9/maintainability)](https://codeclimate.com/github/khihadysucahyo/go-clean-arch-boilerplate/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/25184a167ac260f033e9/test_coverage)](https://codeclimate.com/github/khihadysucahyo/go-clean-arch-boilerplate/test_coverage)
# go-clean-arch

## Description
This is an example of implementation of Clean Architecture in Go (Golang) projects.

Rule of Clean Architecture by Uncle Bob
 * Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
 * Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
 * Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
 * Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
 * Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has  4 Domain layer :
 * Models Layer
 * Repository Layer
 * Usecase Layer  
 * Delivery Layer

#### The diagram:

![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)

The original explanation about this project's structure  can read from this medium's post : https://medium.com/@imantumorang/golang-clean-archithecture-efd6d7c43047.

It may different already, but the concept still the same in application level, also you can see the change log from v1 to current version in Master.

### How To Run This Project
> Make Sure you have run the content.sql in your mysql


Since the project already use Go Module, I recommend to put the source code in any folder but GOPATH.

#### Run the Testing

```bash
$ make test
```

#### Run the Applications
Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/bxcodec/go-clean-arch.git

#move to project
$ cd go-clean-arch

# Build the docker image first
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker ps

# Execute the call
$ curl localhost:6060/contents

# Stop
$ make stop
```


### Tools Used:
In this project, I use some tools listed below. But you can use any simmilar library that have the same purposes. But, well, different library will have different implementation type. Just be creative and use anything that you really need. 

- All libraries listed in [`go.mod`](https://github.com/bxcodec/go-clean-arch/blob/master/go.mod) 
- ["github.com/vektra/mockery".](https://github.com/vektra/mockery) To Generate Mocks for testing needs.
# go-echo-clean-arch
