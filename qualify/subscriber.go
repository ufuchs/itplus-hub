//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//
// [ The true sign of intelligence is not knowledge but imagination. ]
// [                                             - Albert Einstein - ]
//

package qualify

import (
	"encoding/json"
	"fmt"
	"sync"

	"ufuchs/itplus/base/fcc"
)

type (
	//
	SubsRunInfo struct {
		device fcc.Device
		event  *fcc.Event
	}

	Subscriber struct {
		hostname      string
		out           chan []byte
		closing       chan chan error
		ri            *SubsRunInfo
		In            chan fcc.Timestamped
		runInfoUpdate chan *SubsRunInfo
		wg            sync.WaitGroup
	}
)

//
//
//
func NewSubscriber(parent *Dispatcher, ri *SubsRunInfo) *Subscriber {

	s := &Subscriber{
		closing:  make(chan chan error),
		hostname: parent.hostname,
		In:       make(chan fcc.Timestamped, 1),
		out:      parent.Out,
		ri:       ri,
	}

	go s.run()

	return s

}

// func validateAgainstEvent(m *Measurement, event *fcc.Event) {
// 	if event != nil {
// 	}

// 	return
// }

//
//
//
func (s *Subscriber) SetRunInfo(ri *SubsRunInfo) {
	s.runInfoUpdate = make(chan *SubsRunInfo, 1)
	s.runInfoUpdate <- ri
}

//
//
//
func (s *Subscriber) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

//
//
//
func (s *Subscriber) run() {

	var (
		extractor = &Extractor{}
	)

	// https://talks.golang.org/2013/advconc.slide#26

	for {
		select {
		case errc := <-s.closing:
			fmt.Printf("    %v\n", s.ri.device.Alias)
			errc <- nil
			return
		case s.ri = <-s.runInfoUpdate:
			s.runInfoUpdate = nil
			break
		case ts := <-s.In:

			m, err := extractor.ExtractData(ts, s.ri.device)
			if err != nil {
				fmt.Println(err)
			}
			m.Host = s.hostname

			var data []byte
			if data, err = json.Marshal(&m); err != nil {
				fmt.Println(err)
			}

			s.out <- data
		}
	}
}
