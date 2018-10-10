package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/windmilleng/tcell"

	"github.com/windmilleng/tilt/internal/hud/view"
	"github.com/windmilleng/tilt/internal/rty"
)

func main() {
	f, err := os.Create("logfile")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	log.Printf("ahhh")

	d, err := NewDemo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("starting\n")

	err = d.Run()
	if err != nil {
		log.Fatal(err)
	}
}

type Demo struct {
	view   view.View
	model  *model
	screen tcell.Screen
}

type model struct {
	lastKey *tcell.EventKey
}

func NewDemo() (*Demo, error) {
	screen, err := tcell.NewTerminfoScreen()
	if err != nil {
		return nil, err
	}

	r := &Demo{
		screen: screen,
		model:  &model{},
	}

	r.view = view.View{
		Resources: []view.Resource{
			view.Resource{
				"fe",
				"fe",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"be",
				"be",
				[]string{"/"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"graphql",
				"graphql",
				[]string{},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"snacks",
				"snacks",
				[]string{"snacks/whoops.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"doggos",
				"doggos",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"elephants",
				"elephants",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"heffalumps",
				"heffalumps",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"aardvarks",
				"aardvarks",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"quarks",
				"quarks",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"boop",
				"boop",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"laurel",
				"laurel",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"hardy",
				"hardy",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"north",
				"north",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"west",
				"west",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
			view.Resource{
				"east",
				"east",
				[]string{"fe/main.go"},
				time.Second,
				view.ResourceStatusFresh,
				"1/1 pods up",
			},
		},
	}

	return r, nil
}

type rtyNavigationState struct {
	resources rty.ScrollState
	stream    rty.ScrollState
}

func (d *Demo) Run() error {
	d.screen.Init()
	defer d.screen.Fini()
	d.screen.Clear()
	screenEvs := make(chan tcell.Event)
	go func() {
		for {
			screenEvs <- d.screen.PollEvent()
		}
	}()

	// initial render
	if err := d.render(); err != nil {
		return err
	}

	for {
		select {
		case ev := <-screenEvs:
			done := d.handleScreenEvent(ev)
			if done {
				return nil
			}
		}
		if err := d.render(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Demo) handleScreenEvent(ev tcell.Event) bool {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyRune:
			d.model.lastKey = ev
			switch ev.Rune() {
			case 'q':
				return true
			case 'j':
				scroll := d.resourcesScroll()
				selected := scroll.Get()
				if nextResource := d.nextResource(selected); nextResource != "" {
					scroll.Select(rty.ComponentID(nextResource))
				}
			case 'k':
				scroll := d.resourcesScroll()
				selected := scroll.Get()
				if nextResource := d.previousResource(selected); nextResource != "" {
					scroll.Select(rty.ComponentID(nextResource))
				}
			}
		}
	}

	return false
}

func (d *Demo) render() error {
	c := d.TopLevel()
	if err := c.Render(rty.NewScreenCanvas(d.screen)); err != nil {
		return err
	}
	d.screen.Show()
	return nil
}

func (d *Demo) TopLevel() rty.Component {
	l := rty.NewFlexLayout(rty.DirVert)

	l.Add(d.header())
	l.Add(d.resources())
	l.Add(d.footer())

	return l
}

func (d *Demo) header() rty.FixedDimComponent {
	b := rty.NewBox()
	b.SetInner(rty.String("header"))
	return rty.NewFixedDimSize(b, 3)
}

func (d *Demo) resources() rty.Component {
	l := rty.NewScrollLayout(rty.DirVert)

	for _, r := range d.view.Resources {
		rc := d.resource(r)
		l.Add(rc)
	}

	return l
}

func (d *Demo) resource(r view.Resource) rty.FixedDimComponent {
	lines := rty.NewLines()
	cl := rty.NewLine()
	cl.Add(rty.String(r.Name))
	cl.Add(rty.NewFillerString('-'))
	cl.Add(rty.String(fmt.Sprintf("%d", r.Status)))
	lines.Add(cl)
	cl = rty.NewLine()
	cl.Add(rty.String(fmt.Sprintf(
		"LOCAL: (watching %v) - ", r.DirectoryWatched)))
	cl.Add(rty.NewTruncatingStrings(r.LatestFileChanges))
	lines.Add(cl)
	cl = rty.NewLine()
	cl.Add(rty.String(
		fmt.Sprintf("  K8S: %v", r.StatusDesc)))
	lines.Add(cl)
	cl = rty.NewLine()
	return lines
}

func (d *Demo) footer() rty.FixedDimComponent {
	b := rty.NewBox()
	b.SetInner(rty.String("footer"))

	return rty.NewFixedDimSize(b, 3)
}