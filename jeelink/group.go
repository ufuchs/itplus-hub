package jeelink

import (
	"context"
	"fmt"
)

type Group struct {
	ctx      context.Context
	JeeLinks []*JeeLink
	Out      chan []byte
}

//
//
//
func NewGroup(ctx context.Context) *Group {
	return &Group{
		ctx:      ctx,
		JeeLinks: []*JeeLink{},
		Out:      make(chan []byte, 256),
	}
}

//
// Add
//
func (g *Group) Add(j *JeeLink) {

	g.JeeLinks = append(g.JeeLinks, j)

	go func() {

		for {
			select {
			case <-j.quit:
				j.reader.Close()
				j.dispatcher.Close()
				// close(j.Out)
				fmt.Println("==> JeeGroup   :", j.GetHostname(), "- unregister fan in")
				return
			case in := <-j.dispatcher.Out:
				g.Out <- in
			}
		}

	}()

}

//
//
//
func (g *Group) Contains(id int) (jee *JeeLink, ok bool) {

	for _, jee = range g.JeeLinks {
		if ok = jee.Id == id; ok {
			return
		}
	}
	return nil, false
}

//
//
//
func (g *Group) Close() {
	for _, jee := range g.JeeLinks {
		jee.Close()
	}
	fmt.Println("\r==> JeeGroup   : all JeeLinks closed")

}

//
// Remove
//
func (g *Group) Remove(id int) {
	var newJeeLinks = g.JeeLinks[:0]
	for _, jee := range g.JeeLinks {
		if jee.Id == id {
			jee.Close()
			jee = nil
			continue
		}
		newJeeLinks = append(newJeeLinks, jee)

	}
	g.JeeLinks = newJeeLinks
}
