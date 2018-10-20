package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/benjojo/bpm"
)

var (
	min                 = flag.Float64("min", 120, "min BPM you are expecting")
	max                 = flag.Float64("max", 200, "max BPM you are expecting")
	progressive         = flag.Bool("progressive", false, "Print the BPM for every period")
	progressiveInterval = flag.Int("interval", 10, "How many seconds for every progressive chunk printed")
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(flag.Arg(0))

	if err != nil {
		log.Fatalf("Unable to open file, Err: %s", err.Error())
	}

	in := make(chan float32)
	out := make(chan float32)

	go bpm.ProgressivelyReadFloatArray(in, out)

	// floats := make([]float32, 0)
	done := make(chan bool)

	go readProgressiveVars(out, done, *progressive, *progressiveInterval)
	for {
		var f float32
		err = binary.Read(file, binary.LittleEndian, &f)

		if err != nil {
			break
		}

		in <- f
	}
	close(in)
	file.Close()

	select {
	case <-done:
		os.Exit(0)
	}
}

func calcChunkLen(second int) int {
	return (bpm.RATE / bpm.INTERVAL) * second
}

func readProgressiveVars(input chan float32, done chan bool, progressive bool, pint int) {
	if progressive {
		maxsize := calcChunkLen(pint)

		nrg := make([]float32, 0)
		for f := range input {
			nrg = append(nrg, f)
			if len(nrg) == maxsize {
				fmt.Printf("%f\n", bpm.ScanForBpm(nrg, *min, *max, 1024, 1024))
				nrg = make([]float32, 0)
			}
		}
		done <- true
	} else {
		nrg := make([]float32, 0)
		for f := range input {
			nrg = append(nrg, f)
		}
		fmt.Printf("%f\n", bpm.ScanForBpm(nrg, *min, *max, 1024, 1024))
		done <- true
	}
}
