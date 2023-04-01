package game

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type DemoMenu struct {
	active bool
	ui     *ebitenui.UI
	window *widget.Window
	res    *uiResources
	game   *Game

	resolutions []MenuResolution

	// held vars that should not get updated in real-time
	newMinLightRGB [3]float32
	newMaxLightRGB [3]float32
}

type MenuResolution struct {
	width, height int
	aspectRatio   MenuAspectRatio
}

type MenuAspectRatio struct {
	w, h, fov int
}

func (r MenuResolution) String() string {
	if r.aspectRatio.w == 0 || r.aspectRatio.h == 0 {
		return fmt.Sprintf("(*) %dx%d", r.width, r.height)
	}
	return fmt.Sprintf("(%d:%d) %dx%d", r.aspectRatio.w, r.aspectRatio.h, r.width, r.height)
}

func createMenu(g *Game) *DemoMenu {
	res, err := NewUIResources()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// using empty background container since settings will be in a window
	bg := widget.NewContainer()
	var ui *ebitenui.UI = &ebitenui.UI{
		Container: bg,
	}

	menu := &DemoMenu{
		game:        g,
		ui:          ui,
		res:         res,
		active:      false,
		resolutions: g.generateMenuResolutions(),
	}

	menu.initMenu()

	return menu
}

func (m *DemoMenu) initMenu() {
	root := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			// It is using a GridLayout with a single column
			widget.GridLayoutOpts.Columns(1),
			// It uses the Stretch parameter to define how the rows will be layed out.
			// - a fixed sized header
			// - a content row that stretches to fill all remaining space
			// - a fixed sized footer
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
			// Padding defines how much space to put around the outside of the grid.
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    20,
				Bottom: 20,
			}),
			// Spacing defines how much space to put between each column and row
			widget.GridLayoutOpts.Spacing(0, 20))),
		widget.ContainerOpts.BackgroundImage(m.res.background),
	)

	// window title
	titleBar := titleBarContainer(m)

	// settings pages
	settings := settingsContainer(m)
	root.AddChild(settings)

	// footer
	footer := footerContainer(m)
	root.AddChild(footer)

	ww, wh := ebiten.WindowSize()
	window := widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(root),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.MinSize(ww/2, wh/2),
		widget.WindowOpts.MaxSize(ww, wh),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Resize: ", args.Rect)
		}),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Move: ", args.Rect)
		}),
	)

	r := image.Rect(0, 0, 2*ww/3, 2*wh/3)
	r = r.Add(image.Point{ww / 8, wh / 16})
	window.SetLocation(r)
	if m.window != nil {
		m.window.Close()
	}
	m.ui.AddWindow(window)
	m.window = window
}

func (g *Game) openMenu() {
	g.paused = true
	g.menu.active = true
	g.mouseMode = MouseModeCursor
	ebiten.SetCursorMode(ebiten.CursorModeVisible)

	// for color menu items [1, 1, 1] represents NRGBA{255, 255, 255}
	g.menu.newMinLightRGB = [3]float32{
		float32(g.minLightRGB.R) * 1 / 255, float32(g.minLightRGB.G) * 1 / 255, float32(g.minLightRGB.B) * 1 / 255,
	}
	g.menu.newMaxLightRGB = [3]float32{
		float32(g.maxLightRGB.R) * 1 / 255, float32(g.maxLightRGB.G) * 1 / 255, float32(g.maxLightRGB.B) * 1 / 255,
	}

	g.menu.initMenu()
}

func (g *Game) generateMenuResolutions() []MenuResolution {
	resolutions := make([]MenuResolution, 0)

	ratios := []MenuAspectRatio{
		{5, 4, 64},
		{4, 3, 68},
		{3, 2, 74},
		{16, 9, 84},
		{21, 9, 100},
	}

	widths := []int{
		640,
		800,
		960,
		1024,
		1280,
		1440,
		1600,
		1920,
	}

	for _, r := range ratios {
		for _, w := range widths {
			h := (w / r.w) * r.h
			resolutions = append(
				resolutions,
				MenuResolution{width: w, height: h, aspectRatio: r},
			)
		}
	}

	return resolutions
}

func (g *Game) closeMenu() {
	g.paused = false
	g.menu.active = false
}

func (m *DemoMenu) layout(w, h int) {
	// TODO: react to game window layout size/scale changes
	//m.mgr.SetDisplaySize(float32(w), float32(h))
}

func (m *DemoMenu) update() {
	if !m.active {
		return
	}

	m.ui.Update()
}

func (m *DemoMenu) draw(screen *ebiten.Image) {
	if !m.active {
		return
	}

	m.ui.Draw(screen)
}
