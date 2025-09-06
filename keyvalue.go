package keyvalue

import (
	"strings"

	fltk "github.com/archeopternix/go-fltk"
)

// --- Data structures ---

// GroupedKeyValueTree represents the main data structure for storing groups of key-value pairs.
// Each group (KVPgroup) contains multiple KVPelements (key-value pairs).
type GroupedKeyValueTree struct {
	Groups []KVPgroup
}

// KVPgroup represents a group of key-value pairs, identified by Name.
type KVPgroup struct {
	Name     string
	Elements []KVPelement
}

// KVPelement represents a single key-value pair within a group.
type KVPelement struct {
	Key   string
	Value string
}

// groupWidgets aggregates the FLTK widgets for a single group displayed in the grid.
type groupWidgets struct {
	box      *fltk.Group   // The surrounding FLTK group box for the group.
	label    *fltk.Box     // The FLTK label for the group name.
	inputs   []*fltk.Input // The value input widgets for all keys in this group.
	elements *[]KVPelement // Pointer to the underlying elements in GroupedKeyValueTree.
}

// KeyValueGrid displays and manages a grid of editable grouped key-value pairs
// using FLTK widgets. It owns the underlying GroupedKeyValueTree structure.
type KeyValueGrid struct {
	*fltk.Group
	tree         *GroupedKeyValueTree     // Underlying data model.
	groupWidgets map[string]*groupWidgets // Mapping from group name to their widgets.
	nextY        int                      // Vertical position for next group.
	parent       *fltk.Window             // Parent FLTK window for widget placement.
	x, y, w, h   int                      // Geometry for the grid area.
}

// NewKeyValueGrid creates and initializes a KeyValueGrid within the given parent window
// at position (x, y) with size (w, h).
func NewKeyValueGrid(parent *fltk.Window, x, y, w, h int) *KeyValueGrid {
	group := fltk.NewGroup(x, y, w, h)
	grid := &KeyValueGrid{
		Group:        group, // Embed the group
		tree:         &GroupedKeyValueTree{},
		groupWidgets: make(map[string]*groupWidgets),
		nextY:        y + 15,
		parent:       parent,
		x:            x,
		y:            y,
		w:            w,
		h:            h,
	}
	parent.Add(group)

	return grid
}

// Resize updates the geometry of the grid and triggers a redraw/layout.
func (grid *KeyValueGrid) Resize(x, y, w, h int) {
	grid.x = x
	grid.y = y
	grid.w = w
	grid.h = h
	grid.Refresh()
}

// Refresh rebuilds the visible grid from the data model, recalculating key column width
// to fit the longest key text in all groups. All widgets are re-created and re-laid-out.
func (grid *KeyValueGrid) Refresh() {
	// Find the longest key
	longestKeyText := ""

	for _, group := range grid.tree.Groups {
		for _, elem := range group.Elements {
			if len(elem.Key) > len(longestKeyText) {
				longestKeyText = elem.Key
			}
		}
	}

	// Get the width in pixels for the longest key
	keyWidth := 70 // Default minimum width

	if longestKeyText != "" {
		//w, _ := fltk.MeasureText(longestKeyText, true)
		w := len(longestKeyText) * 6
		keyWidth = w + 16 // Add padding for readability
		if keyWidth < 70 {
			keyWidth = 70
		}
	}

	// Clean up old widgets
	for _, gw := range grid.groupWidgets {
		gw.box.Hide()
		grid.parent.Remove(gw.box)
		gw.label.Hide()
		grid.parent.Remove(gw.label)
		for _, in := range gw.inputs {
			in.Hide()
			grid.parent.Remove(in)
		}
	}

	grid.groupWidgets = make(map[string]*groupWidgets)
	grid.nextY = grid.y + 15

	// Layout all groups within the area
	for i := range grid.tree.Groups {
		group := &grid.tree.Groups[i]
		if len(group.Elements) == 0 {
			continue
		}
		gw := grid.drawGroup(group, keyWidth)
		grid.groupWidgets[group.Name] = gw
	}
	grid.parent.Redraw()
}

