# bpm
[![GoDoc](https://godoc.org/github.com/benjojo/bpm?status.svg)](https://godoc.org/github.com/benjojo/bpm)


Library and tool for dealing with beats per second detection in Go

This is a direct port of Mark Hills [bpm-tools](http://www.pogo.org.uk/~mark/bpm-tools/).

For that reason, it is also under the same licence as that utility, ( GPLv2 ).

You can use this version as a libary too, In testing you may find it produces slightly
different results than bpm-tools. This is because both tools use a random float on mid
point selection, Ideally it would use a static one, but I won't replicate drand exactly
to just get 1:1 matching with the offical tools.

## Usage of the command line util

You need to feed the command line utility PCM 32 bit little edian floats (mono), there are
two easy ways to do this:

Sox:

```
sox "$FILE" -r 44100 -e float -c 1 -t raw - | ./cmd /dev/stdin
```

ffmpeg:

```
ffmpeg -v quiet -i 1479012090.ts -f f32le -ac 1 -c:a pcm_f32le -ar 44100 pipe:1 | ./cmd /dev/stdin
```

You can also ask for a "snapshot" per second of the BPM calculated, using progressive mode:

```
$ ffmpeg -v quiet -i 1479012090.ts -f f32le -ac 1 -c:a pcm_f32le -ar 44100 pipe:1 | ./cmd -progressive=true /dev/stdin
179.020979
193.086109
191.044776
```

The first value is the BPM for the first 10 seconds, the 2nd for 10-20 seconds, 3rd is the 20-30th (and more until you stop
giving it data)


