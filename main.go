//
// Copyright(c) 2017 Uli Fuchs <ufuchs@gmx.com>
// MIT Licensed
//
// [ Zeit ist das, was man an der Uhr abliest              ]
// [                                   - Albert Einstein - ]

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ufuchs/zeroconf"

	"hidrive.com/ufuchs/itplus/base/fcc"
	"hidrive.com/ufuchs/itplus/base/zvous"
	"hidrive.com/ufuchs/itplus/hub/app"
	"hidrive.com/ufuchs/itplus/hub/jeelink"
	"hidrive.com/ufuchs/itplus/hub/socket"
)

//
//
//
func prepare() {

	var err error

	if app.BaseDir, err = os.Getwd(); err != nil {
		fcc.Fatal(err)
	}

	svc := app.NewConfigService().
		RetrieveAll()

	if svc.LastErr != nil {
		fcc.Fatal(svc.LastErr)
	}

}

//
//
//
func handleSignals(cancel context.CancelFunc, sigch <-chan os.Signal) {
	for {
		select {
		case <-sigch:
			fmt.Printf("\r")
			cancel()
		}
	}
}

//
// 	https://github.com/kelseyhightower/envconfig
//
func main() {

	var (
		err         error
		sigs        = make(chan os.Signal, 2)
		mainWG      sync.WaitGroup
		hostname    string
		serviceName = zvous.AVAHI_MEASUREMENT
	)

	prepare()

	ctx, cancel := context.WithCancel(context.Background())
	mainCtx := context.WithValue(ctx, 0, &mainWG)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go handleSignals(cancel, sigs)

	if hostname, err = os.Hostname(); err != nil {
		fcc.Fatal(err)
	}

	//	var configName = "aleta"

	server, err := zeroconf.Register(hostname, serviceName, "local.", 8080, []string{}, nil)
	defer server.Shutdown()
	if err != nil {
		fcc.Fatalf("==> Zeroconf : Registering service '%v' failed - %v", serviceName, err)
	}

	jeeGroup := jeelink.NewGroup(mainCtx)
	//	jeeDiscovery := discovery.NewTCPDiscoverService(discovery.SERVICENAME, zeroconf.IPv4, 4)
	jeeDiscovery := zvous.NewZCBrowserService(zvous.AVAHI_DATA, zeroconf.IPv4, 4)

	jeeGlue := NewGlue(jeeDiscovery, jeeGroup)

	go jeeGlue.Run()

	var hub = socket.NewHub()

	hub.In = jeeGroup.Out

	go socket.Run(8080, hub)

	select {
	case <-ctx.Done():
		jeeGlue.Close()

		//jeeGroup.Close()

		//fmt.Println("\r")
		time.Sleep(1000 * time.Millisecond)
		mainWG.Wait()
		fmt.Println("==> App stopped...")
	}

}
