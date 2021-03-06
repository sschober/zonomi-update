# zonomi-update

A simple, fast and small command-line tool to update your 
[zonomi](https://zonomi.com) domain.

## Overview

zonomi-update updates your zonomi domains to your current external
IP address. I logs its operation to syslog by default.

## Usage

Use it via a crontab entry like so:

    */10 * * * * $HOME/bin/zonomi-update

## Configuration

It can be configured via `$HOME/.zonomi-update`. Use

    zonomi-update -dumpflags > ~/.zonomi-update

to get a skeleton config file.

**It is absolutely neccessary to set your `domain` and `apikey`!**

## Installation

To install it, first install the [`iniflags`](github.com/vharitonsky/iniflags) 
package

    go get -u -a github.com/vharitonsky/iniflags

then build and install the `zonomi-update` itself via

    go get -a -u github.com/sschober/zonomi-update
    
# Author

Sven Schober <sv3sch@gmail.com>