package main

import (
	"fmt"

	"ufuchs/itplus/base/zvous"
	"ufuchs/itplus/hub/jeelink"
)

type Glue struct {
	quit         chan struct{}
	jeeDiscovery *zvous.ZCBrowserService
	jeeGroup     *jeelink.Group
}

//
func NewGlue(jeeDiscovery *zvous.ZCBrowserService, jeeGroup *jeelink.Group) *Glue {
	return &Glue{
		quit:         make(chan struct{}),
		jeeDiscovery: jeeDiscovery,
		jeeGroup:     jeeGroup,
	}
}

//
// Close
//
func (g *Glue) Close() {
	close(g.quit)
	return
}

//
func (g *Glue) Run() {

	for {
		select {
		case <-g.quit:
			fmt.Println("==> Glue   : start finalizing")
			g.jeeDiscovery.Close()
			g.jeeGroup.Close()
			fmt.Println("==> Glue   : finalized")
			return
		case conns := <-g.jeeDiscovery.Out:

			var (
				ok  bool
				err error
				jee *jeelink.JeeLink
			)

			for _, conn := range conns {

				if conn.Discharge {
					g.jeeGroup.Remove(conn.ID)
					fmt.Printf("    JeeGroup : %v - removed\n", conn.GetIdentifier())
					continue
				}

				if jee, ok = g.jeeGroup.Contains(conn.ID); ok {
					continue
				}

				if jee, err = jeelink.Factory("home", conn); err != nil {
					fmt.Println(err)
					continue
				}

				g.jeeGroup.Add(jee)
				fmt.Printf("    JeeGroup   : %v - added\n", conn.GetIdentifier())

			}

		}
	}
}
