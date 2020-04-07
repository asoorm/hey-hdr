package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/codahale/hdrhistogram"
	"io"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type hey struct {
	ResponseTime  int64
	DNSPlusDialup time.Duration
	DNS           time.Duration
	RequestWrite  time.Duration
	ResponseDelay float64
	ResponseRead  float64
	StatusCode    int
	Offset        float64
}

var outFile string

func main() {
	flag.StringVar(&outFile, "out", "", "file to write hdr e.g. `hdr.csv`")
	flag.Parse()

	var h hey

	hist := hdrhistogram.New(0, 6E10, 4)

	switch flag.NArg() {
	case 0:
		r := csv.NewReader(os.Stdin)
		// Read and throw away header
		_, _ = r.Read()
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			fatalOnErr(err)

			responseTime, _ := strconv.ParseFloat(record[0], 64)
			h.ResponseTime = int64(1E6 * responseTime)
			hist.RecordValue(h.ResponseTime)
			//fmt.Printf("RAW: %s RECORDING: %d\n", record[0], h.ResponseTime)
		}

		//bars := hist.Distribution()
		//for k, v := range bars {
		//	if v.Count > 0 {
		//		fmt.Printf("bar: %d, %#v\n", k, v)
		//	}
		//}

		fmt.Fprintf(os.Stdout, "  Count: %d\n", hist.TotalCount())
		fmt.Fprintf(os.Stdout, "    Max: %s\n", time.Duration(hist.Max())*time.Microsecond)
		fmt.Fprintf(os.Stdout, "   Mean: %s\n", time.Duration(hist.Mean())*time.Microsecond)
		fmt.Fprintf(os.Stdout, "    P50: %s\n", time.Duration(hist.ValueAtQuantile(50))*time.Microsecond)
		fmt.Fprintf(os.Stdout, "    P95: %s\n", time.Duration(hist.ValueAtQuantile(95))*time.Microsecond)
		fmt.Fprintf(os.Stdout, "    P99: %s\n", time.Duration(hist.ValueAtQuantile(99))*time.Microsecond)
		fmt.Fprintf(os.Stdout, "   P999: %s\n", time.Duration(hist.ValueAtQuantile(99.9))*time.Microsecond)
		fmt.Fprintf(os.Stdout, "  P9999: %s\n", time.Duration(hist.ValueAtQuantile(99.99))*time.Microsecond)
		fmt.Fprintf(os.Stdout, " P99999: %s\n", time.Duration(hist.ValueAtQuantile(99.999))*time.Microsecond)

		break
	default:
		fmt.Fprint(os.Stderr, "input must be from stdin\n")
		os.Exit(1)
	}

	if outFile == "" {
		os.Exit(0)
	}

	file, err := os.Create(outFile)
	fatalOnErr(err)
	defer file.Close()

	tw := tabwriter.NewWriter(file, 0, 8, 2, ' ', tabwriter.StripEscape)
	_, err = fmt.Fprintf(tw, "Value(ms)\tPercentile\tTotalCount\t1/(1-Percentile)\n")
	fatalOnErr(err)

	total := float64(hist.TotalCount())
	for _, q := range logarithmic {
		value := (time.Duration(hist.ValueAtQuantile(q * 100)) * time.Microsecond).Seconds()*1000
		oneBy := oneByQuantile(q)

		count := int64((q * total) + 0.5) // Count at quantile
		_, err = fmt.Fprintf(tw, "%.3f\t%f\t%d\t%f\n", value, q, count, oneBy)
		fatalOnErr(err)
	}

	fatalOnErr(tw.Flush())
}

func oneByQuantile(q float64) float64 {
	if q < 1.0 {
		return 1 / (1 - q)
	}
	return float64(10000000)
}

var logarithmic = []float64{
	0.00,
	0.100,
	0.200,
	0.300,
	0.400,
	0.500,
	0.550,
	0.600,
	0.650,
	0.700,
	0.750,
	0.775,
	0.800,
	0.825,
	0.850,
	0.875,
	0.8875,
	0.900,
	0.9125,
	0.925,
	0.9375,
	0.94375,
	0.950,
	0.95625,
	0.9625,
	0.96875,
	0.971875,
	0.975,
	0.978125,
	0.98125,
	0.984375,
	0.985938,
	0.9875,
	0.989062,
	0.990625,
	0.992188,
	0.992969,
	0.99375,
	0.994531,
	0.995313,
	0.996094,
	0.996484,
	0.996875,
	0.997266,
	0.997656,
	0.998047,
	0.998242,
	0.998437,
	0.998633,
	0.998828,
	0.999023,
	0.999121,
	0.999219,
	0.999316,
	0.999414,
	0.999512,
	0.999561,
	0.999609,
	0.999658,
	0.999707,
	0.999756,
	0.99978,
	0.999805,
	0.999829,
	0.999854,
	0.999878,
	0.99989,
	0.999902,
	0.999915,
	0.999927,
	0.999939,
	0.999945,
	0.999951,
	0.999957,
	0.999963,
	0.999969,
	0.999973,
	0.999976,
	0.999979,
	0.999982,
	0.999985,
	0.999986,
	0.999988,
	0.999989,
	0.999991,
	0.999992,
	0.999993,
	0.999994,
	0.999995,
	0.999996,
	0.999997,
	0.999998,
	0.999999,
	1.0,
}

func fatalOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
