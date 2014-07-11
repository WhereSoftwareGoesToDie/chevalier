Chevalier
=========

A search engine for Vaultaire data sources.

Requirements
============

Chevalier requires [Go](https://golang.org) 1.2 or greater,
[ZeroMQ](http://zeromq.org) 4.0 or greater, and [Protocol
Buffers](https://code.google.com/p/protobuf/). (It also requires other
things, but `go get` should be able to find them.)

It talks to a [Vaultaire](https://github.com/anchor/vaultaire) cluster
of version 2.0 or greater, so you'll need one of those also.
