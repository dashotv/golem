# Golem

A toolchain, similar to Rails or other application frameworks, that instead
generates idiomatic go.

Supports generation of a simple application, with support for:
* Command-line
* Web
* Database Models
* Background Jobs

This functionality is built on the backs of several technologies:
* `Cobra` and `Viper` for command-line and configuration
* `Gin Gonic` for web routing
* `MGM` for Mongodb connections and models
* `JobRunner` for background jobs and management

## Why

I found myself building very similar apps over and over, I was interested
in things like `Kallax` and other tools that help solve repitition in go
applications by using code generation tools.

I have existing Ruby on Rails applications that use Mongodb databases and
wanted something to simplify replacing those RoR apps with go-based apps

I wanted to do this in a way that was idiomatic go, rather than using
runtime reflection and things of that nature.

Utlimately, I'd like to see this support additional `dialects` for generating
code compatible with other simple frameworks. Like the aforementioned `Kallax`
for integrating with `postgresql` compatible databases, or `Gorilla Mux` for
web requests and routing.

## Quickstart

This assumes you have a working go installation

> go get -u github.com/dashotv/golem

> golem new <name> <repo>

* \<name\> is the application name
* \<repo\> is the full package repo name

Example:

> golem new blarg github.com/dashotv/blarg

This will generate an application similar to the following tree:

```
blarg
├── LICENSE
├── application
│   └── app.go
├── cmd
│   ├── root.go
│   └── server.go
├── config
│   └── config.go
├── main.go
├── models
│   ├── connector.go
│   ├── document.go
│   ├── hello.go
│   └── schema.go
└── server
    ├── releases
    │   └── routes.go
    ├── routes.go
    └── server.go
```

This structure is a `cobra` application, with the addition of a few additional
pacakges. See `Anatomy` for more information.
`config` package
that manages configuration of the application and a `server` subcommand that
runs the application for you. The `config` is also configured to automaticall
unmarshal for you, with all the support from `viper` for environment and
command line options overriding the settings.

The `.golem` directory contins the `golem` configuration and additional files
used for generation.

* `.golem/models`: contains yaml model definition files
* `.golem/routes.yaml`: contains route definition files
* `.golem/jobs.yaml`: contains job definition files

## Anatomy

A simple overview of the generated files.

### Application

The application package manages shared configuration and clients with all other
aspects of the system.

### Config

A simple configuration structure setup to use `Viper` unmarshalling. Because the
generated application is compatible with `Cobra` and `Viper` you can use
functionality from those tools to support additional functionality like
environment variable overrides.

There is also a stubbed validation function which is configured to be called
after `Viper` loads the configuration. If you look in `cmd/root.go` you can
find the wiring that makes this possible, as well as the definitions for where
to search for configurations files at runtime.

### Models

The generated model structures are placed here and include a `Connector` class
which manges the connections to the database(s). Use the `models.NewConnector()`
function to get an instance of the `Connector`.

The generated models include a model object, a `Store` structure which handles interaction with
the database, and a `Query` structure that gives you an easy interface to generating
queries to the data.

### Server

`Golem` assumes that you are building an application to serve web requests and the
`server` package is where all of this is wired up.

For each of the `Group` of `Routes` defined in `.golem/routes.yaml` it generates
a package which implements handler functions that automatically convert the
parameters.
