# tel[l me ]net

`tellmenet` is a small TCP server that writes everything it knows about a
connection when it’s opened, then closes it.

A live version runs at `whoami.bfontaine.net`, e.g.:

    $ telnet whoami.bfontaine.net
    ...
    Network: tcp
    IP: ...
    Global unicast: true
    Multicast: false
    Interface-local multicast: false
    Link-local multicast: false
    Link-local unicast: false
    Loopback: false
    Port: 59683

## Usage

The default is to listen for connection on the port 23 (telnet). Note you need
to be `sudo` in order to run it in a port <1000.
