package main

import (
	"bytes"
	"fmt"
	"log"
	"net"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func main() {
	// Contoh: hasil benchmark dari file atau stdin
	data := []byte(`PUT_BENCHMARK_RESULTS_HERE`)
	dec := vegeta.NewDecoder(bytes.NewReader(data))

	var m vegeta.Metrics
	for {
		var r vegeta.Result
		if err := dec.Decode(&r); err != nil {
			break
		}
		m.Add(&r)
	}

	// Hubungkan ke StatsD
	conn, err := net.Dial("udp", "127.0.0.1:8125")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Contoh kirim success rate
	fmt.Fprintf(conn, "greenfunction.success_rate:%.2f|g\n", m.Success*100)
	fmt.Fprintf(conn, "greenfunction.throughput:%.2f|g\n", m.Throughput)
	fmt.Fprintf(conn, "greenfunction.latency_mean:%.2f|g\n", m.Latencies.Mean.Seconds()*1000) // ms

	fmt.Println("Metrics sent to StatsD")
}
