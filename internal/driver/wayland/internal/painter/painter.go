package painter

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/internal"
	"fyne.io/fyne/v2/internal/driver"
)

// Paint is the main entry point for a simple software painter.
// The canvas to be drawn is passed in as a parameter and the return is an
// image containing the result of rendering.
func Paint(dst *image.RGBA, c fyne.Canvas) {
	paint := func(obj fyne.CanvasObject, pos, clipPos fyne.Position, clipSize fyne.Size) bool {
		w := fyne.Min(clipPos.X+clipSize.Width, c.Size().Width)
		h := fyne.Min(clipPos.Y+clipSize.Height, c.Size().Height)
		clip := image.Rect(
			internal.ScaleInt(c, clipPos.X),
			internal.ScaleInt(c, clipPos.Y),
			internal.ScaleInt(c, w),
			internal.ScaleInt(c, h),
		)
		switch o := obj.(type) {
		case *canvas.Image:
			drawImage(c, o, pos, dst, clip)
		case *canvas.Text:
			drawText(c, o, pos, dst, clip)
		case gradient:
			drawGradient(c, o, pos, dst, clip)
		case *canvas.Circle:
			drawCircle(c, o, pos, dst, clip)
		case *canvas.Line:
			drawLine(c, o, pos, dst, clip)
		case *canvas.Raster:
			drawRaster(c, o, pos, dst, clip)
		case *canvas.Rectangle:
			drawRectangle(c, o, pos, dst, clip)
		}

		return false
	}

	driver.WalkVisibleObjectTree(c.Content(), paint, nil)
	for _, o := range c.Overlays().List() {
		driver.WalkVisibleObjectTree(o, paint, nil)
	}
}
