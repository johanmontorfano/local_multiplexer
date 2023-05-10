package net_test

import (
	"johanmnto/epr/net"
	"testing"
)

func TestConvert(t *testing.T) {
	outgoingES, err := net.MakeOutgoingEventStreamFromURL("https://sse.dev/test")
	if err != nil {
		t.FailNow()
	}

	if err := outgoingES.ReadEvents(func(readLine string) { println(readLine) }); err != nil {
		t.FailNow()
	}
}
