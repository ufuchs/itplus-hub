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

// //
// //
// //
// func newReadResult(color string) *readResult {
// 	return &readResult{
// 		color: []byte(color),
// 		buf:   make([]byte, BUFSIZE, BUFSIZE*2),
// 	}
// }

// //
// //
// //
// func run(conn net.Conn) {
// 	var (
// 		//text     = make([]byte, BUFSIZE, BUFSIZE*2)
// 		//firstRun = true
// 		wg     sync.WaitGroup
// 		textOK = false
// 	)

// 	var (
// 		INa      = newReadResult("green")
// 		INb      = newReadResult("red")
// 		INin     *readResult
// 		INout    *readResult
// 		INtoggle = false
// 	)

// 	var (
// 		i             uint64
// 		j             int
// 		OUTa          = make([]byte, ROUNDS*100, ROUNDS*100*2)
// 		OUTb          = make([]byte, ROUNDS*100, ROUNDS*100*2)
// 		OUTin, OUTout *[]byte
// 		OUTtoggle     = false
// 	)
// 	OUTb = OUTb[:0]
// 	OUTin = &OUTb
// 	OUTout = &OUTa

// 	//text = text[:0]

// 	for {

// 		if INtoggle {
// 			INin = INb
// 			INout = INa
// 		} else {
// 			INin = INa
// 			INout = INb
// 		}
// 		INtoggle = !INtoggle

// 		wg.Add(2)

// 		go func(r *readResult, wg *sync.WaitGroup) {
// 			r.n, r.err = conn.Read(r.buf)
// 			//			fmt.Println(r.n, string(r.buf))
// 			wg.Done()
// 		}(INin, &wg)

// 		go func(r *readResult, text *[]byte, wg *sync.WaitGroup) {

// 			// if textOK {
// 			// 	text = text[:0]
// 			// 	textOK = false
// 			// }

// 			for i := 0; i < r.n; i++ {

// 				var c = r.buf[i]

// 				switch c {
// 				case 0:
// 					continue
// 				case 10:
// 					*text = append(*text, c)
// 					if i < r.n {

// 						continue
// 					}
// 					textOK = true
// 					fmt.Println("...")
// 					goto end
// 				default:
// 					*text = append(*text, c)
// 				}

// 			}

// 		end:

// 			wg.Done()

// 		}(INout, OUTin, &wg)

// 		wg.Wait()

// 		//fmt.Println(string(text))

// 		if i == ROUNDS {

// 			if OUTtoggle {
// 				OUTb = OUTb[:0]
// 				OUTin = &OUTb
// 				OUTout = &OUTa
// 			} else {
// 				OUTa = OUTa[:0]
// 				OUTin = &OUTa
// 				OUTout = &OUTb
// 			}
// 			OUTtoggle = !OUTtoggle

// 			go func(b *[]byte, j int) {
// 				s := fmt.Sprintf("aaa-%02d.txt", j)
// 				ioutil.WriteFile(s, *b, 0644)
// 			}(OUTout, j)

// 			j++
// 			i = 0
// 		} else {
// 			i++
// 		}

// 	}

// }

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
