package main

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/immesys/bw2bind"
	"github.com/immesys/spawnpoint/spawnable"
)

var total int64

func main() {
	cl := bw2bind.ConnectOrExit("")
	cl.SetEntityFromEnvironOrExit()
	params := spawnable.GetParamsOrExit()
	from := params.MustString("fromuri")
	to := params.MustString("touri")
	if !strings.HasSuffix(from, "/") {
		from += "/"
	}
	if !strings.HasSuffix(to, "/") {
		to += "/"
	}
	mq := make(chan *bw2bind.SimpleMessage, 1000)
	go sub(cl, from, mq)
	go pub(cl, from, to, mq)
	for {
		time.Sleep(5 * time.Second)
		lt := atomic.LoadInt64(&total)
		fmt.Printf("Forwarded %d messages\n", lt)
	}
}

func sub(cl *bw2bind.BW2Client, from string, mq chan *bw2bind.SimpleMessage) {
	sc := cl.SubscribeOrExit(&bw2bind.SubscribeParams{
		AutoChain: true,
		URI:       from + "*",
	})
	for m := range sc {
		mq <- m
	}
	panic("Subscribe channel ended\n")
}

func pub(cl *bw2bind.BW2Client, from string, to string, mq chan *bw2bind.SimpleMessage) {
	for m := range mq {
		suffix := strings.TrimPrefix(m.URI, from)
		desturi := to + suffix
		cl.PublishOrExit(&bw2bind.PublishParams{
			PayloadObjects: m.POs,
			URI:            to + suffix,
			AutoChain:      true,
		})
		atomic.AddInt64(&total, 1)
	}
	panic("Publish channel ended\n")
}
