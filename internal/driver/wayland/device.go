package wayland

import (
	"fyne.io/fyne/v2"
)

type device driver

func newDevice(d *driver) fyne.Device {
	return (*device)(d)
}

func (d *device) Orientation() fyne.DeviceOrientation {
	// TODO orientation
	return fyne.OrientationHorizontalLeft
}

func (d *device) IsMobile() bool {
	// TODO ismobile, probably guess from viewport
	return false
}

func (d *device) HasKeyboard() bool {
	// TODO: check inputs
	return true
}

func (d *device) SystemScaleForWindow(w fyne.Window) float32 {
	// TODO: get scaling
	return -1
}
