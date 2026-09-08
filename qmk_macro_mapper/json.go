package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Layout struct {
	Version       int             `json:"version,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	Documentation string          `json:"documentation,omitempty"`
	Author        string          `json:"author,omitempty"`
	Keyboard      string          `json:"keyboard,omitempty"`
	Keymap        string          `json:"keymap,omitempty"`
	Layout        string          `json:"layout,omitempty"`
	Macros        [][]interface{} `json:"macros,omitempty"`
	Layers        [][]string      `json:"layers,omitempty"`
}

type LayoutSource struct {
	Version       int                   `json:"version,omitempty"`
	Notes         string                `json:"notes,omitempty"`
	Documentation string                `json:"documentation,omitempty"`
	Author        string                `json:"author,omitempty"`
	Keyboard      string                `json:"keyboard"`
	Keymap        string                `json:"keymap"`
	Layout        string                `json:"layout"`
	UnusedRows    int                   `json:"unused_rows,omitempty"`
	Layers        map[string]SplitLayer `json:"layers"`
}

type SplitLayer struct {
	Left  Hand `json:"left"`
	Right Hand `json:"right"`
}

type Hand struct {
	Rows [][]string `json:"rows"`
}

type MacroConfig struct {
	Name    string        `json:"name"`
	Keycode string        `json:"keycode"`
	Macro   []interface{} `json:"macro"`
}

func loadLayoutJSON(file string) (Layout, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return Layout{}, err
	}

	var source LayoutSource
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&source); err != nil {
		return Layout{}, fmt.Errorf("parse layout JSON: %w", err)
	}
	if source.Keyboard == "" || source.Keymap == "" || source.Layout == "" {
		return Layout{}, fmt.Errorf("keyboard, keymap, and layout are required")
	}
	if source.UnusedRows < 0 {
		return Layout{}, fmt.Errorf("unused_rows cannot be negative")
	}

	layers, err := flattenLayers(source.Layers, source.UnusedRows)
	if err != nil {
		return Layout{}, err
	}

	return Layout{
		Version:       source.Version,
		Notes:         source.Notes,
		Documentation: source.Documentation,
		Author:        source.Author,
		Keyboard:      source.Keyboard,
		Keymap:        source.Keymap,
		Layout:        source.Layout,
		Layers:        layers,
	}, nil
}

func flattenLayers(source map[string]SplitLayer, unusedRows int) ([][]string, error) {
	if len(source) == 0 {
		return nil, fmt.Errorf("at least one layer is required")
	}

	layerNumbers := make([]int, 0, len(source))
	for key := range source {
		number, err := strconv.Atoi(key)
		if err != nil || number < 0 || strconv.Itoa(number) != key {
			return nil, fmt.Errorf("layer key %q must be a non-negative integer", key)
		}
		layerNumbers = append(layerNumbers, number)
	}
	sort.Ints(layerNumbers)

	layers := make([][]string, len(layerNumbers))
	expectedKeys := 0
	for position, number := range layerNumbers {
		if number != position {
			return nil, fmt.Errorf("layers must be contiguous from 0; layer %d is missing", position)
		}

		layer, err := flattenLayer(source[strconv.Itoa(number)], unusedRows)
		if err != nil {
			return nil, fmt.Errorf("layer %d: %w", number, err)
		}
		if position == 0 {
			expectedKeys = len(layer)
		} else if len(layer) != expectedKeys {
			return nil, fmt.Errorf("layer %d has %d keys; expected %d", number, len(layer), expectedKeys)
		}
		layers[position] = layer
	}

	return layers, nil
}

func flattenLayer(layer SplitLayer, unusedRows int) ([]string, error) {
	leftRows := layer.Left.Rows
	rightRows := layer.Right.Rows
	if len(leftRows) == 0 {
		return nil, fmt.Errorf("each hand must contain at least one row")
	}
	if len(leftRows) != len(rightRows) {
		return nil, fmt.Errorf("left has %d rows and right has %d", len(leftRows), len(rightRows))
	}

	keys := make([]string, 0)
	for rowIndex := range leftRows {
		left := leftRows[rowIndex]
		right := rightRows[rowIndex]
		if len(left) == 0 || len(left) != len(right) {
			return nil, fmt.Errorf("row %d must have the same non-zero width on both hands", rowIndex)
		}
		for _, key := range append(append([]string{}, left...), right...) {
			if strings.TrimSpace(key) == "" {
				return nil, fmt.Errorf("row %d contains an empty keycode", rowIndex)
			}
			keys = append(keys, key)
		}
	}

	unusedRowWidth := len(leftRows[len(leftRows)-1])
	for range unusedRows {
		for range unusedRowWidth * 2 {
			keys = append(keys, "KC_NO")
		}
	}

	return keys, nil
}

func loadMacroJSON(file string) ([]MacroConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var macros []MacroConfig
	if err := json.Unmarshal(data, &macros); err != nil {
		return nil, fmt.Errorf("parse macro JSON: %w", err)
	}
	return macros, nil
}

func writeOutputToJSON(layout Layout, file string) error {
	data, err := json.MarshalIndent(layout, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0o644)
}
