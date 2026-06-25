package main

import (
    "fmt"
    "github.com/SherClockHolmes/webpush-go"
)

func main() {
    privateKey, publicKey, _ := webpush.GenerateVAPIDKeys()
    fmt.Println("Private Key:", privateKey)
    fmt.Println("Public Key:", publicKey)
}