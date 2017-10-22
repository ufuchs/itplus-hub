//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//

// [ Quality means doing it right when no one is looking. ]
// [                                       - Henry Ford - ]

// @see: https://talks.golang.org/2013/advconc/dedupermain/dedupermain.go
// @see: http://marcio.io/2015/07/singleton-pattern-in-go/

package qualify

import (
	"fmt"
	"time"

	"ufuchs/itplus/base/fcc"
)

type (
	//
	Dispatcher struct {
		closing       chan chan error
		hostname      string
		Out           chan []byte
		ri            *fcc.RunInfo
		seen          map[int]bool
		subscriber    map[int]*Subscriber
		In            <-chan []byte
		runInfoUpdate chan *fcc.RunInfo
	}
)

//
//
//
func NewDispatcher(hostname string) *Dispatcher {

	return &Dispatcher{
		closing:    make(chan chan error),
		hostname:   hostname,
		Out:        make(chan []byte, 256),
		seen:       make(map[int]bool),
		subscriber: make(map[int]*Subscriber),
	}
}

//
//
//
func (d *Dispatcher) Close() error {
	errc := make(chan error)
	d.closing <- errc
	return <-errc
}

//
//
//
func (d *Dispatcher) Run() {

	//
	//
	//
	var registerInfo = func(addr int, alias string) {
		fmt.Printf("==> Dispatcher : %v - register device: %0.2d - %v\n", d.hostname, addr, alias)
	}

	//
	//
	//
	closeSubs := func() {
		fmt.Printf("==> Dispatcher : %v - unregister devices\n", d.hostname)
		switch {
		case len(d.subscriber) == 0:
			fmt.Println("      no device has been registered")
			break
		default:
			for _, subs := range d.subscriber {
				subs.Close()
			}
		}
		close(d.Out)
	}

	////////////////////////////////////////////////////////////////////////////

	for {

		select {
		case errc := <-d.closing:
			closeSubs()
			errc <- nil
			return
		case d.ri = <-d.runInfoUpdate:
			d.runInfoUpdate = nil
			fmt.Printf("==> Dispatcher : %v - set new run info\n", d.hostname)
		case in, ok := <-d.In:

			var num int

			if !ok || len(in) == 0 {
				continue
			}

			addr := getDeviceAddr(in)

			if num, ok = d.ri.DeviceAddr[addr]; !ok {
				// 1. this device is not one of ours
				// 2. or it is not registered
				if d.ri.ShowUnknownDevices {
					fmt.Printf("==> Dispatcher : %v - unknown device: %0.2d --> %v\n", d.hostname, addr, string(in))
				}
				continue
			}

			ts := fcc.Timestamped{
				When: time.Now().Unix(),
				Addr: addr,
				Num:  num,
				Data: string(in),
			}

			if d.seen[addr] {
				d.subscriber[addr].In <- ts
				continue
			}

			ev, OK := d.ri.Events[num]
			if !OK {
				ev = nil
			}

			device := d.ri.Devices[num]

			registerInfo(addr, device.Alias)

			d.subscriber[addr] = NewSubscriber(d, &SubsRunInfo{
				device: device,
				event:  ev,
			})

			d.seen[addr] = true

			d.subscriber[addr].In <- ts

		}

	}

}

//
//
//
func (d *Dispatcher) SetRunInfo(ri *fcc.RunInfo) {
	d.runInfoUpdate = make(chan *fcc.RunInfo, 1)
	d.runInfoUpdate <- ri
}

//
// Gets the device address as decimal
//
func getDeviceAddr(buf []byte) int {

	var val int
	var c byte
	var i int

	// OK CC 24 130 4 106 125
	// OK 9 24 130 4 106 125
	//    ^-- we have to skip this
	for i = 3; buf[i] != 32; i++ {
	}

	// OK 9 24 130 4 106 125
	//      ^-- and this we are need
	val = 0
	for i = i + 1; buf[i] != 32; i++ {

		// input is decimal only
		val = val * 10

		c = buf[i]
		switch {
		case c >= 'A' && c <= 'E':
			c = c - 55
			break
		case c >= '0' && c <= '9':
			c = c - 48
			break

		}

		val = val + int(c)

	}

	return val

}
