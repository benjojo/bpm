package bpm

import (
	"math"
	"math/rand"
)

const (
	RATE     = 44100
	INTERVAL = 128
)

/*
 * Sample from the metered energy
 *
 * No need to interpolate and it makes a tiny amount of difference; we
 * take a random sample of samples, any errors are averaged out.
 */
func sample(nrg []float32, offset float64) float64 {
	n := math.Floor(offset)
	i := int64(n)

	if n >= 0.0 && n < float64(len(nrg)) {
		return float64(nrg[i])
	}
	return 0.0
}

/*
 * Test an autodifference for the given interval
 */
func autodifference(nrg []float32, interval float64) float64 {
	var diff, total float64
	beats := [...]float64{-32, -16, -8, -4, -2, -1, 1, 2, 4, 8, 16, 32}
	nobeats := [...]float64{-0.5, -0.25, 0.25, 0.5}

	mid := rand.Float64() * float64(len(nrg))
	v := sample(nrg, mid)

	diff, total = 0.0, 0.0 // Just be sure I guess?

	for n := 0; n < (len(beats) / 2); n++ {

		y := sample(nrg, mid+beats[n]*interval)

		w := 1.0 / math.Abs(beats[n])
		diff += w * math.Abs(y-v)
		total += w
	}

	for n := 0; n < (len(nobeats) / 2); n++ {

		y := sample(nrg, mid+nobeats[n]*interval)

		w := math.Abs(nobeats[n])
		diff -= w * math.Abs(y-v)
		total += w
	}

	return diff / total
}

/*
 * Beats-per-minute to a sampling interval in energy space
 */
func bpmToInterval(bpm float64) float64 {
	var beatsPerSecond, samplesPerBeat float64

	beatsPerSecond = bpm / 60
	samplesPerBeat = RATE / beatsPerSecond
	return samplesPerBeat / INTERVAL
}

/*
 * Sampling interval in enery space to beats-per-minute
 */
func intervalToBpm(interval float64) float64 {
	var samplesPerBeat, beatsPerSecond float64

	samplesPerBeat = interval * INTERVAL
	beatsPerSecond = float64(RATE) / samplesPerBeat
	return beatsPerSecond * 60
}

// ScanForBpm Scan a range of BPM values for the one with the minimum autodifference,
// needs to be fed a sampled input, you can use ReadFloatArray to do that.
func ScanForBpm(nrg []float32, slowest, fastest float64, steps, samples int) float64 {
	slowest = bpmToInterval(slowest)
	fastest = bpmToInterval(fastest)
	step := (slowest - fastest) / float64(steps)

	height := math.Inf(0)
	trough := math.NaN()

	for interval := fastest; interval <= slowest; interval += step {
		t := float64(0.0)

		for s := 0; s < samples; s++ {
			t += autodifference(nrg, interval)
		}

		if t < height {
			trough = interval
			height = t
		}
	}

	return intervalToBpm(trough)
}

// ReadFloatArray takes a pcm_s16le file and processes it into
// an array that can be used by ScanForBpm
func ReadFloatArray(samples []float32) []float32 {
	var v, n float64
	nrg := make([]float32, 0)

	for i := 0; i < len(samples); i++ {
		z := math.Abs(float64(samples[i]))
		if z > v {
			v += (z - v) / 8
		} else {
			v -= (v - z) / 512
		}

		n++
		if n == INTERVAL {
			n = 0
			nrg = append(nrg, float32(v))
		}
	}

	return nrg
}
