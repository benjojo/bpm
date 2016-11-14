package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/benjojo/bpm"
)

func main() {
	min := flag.Float64("min", 120, "min BPM you are expecting")
	max := flag.Float64("max", 200, "min BPM you are expecting")
	flag.Parse()

	if flag.Arg(0) == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(flag.Arg(0))

	if err != nil {
		log.Fatalf("Unable to open file, Err: %s", err.Error())
	}

	floats := make([]float32, 0)

	for {
		var f float32
		err = binary.Read(file, binary.LittleEndian, &f)

		if err != nil {
			break
		}

		floats = append(floats, f)
	}

	nrg := bpm.ReadFloatArray(floats)
	bpmr := bpm.ScanForBpm(nrg, *min, *max, 1024, 1024)

	fmt.Printf("%f", bpmr)
}
