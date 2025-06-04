package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
)

func main() {
    log.Println("Starting alert testing...")
    
    testHighErrorRate()
    time.Sleep(30 * time.Second)
    
    testHighLatency()
    time.Sleep(30 * time.Second)
    
    log.Println("Alert testing completed")
}

func testHighErrorRate() {
    log.Println("Testing high error rate alert...")
    
    var wg sync.WaitGroup
    
    for i := 0; i < 200; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            http.Get("http://localhost:8080/nonexistent-endpoint")
        }()
        
        if i%20 == 0 {
            time.Sleep(100 * time.Millisecond)
        }
    }
    
    wg.Wait()
    log.Println("High error rate test completed")
}

func testHighLatency() {
    log.Println("Testing high latency alert...")
    
    var wg sync.WaitGroup
    
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            http.Get("http://localhost:8080/health?delay=1000")
        }()
        time.Sleep(100 * time.Millisecond)
    }
    
    wg.Wait()
    log.Println("High latency test completed")
}