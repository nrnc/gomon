# Gomon

This is the source code for the second project in the Udemy course Working with Websockets in Go (Golang) by Trevor Sawler. Created while i was taking the course


A Simple Monitoring service

## Build

Build in the normal way on Mac/Linux:

~~~
go build -o gomon cmd/web/*.go
~~~

Or on Windows:

~~~
go build -o gomon.exe cmd/web/.
~~~

Or for a particular platform:

~~~
env GOOS=linux GOARCH=amd64 go build -o vigilate cmd/web/*.go
~~~

## Requirements

Gomon requires:
- Postgres 11 or later (db is set up as a repository, so other databases are possible)
- An account with [Pusher](https://pusher.com/), or a Pusher alternative 
(like [ipê](https://github.com/dimiro1/ipe))

## Run

First, make sure ipê is running (if you're using ipê):

On Mac/Linux
~~~
cd ipe
./ipe 
~~~

On Windows
~~~
cd ipe
ipe.exe
~~~

Run with flags:

~~~
./gomon \
-dbuser='yourusername' \
-pusherHost='localhost' \
-pusherPort='4001' \
-pusherKey='123abc' \
-pusherSecret='abc123' \
-pusherApp="1" \
-pusherSecure=false
~~~~

