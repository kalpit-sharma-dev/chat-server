package main

import (
    "fmt"
    "runtime"
    "time"
)

func printStats() {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)

    fmt.Println("-------- Runtime Stats --------")
    fmt.Printf("CPUs: %d\n", runtime.NumCPU())
    fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    fmt.Println("-------- Memory Stats --------")
    fmt.Printf("Alloc: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
    fmt.Printf("Total Alloc: %.2f MB\n", float64(memStats.TotalAlloc)/1024/1024)
    fmt.Printf("Sys: %.2f MB\n", float64(memStats.Sys)/1024/1024)
    fmt.Printf("Heap Alloc: %.2f MB\n", float64(memStats.HeapAlloc)/1024/1024)
    fmt.Printf("Heap Sys: %.2f MB\n", float64(memStats.HeapSys)/1024/1024)
    fmt.Printf("Num GC: %v\n", memStats.NumGC)
    fmt.Println("------------------------------")
}

func main() {
    // Simulate load (optional)
    go func() {
        for {
            _ = make([]byte, 1024*1024) // Allocate 1 MB per iteration
            time.Sleep(100 * time.Millisecond)
        }
    }()

    // Print stats every 5 seconds
    for {
        printStats()
        time.Sleep(5 * time.Second)
    }
}
