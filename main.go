package main

import (
	"fmt"

	fltk "github.com/pwiecz/go-fltk"
)

func main() {
	win := fltk.NewWindow(500, 420)
	win.SetLabel("KeyValueGrid Demo")

	fmt.Println("check1")
	grid := NewKeyValueGrid(win, 20, 20, 460, 350)
	fmt.Println("check2")

	// Add some groups and key-value pairs
	grid.Add("User", "Name", "Alice")
	grid.Add("User", "Email", "alice@example.com")
	grid.Add("Settings", "Theme", "Dark")
	grid.Add("Settings", "Language", "en-US")
	grid.Add("Network", "Proxy", "")
	grid.Add("Network", "Timeout", "30s")
	fmt.Println("check3")

	/*	// Button to print data from UI edits
		printBtn := fltk.NewButton(20, 380, 120, 30, "Print Data")
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
	*/
	// Button to clear all groups
	clearBtn := fltk.NewButton(360, 380, 120, 30, "Clear All")
	clearBtn.SetCallback(func() {
		grid.ClearAll()
	})

	win.End()
	win.Show()
	fltk.Run()
}
