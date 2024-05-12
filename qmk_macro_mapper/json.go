package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type MacroAction struct {
	Action   string   `json:"action"`
	Keycodes []string `json:"keycodes"`
}

type Macro struct {
	Data interface{}
}

type Layout struct {
	Version       int             `json:"version,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	Documentation string          `json:"documentation,omitempty"`
	Keyboard      string          `json:"keyboard,omitempty"`
	Keymap        string          `json:"keymap,omitempty"`
	Layout        string          `json:"layout,omitempty"`
	Macros        [][]interface{} `json:"macros,omitempty"`
	Layers        [][]string      `json:"layers,omitempty"`
}

type MacroConfig struct {
	Name    string        `json:"name"`
	Keycode string        `json:"keycode"`
	Macro   []interface{} `json:"macro"`
}

func load_layout_json(file string) (layout Layout, err error) {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		return Layout{}, err
	}

	fmt.Println("Successfully Opened layout file", file)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Layout{}, err
	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &layout)

	return layout, nil
}

func load_macro_json(file string) (macro_config []MacroConfig, err error) {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		return []MacroConfig{}, err
	}

	fmt.Println("Successfully Opened macro config file", file)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return []MacroConfig{}, err
	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &macro_config)

	return macro_config, nil
}

func write_output_to_json(layout Layout, out_file string) error {
	jsonData, err := json.MarshalIndent(layout, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.Create(out_file)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
