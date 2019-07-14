# vaali

## Dependencies
Vaali is a golang based web/cli application framework with following which, with thin wrappers, integrates:
* Command line interface with https://github.com/urfave/cli
* MongoDB - with mgo
* Web app framework with Labstack echo - https://github.com/labstack/echo

## Features
Vaali provides following features using above mentioned libraries:
* User management, configurable authentication and authorization schemes
* Access control on API based on a pre-defined URL scheme
* Decalrative filters for data types stored in MongoDB - helps to populate and apply the filter in a data type independent way
* Declarative tables for data types stored in MongoDB - based on column specs generates data for displaying on a web page
* User friendly commandline - if user misses a required argument, it is prompted before the command is executed
* Module system which helps statically plugin REST endpoints, commands, mongo db indices etc when using Vaali as a library
* EMail based user registration
* JSON based configuration system

## Installing
* Using dep:
```dep ensure add github.com/varunamachi/vaali```
* Using go get:
```go get -v github.com/varunamachi/vaali```
