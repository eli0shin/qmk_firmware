package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	layout_file, macros_file, out_file := get_arguments()

	layout, err := load_layout_json(layout_file)
	if err != nil {
		fmt.Println("error loading layout json", err)
		return
	}

	macro_config, err := load_macro_json(macros_file)
	if err != nil {
		fmt.Println("error loading macro json", err)
		return
	}

	macro_map, macro_array := transform_macros_to_map(macro_config)

	layout.Macros = macro_array

	map_macros_to_layout(layout.Layers, macro_map)

	output_err := write_output_to_json(layout, out_file)
	if output_err != nil {
		fmt.Println("Error writing out file to json", output_err)
		return
	}

	fmt.Println("Successfully wrote layout to", out_file)
}

func map_macros_to_layout(layers [][]string, macro_map map[string]IndexedMacro) {
	macro_regex := regexp.MustCompile(`^ANY\(M\*(.*?)\)`)
	for layer_index, layer := range layers {
		for key_index, key := range layer {

			fmt.Println("key", key)
			match := macro_regex.FindStringSubmatch(key)
			if match != nil {

				fmt.Println("match", match)
				macro, ok := macro_map[match[1]]
				fmt.Println("macro", macro)
				if ok {
					layers[layer_index][key_index] = fmt.Sprintf("ANY(QK_MACRO_%d)", macro.Index)
				} else {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("Cannot find macro by keycode %s\r\n", match[1]))
				}
			}
		}
	}
}

type IndexedMacro struct {
	Index int
	Macro []interface{}
}

func transform_macros_to_map(macro_config []MacroConfig) (macro_map map[string]IndexedMacro, macro_array [][]interface{}) {
	// Initialize the map
	macro_map = make(map[string]IndexedMacro)

	// Initialize the slice with the appropriate length, filled with nil slices
	macro_array = make([][]interface{}, len(macro_config))

	for index, macro := range macro_config {
		macro_map[macro.Keycode] = IndexedMacro{index, macro.Macro}
		macro_array[index] = macro.Macro
	}
	return macro_map, macro_array
}

func get_arguments() (layout_file string, macros_file string, out_file string) {
	args_without_prog := os.Args[1:]
	return args_without_prog[0], args_without_prog[1], args_without_prog[2]
}
