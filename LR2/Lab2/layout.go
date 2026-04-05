package main

const (
	// Panel dimensions
	LeftPanelWidth = 25
	PanelGap       = 4
	PanelMargin    = 2

	// Border & padding constants
	BorderWidth   = 2 // Rounded border takes 2 chars (1 each side)
	PaddingWidth  = 2 // Padding of 1 on each side = 2 total
	OverheadWidth = BorderWidth + PaddingWidth
)

type Layout struct {
	WindowWidth  int
	WindowHeight int

	// Left panel
	LeftPanelHeight int

	// Right panels
	RightPanelWidth   int
	TopPanelHeight    int
	BottomPanelHeight int

	// Internal component dimensions
	TextareaInputWidth   int
	TextareaInputHeight  int
	TextareaOutputWidth  int
	TextareaOutputHeight int
	FilepickerHeight     int
}

func NewLayout(width, height int) Layout {
	l := Layout{
		WindowWidth:  width,
		WindowHeight: height,
	}

	// Calculate panel dimensions
	totalH := height - PanelMargin
	l.RightPanelWidth = width - LeftPanelWidth - 8 // Matches existing calculation: msg.Width - leftWidth - 8

	// Split vertical space for right panels
	contentHeight := totalH - 8 // Matches existing calculation: totalH - 8
	l.TopPanelHeight = contentHeight / 2
	l.BottomPanelHeight = contentHeight/2 + contentHeight%2

	l.LeftPanelHeight = totalH - 4 // Matches existing calculation: totalH - 4

	// Internal dimensions (subtract borders + padding + label line)
	l.TextareaInputWidth = l.RightPanelWidth - 2
	l.TextareaInputHeight = l.TopPanelHeight - 3

	l.TextareaOutputWidth = l.RightPanelWidth - 2
	l.TextareaOutputHeight = l.BottomPanelHeight - 3

	l.FilepickerHeight = l.TopPanelHeight - 3

	return l
}
