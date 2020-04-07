package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/codahale/hdrhistogram"
	"io"
	"os"
	"strconv"
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

func main() {
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
}

func fatalOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
