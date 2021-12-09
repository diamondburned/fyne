package wayland

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/driver/common"
	"fyne.io/fyne/v2/internal/painter"

	fynedriver "fyne.io/fyne/v2/internal/driver"
)

type driver struct {
	state *state
	appID string
}

// New creates a new Wayland driver for Fyne.
func New(id string) fyne.Driver {
	return &driver{
		state: newWaylandState(),
		appID: id,
	}
}

func (d *driver) CreateWindow(title string) fyne.Window {
	return newWindow(d.state, d.appID, title)
}

func (d *driver) AllWindows() []fyne.Window {
	d.state.mu.Lock()
	defer d.state.mu.Unlock() // try not do this too often

	windows := make([]fyne.Window, 0, len(d.state.windows))
	for w := range d.state.windows {
		windows = append(windows, w)
	}

	return windows
}

func (d *driver) RenderedTextSize(text string, fontSize float32, style fyne.TextStyle) (fyne.Size, float32) {
	return painter.RenderedTextSize(text, fontSize, style)
}

func (d *driver) CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	return common.CanvasForObject(obj)
}

func (d *driver) AbsolutePositionForObject(obj fyne.CanvasObject) fyne.Position {
	c := common.CanvasForObject(obj)
	if c == nil {
		return fyne.Position{}
	}

	wlcanvas := c.(*canvas)
	return fynedriver.AbsolutePositionForObject(obj, wlcanvas.ObjectTrees())
}

func (d *driver) Device() fyne.Device {
	return newDevice(d)
}

func (d *driver) Run() {
	d.state.run()
}

func (d *driver) Quit() {
	d.state.quit()
}

func (d *driver) StartAnimation(*fyne.Animation) {}

func (d *driver) StopAnimation(*fyne.Animation) {}
