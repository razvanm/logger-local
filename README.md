# logger-local

This is small demo client/server for Syncbase. It assumes a proper Vanadium
environment ([instructions](https://vanadium.github.io/installation/))
and the `v23:syncbase` profile installed (`jiri profile install v23:syncbase`).

## Setup

The setup includes two `syncbased` servers (one for the server and one for the
client) that publish their address to a local `mountabled` server. The server
program connects to their dedicated `syncbased` server. A setup program creates
the databases, the collections and the _syncgroup_ used to synchronize the two.

## How to run demo

Build the binaries, the certs:

    make

Start the servers (one `mounttabled` and two `syncbased`):

    make start-servers

Create the hierarchies on the two syncbases:

    make setup

Run the server:

    make server

In a separate terminal run the client (the data will show up on the
terminal running the server):

    make client

Stop the servers:

    make kill-servers

Clean everything:

    make clean