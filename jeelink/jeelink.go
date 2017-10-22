//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//

// [ Die Kunst ist, einmal mehr aufzustehen, als man umgeworfen wird.          ]
// [                                                       -Winston Churchill- ]

package jeelink

import (
	"errors"
	"fmt"

	"ufuchs/itplus/base/fcc"
	"ufuchs/itplus/base/zvous"
	"ufuchs/itplus/hub/app"
	"ufuchs/itplus/hub/qualify"
)

var (
	//
	ErrMissingJEELINKEnvironment = errors.New("==> Environment 'export JEELINK_PORT=/dev/...' is missing.")
	//
	ErrMissingJEELINK = errors.New("==> JEELINK missing. Did you insert any one?")
)

const (
	JEELINK_PORT = "JEELINK_PORT"
	BAUD         = 57600
	BUFSIZE      = 384
)

//
//
//
type JeeLink struct {
	Id         int
	quit       chan struct{}
	reader     ReadCloseFlusher
	Out        chan []byte
	dispatcher *qualify.Dispatcher
	err        error
}

//
//
//
func Factory(configName string, entry *zvous.ServiceEntry) (*JeeLink, error) {

	var (
		err    error
		ri     *fcc.RunInfo
		jee    *JeeLink
		reader ReadCloseFlusher
	)

	fmt.Printf("==> Jeelink    : %v - initialize service at '%v'\n", entry.GetIdentifier(), entry.ExtractConn())

	if ri, err = app.GetRuninfo(app.ConfigFilesDir, configName, 0); err != nil {
		return nil, err
	}

	if reader, err = NewJeelinkReader(entry); err != nil {
		return nil, err
	}

	dispatcher := qualify.NewDispatcher(reader.GetHostname())
	dispatcher.SetRunInfo(ri)

	jee = NewJeeLink(reader, entry, dispatcher)

	if err := jee.reader.Flush(); err != nil {
		return nil, err
	}

	go jee.Run()
	go jee.dispatcher.Run()

	return jee, nil

}

//
//
//
func NewJeeLink(reader ReadCloseFlusher, c *zvous.ServiceEntry, dispatcher *qualify.Dispatcher) *JeeLink {

	var jee = &JeeLink{
		quit:       make(chan struct{}),
		Id:         c.ID,
		reader:     reader,
		dispatcher: dispatcher,
		Out:        make(chan []byte, BUFSIZE),
	}

	jee.dispatcher.In = jee.Out

	return jee

}

//
//
//
func (j *JeeLink) GetHostname() string {
	return j.reader.GetHostname()
}

//
//
//
func (j *JeeLink) Close() {
	close(j.quit)
	return
}

//
//
//
func (j *JeeLink) Run() {

	type readResult struct {
		n   int
		err error
	}

	var (
		nextDelimIndex = 0
		delimiter      = []byte{13, 10}
		received       = make([]byte, BUFSIZE)
		text           = []byte{}
		readDone       = make(chan readResult, 1)
	)

	// var read = func(received []byte) <-chan readResult {
	// 	n, err := j.reader.Read(received)
	// 	readDone <- readResult{n, err}
	// 	return readDone
	// }

	for {

		go func(received []byte) {
			n, err := j.reader.Read(received)
			readDone <- readResult{n, err}
		}(received)

		select {
		case <-j.quit:
			j.reader.Close()
			j.dispatcher.Close()
			close(j.Out)
			fmt.Printf("==> Jeelink    : %v - unregistered\n", j.GetHostname())
			return
			//		case r := <-read(received):
		case r := <-readDone:

			var c byte
			for i := 0; i < r.n; i++ {

				if c = received[i]; c == 0 {
					continue
				}

				if c == delimiter[nextDelimIndex] {
					nextDelimIndex++
				} else {
					text = append(text, c)
				}

				if nextDelimIndex == 2 {

					nextDelimIndex = 0

					if len(text) == 0 {
						continue
					}

					// if isMeasurement(text) {
					// 	j.Out <- append([]byte(nil), text...)
					// } else {
					// 	log.Println("X", (text[:6]))
					// }

					if isMeasurement(text) {
						j.Out <- append([]byte(nil), text...)
					}

					text = text[:0]
					continue
				}

			}

		}

	}

}

//
//
//
func isMeasurement(buf []byte) bool {
	var ok bool

	// still unhandled...
	// 10 91 76 97 67 114

	// 79 75 32 86 65
	if len(buf) > 5 {
		ok = buf[0] == 79 && buf[1] == 75
		if buf[3] == 86 && buf[4] == 65 {
			return false
		}
	}
	return ok
}
