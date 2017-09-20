package jeelink

import (
	"fmt"
	"io"
	"net"

	"github.com/tarm/serial"
	"hidrive.com/ufuchs/itplus/base/fcc"
	"hidrive.com/ufuchs/itplus/base/zvous"
)

type Flusher interface {
	Flush() (err error)
}

// ReadCloser is the interface that groups the basic Read and Close methods.
type ReadCloseFlusher interface {
	io.Reader
	io.Closer
	Flusher
	GetHostname() string
}

type USBReader struct {
	reader   *serial.Port
	hostname string
}

// https://www.hugopicado.com/2016/09/26/simple-data-processing-pipeline-with-golang.html
type TCPReader struct {
	reader   net.Conn
	hostname string
}

//
//
//
func NewJeelinkReader(c *zvous.ServiceEntry) (ReadCloseFlusher, error) {

	var (
		err           error
		jeelinkReader ReadCloseFlusher
	)

	switch c.Itype {
	case zvous.NET:
		jeelinkReader, err = NewTCPReader(c.ExtractConn(), c.ExtractHostname())
	case zvous.USB:
		jeelinkReader, err = NewUSBReader(c.ExtractHostname())
	}

	return jeelinkReader, err

}

//
// NewUSBReader
//
func NewUSBReader(hostname string) (*USBReader, error) {

	c, err := NewSerialConfig()
	if err != nil {
		return nil, err
	}

	if err = IsDeviceConnected(c); err != nil {
		return nil, err
	}

	var port *serial.Port
	if port, err = serial.OpenPort(c); err != nil {
		fcc.Fatalf("==> JEELINK throws '%v'", err)
	}

	s := fmt.Sprintf("==> JEELINK is connected with '%v'", c.Name)
	fmt.Println(s)

	return &USBReader{
		hostname: hostname,
		reader:   port,
	}, err

}

//
//
//
func (r *USBReader) Flush() error {
	return r.reader.Flush()
}

//
//
//
func (r *USBReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

//
//
//
func (r *USBReader) Close() error {
	return r.reader.Close()
}

//
//
//
func (r *USBReader) GetHostname() string {
	return r.hostname
}

//
// NewUSBReader
//
func NewTCPReader(ip string, hostname string) (*TCPReader, error) {

	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return nil, err
	}

	fmt.Printf("    Jeelink    : %v - connection established\n", conn.RemoteAddr())

	return &TCPReader{
		hostname: hostname,
		reader:   conn,
	}, nil

}

//
//
//
func (r *TCPReader) Flush() error {
	return nil
}

//
//
//
func (r *TCPReader) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

//
//
//
func (r *TCPReader) Close() error {
	return r.reader.Close()
}

//
//
//
func (r *TCPReader) GetHostname() string {
	return r.hostname
}
