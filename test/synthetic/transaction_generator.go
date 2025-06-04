package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"
)

type Transaction struct {
    Amount    float64 `json:"amount"`
    Currency  string  `json:"currency"`
    Type      string  `json:"type"`
    Recipient string  `json:"recipient"`
}

func main() {
    client := &http.Client{Timeout: 10 * time.Second}
    
    log.Println("Starting synthetic transaction generator...")
    
    for {
        generateTransaction(client)
        time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
    }
}

func generateTransaction(client *http.Client) {
    transaction := Transaction{
        Amount:    rand.Float64() * 1000,
        Currency:  []string{"USD", "ETH", "BTC", "USDT"}[rand.Intn(4)],
        Type:      []string{"payment", "transfer", "exchange"}[rand.Intn(3)],
        Recipient: fmt.Sprintf("synthetic_user_%d", rand.Intn(100)),
    }

    jsonData, _ := json.Marshal(transaction)
    
    resp, err := client.Post("http://localhost:8080/health", 
        "application/json", bytes.NewBuffer(jsonData))
    
    if err != nil {
        log.Printf("Error: %v", err)
        return
    }
    defer resp.Body.Close()
    
    log.Printf("Generated transaction: %+v (Status: %d)", transaction, resp.StatusCode)
}