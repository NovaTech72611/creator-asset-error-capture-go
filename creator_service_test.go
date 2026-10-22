package main

import (
	"errors"
	"testing"
)

func TestDeliverCapturesProcessingErrorAndSkipsSubscriber(t *testing.T) {
	cases := []struct {
		name                       string
		processErr                 error
		wantCaptured, wantNotified bool
	}{
		{"processed asset", nil, false, true},
		{"rejected content", errors.New("unsupported format"), true, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			captured, notified := false, false
			result, err := deliver(Delivery{AssetID: "asset-7", SubscriberID: "sub-2", Content: "raw"}, func(string) (string, error) {
				if tc.processErr != nil {
					return "", tc.processErr
				}
				return "ready", nil
			}, func(string, string) error { notified = true; return nil }, func(map[string]any) error { captured = true; return nil })
			if (err != nil) != tc.wantCaptured {
				t.Fatalf("error=%v", err)
			}
			if captured != tc.wantCaptured || notified != tc.wantNotified {
				t.Fatalf("captured=%v notified=%v", captured, notified)
			}
			if tc.wantNotified && (!result.SubscriberNotified || result.AssetID != "asset-7") {
				t.Fatalf("result=%+v", result)
			}
		})
	}
}
