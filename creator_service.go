package main

import "fmt"

type Delivery struct{ AssetID, SubscriberID, Content string }
type DeliveryResult struct {
	AssetID            string
	SubscriberNotified bool
}

type processor func(string) (string, error)
type notifier func(string, string) error

func deliver(d Delivery, process processor, notify notifier, capture func(map[string]any) error) (DeliveryResult, error) {
	content, err := process(d.Content)
	if err != nil {
		_ = capture(map[string]any{"title": "content processing failed", "message": err.Error(), "level": "error", "fingerprint": []string{"content-processing", d.AssetID}, "exception": fmt.Sprintf("%T: %v", err, err), "context": map[string]any{"asset_id": d.AssetID, "subscriber_id": d.SubscriberID}})
		return DeliveryResult{}, err
	}
	if err := notify(d.SubscriberID, content); err != nil {
		return DeliveryResult{}, err
	}
	return DeliveryResult{AssetID: d.AssetID, SubscriberNotified: true}, nil
}
