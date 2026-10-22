package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	c, err := newClient()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	result, err := deliver(Delivery{AssetID: "asset-demo", SubscriberID: "subscriber-demo", Content: "edition-1"}, func(string) (string, error) { return "", errors.New("unsupported demo content") }, func(id, content string) error { fmt.Printf("subscriber %s updated with %s\n", id, content); return nil }, func(p map[string]any) error { _, e := c.capture(p); return e })
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("asset %s delivered; subscriber notified=%t\n", result.AssetID, result.SubscriberNotified)
}
