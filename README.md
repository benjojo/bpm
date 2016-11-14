# bpm
[![GoDoc](https://godoc.org/github.com/benjojo/bpm?status.svg)](https://godoc.org/github.com/benjojo/bpm)


Library and tool for dealing with beats per second detection in Go

This is a direct port of Mark Hills [bpm-tools](http://www.pogo.org.uk/~mark/bpm-tools/).

For that reason, it is also under the same licence as that utility, ( GPLv2 ).

You can use this version as a libary too, In testing you may find it produces slightly
different results than bpm-tools. This is because both tools use a random float on mid
point selection, Ideally it would use a static one, but I won't replicate drand exactly
to just get 1:1 matching with the offical tools.
