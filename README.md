# pcall

## Name

*pcall* - resolve A and AAA records by running an arbitrary external command

## Description

The *pcall* plugin provide a quick way to extend CoreDNS without writing a new plugin

## Syntax

~~~ txt
pcall {
    run /path/to/externa/command
}
~~~

## Examples

Start a server on the default port and load the *pcall* plugin.

~~~ corefile
example.org {
    pcall {
        run /path/to/externa/command
    }
}
~~~

*pcall* will run the command with query type and query name as paramaters in the respective order