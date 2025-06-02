package main

import (
	fltk "github.com/pwiecz/go-fltk"
)

func main() {
	fltk.InitStyles()

	win := fltk.NewWindow(500, 420)
	win.SetLabel("KeyValueGrid Demo")
	win.Resizable(win)

	// Create the KeyValueGrid (replace NewKeyValueGrid with correct import if needed)
	grid := NewKeyValueGrid(win, 20, 20, 460, 350)

	// Add some groups and key-value pairs after layout is set up
	grid.Add("User Account Information", "Name", "Alice")
	grid.Add("User", "Email", "alice@example.com")
	grid.Add("Settings", "Theme", "Dark")
	grid.Add("Settings", "Language", "en-US")
	grid.Add("Network", "Proxy", "")
	grid.Add("Network", "Timeout", "30s")

	win.End()
	win.Show()
	fltk.Run()
}
