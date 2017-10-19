package sse

import (
	"time"
)

type ClientOption interface {
	Apply(*SSEClient) error
}

func (wf WantFlag) Apply(ssec *SSEClient) error {
	ssec.demand(wf)
	return nil
}

type ReconnectTime time.Duration

func (rt ReconnectTime) Apply(ssec *SSEClient) error {
	ssec.reconnectTime = time.Duration(rt)
	return nil
}

type ServerOption interface {
	Apply(*EventSink) error
}

type KeepAliveTime time.Duration

func (kat KeepAliveTime) Apply(sink *EventSink) error {
	sink.keepAliveTime = time.Duration(kat)
	return nil
}
