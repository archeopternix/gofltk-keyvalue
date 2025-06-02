package main

import (
	"fmt"

	fltk "github.com/pwiecz/go-fltk"
)

func main() {
	fltk.InitStyles()

	win := fltk.NewWindow(500, 420)
	win.SetLabel("KeyValueGrid Demo")
	win.Resizable(win)

	// Create the KeyValueGrid (replace NewKeyValueGrid with correct import if needed)
	grid := NewKeyValueGrid(win, 20, 20, 460, 350)

	// Button row: horizontal pack for buttons
	buttonRow := fltk.NewPack(0, win.H()-50, 500, 40)
	buttonRow.SetType(fltk.HORIZONTAL)

	// Button to print data from UI edits
	printBtn := fltk.NewButton(20, 0, 120, 30, "Print Data")
	printBtn.SetCallback(func() {
		tree := grid.GetData()
		fmt.Println("==== Data from UI ====")
		for _, group := range tree.Groups {
			fmt.Printf("[%s]\n", group.Name)
			for _, elem := range group.Elements {
				fmt.Printf("%s = %s\n", elem.Key, elem.Value)
			}
		}
	})
	buttonRow.Add(printBtn)

	// Button to clear all groups
	clearBtn := fltk.NewButton(360, 0, 120, 30, "Clear All")
	clearBtn.SetCallback(func() {
		fmt.Println("click")
		grid.ClearAll()
	})
	buttonRow.Add(clearBtn)
	buttonRow.End()
	win.Add(buttonRow)

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