// drawGroup creates and places the widgets for a single group in the grid.
// keyWidth determines the width of the key label column.
func (grid *KeyValueGrid) drawGroup(group *KVPgroup, keyWidth int) *groupWidgets {
	const labelHeight = 24
	const inputHeight = 25
	const inputPad = 4

	boxWidth := grid.w - 2*20
	boxLeft := grid.x + 20

	numElems := len(group.Elements)
	boxHeight := labelHeight/2 + numElems*inputHeight + (numElems+1)*inputPad

	box := fltk.NewGroup(boxLeft, grid.nextY+labelHeight/4, boxWidth, boxHeight, "")
	box.SetBox(fltk.EMBOSSED_BOX)
	box.Begin()

	//	ToDo: crash: labelwidth, _ := fltk.MeasureText(group.Name, false)
	labelwidth := len(group.Name) * 6
	label := fltk.NewBox(fltk.FLAT_BOX, boxLeft+10, grid.nextY-7, labelwidth+20, labelHeight, group.Name)
	label.SetAlign(fltk.ALIGN_LEFT | fltk.ALIGN_INSIDE)
	label.SetLabelSize(10)

	inputs := []*fltk.Input{}
	for i := range group.Elements {
		elem := &group.Elements[i]
		y := grid.nextY + labelHeight/2 + inputPad + i*(inputHeight+inputPad)
		keyLabel := fltk.NewBox(fltk.NO_BOX, boxLeft+15, y, keyWidth, inputHeight, elem.Key)
		keyLabel.SetAlign(fltk.ALIGN_INSIDE | fltk.ALIGN_LEFT)
		keyLabel.SetLabelSize(12)

		valInput := fltk.NewInput(boxLeft+15+keyWidth+5, y, boxWidth-(keyWidth+35), inputHeight)
		valInput.SetValue(elem.Value)
		inputs = append(inputs, valInput)
	}

	box.End()

	grid.parent.Add(box)
	grid.parent.Add(label)

	gw := &groupWidgets{
		box:      box,
		label:    label,
		inputs:   inputs,
		elements: &group.Elements,
	}
	grid.nextY += boxHeight + labelHeight/2 + 10
	return gw
}

// Add inserts or updates a key-value pair in the specified group.
// If groupName does not exist, it is created. If key does not exist in group, it is added.
// If key exists, its value is updated (if value is non-empty).
// If only groupName is given, creates group if not present.
func (grid *KeyValueGrid) Add(groupName, key, value string) bool {
	groupName = strings.TrimSpace(groupName)
	key = strings.TrimSpace(key)

	if groupName == "" {
		return false
	}
	// Find or create group
	var group *KVPgroup
	for i := range grid.tree.Groups {
		if grid.tree.Groups[i].Name == groupName {
			group = &grid.tree.Groups[i]
			break
		}
	}

	if group == nil {
		grid.tree.Groups = append(grid.tree.Groups, KVPgroup{Name: groupName})
		group = &grid.tree.Groups[len(grid.tree.Groups)-1]
		if key == "" {
			grid.Refresh()
			return true
		}
	}
	if key == "" {
		grid.Refresh()
		return true
	}

	// Find or add key
	for i := range group.Elements {
		if group.Elements[i].Key == key {
			if value != "" {
				group.Elements[i].Value = value
			}
			grid.Refresh()
			return true
		}
	}

	group.Elements = append(group.Elements, KVPelement{Key: key, Value: value})

	grid.Refresh()

	return true

}

// Delete removes an entire group (if key is empty) or a specific key in the group.
// Returns true if something was deleted.
func (grid *KeyValueGrid) Delete(groupName, key string) bool {
	groupName = strings.TrimSpace(groupName)
	key = strings.TrimSpace(key)
	for gi := range grid.tree.Groups {
		group := &grid.tree.Groups[gi]
		if group.Name == groupName {
			if key == "" {
				grid.tree.Groups = append(grid.tree.Groups[:gi], grid.tree.Groups[gi+1:]...)
				grid.Refresh()
				return true
			}
			for ei, elem := range group.Elements {
				if elem.Key == key {
					group.Elements = append(group.Elements[:ei], group.Elements[ei+1:]...)
					grid.Refresh()
					return true
				}
			}
			break
		}
	}
	return false
}

// GetData synchronizes user-edited values from the grid UI to the data model,
// and returns a pointer to the root GroupedKeyValueTree.
func (grid *KeyValueGrid) GetData() *GroupedKeyValueTree {
	for _, gw := range grid.groupWidgets {
		for i, input := range gw.inputs {
			if i < len(*gw.elements) {
				(*gw.elements)[i].Value = input.Value()
			}
		}
	}
	return grid.tree
}

func (grid *KeyValueGrid) ClearAll() {
	// Remove all widgets from the parent window and hide them
	for _, gw := range grid.groupWidgets {
		gw.box.Hide()
		grid.parent.Remove(gw.box)
		gw.label.Hide()
		grid.parent.Remove(gw.label)
		for _, in := range gw.inputs {
			in.Hide()
			grid.parent.Remove(in)
		}
	}
	// Clear widget map
	grid.groupWidgets = make(map[string]*groupWidgets)
	// Clear the underlying data model
	grid.tree.Groups = nil
	// Reset vertical position
	grid.nextY = grid.y + 15
	// Redraw the parent window to update the UI
	grid.parent.Redraw()
}
