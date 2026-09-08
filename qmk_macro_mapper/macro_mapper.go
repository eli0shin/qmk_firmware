package main

import (
	"fmt"
	"os"
	"regexp"
)

var macroKeycodePattern = regexp.MustCompile(`^ANY\(M\*(.*?)\)$`)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "macro_mapper:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("usage: macro_mapper <layout.json> <macros.json> <output.json>")
	}

	layout, err := loadLayoutJSON(args[0])
	if err != nil {
		return fmt.Errorf("load layout %q: %w", args[0], err)
	}

	macroConfig, err := loadMacroJSON(args[1])
	if err != nil {
		return fmt.Errorf("load macros %q: %w", args[1], err)
	}

	macroMap, macroArray := transformMacrosToMap(macroConfig)
	layout.Macros = macroArray
	if err := mapMacrosToLayout(layout.Layers, macroMap); err != nil {
		return err
	}

	if err := writeOutputToJSON(layout, args[2]); err != nil {
		return fmt.Errorf("write output %q: %w", args[2], err)
	}

	fmt.Println("Successfully wrote layout to", args[2])
	return nil
}

func mapMacrosToLayout(layers [][]string, macroMap map[string]IndexedMacro) error {
	for layerIndex, layer := range layers {
		for keyIndex, key := range layer {
			match := macroKeycodePattern.FindStringSubmatch(key)
			if match == nil {
				continue
			}

			macro, ok := macroMap[match[1]]
			if !ok {
				return fmt.Errorf("layer %d key %d refers to unknown macro %q", layerIndex, keyIndex, match[1])
			}
			layers[layerIndex][keyIndex] = fmt.Sprintf("ANY(QK_MACRO_%d)", macro.Index)
		}
	}
	return nil
}

type IndexedMacro struct {
	Index int
	Macro []interface{}
}

func transformMacrosToMap(macroConfig []MacroConfig) (map[string]IndexedMacro, [][]interface{}) {
	macroMap := make(map[string]IndexedMacro)
	macroArray := make([][]interface{}, len(macroConfig))

	for index, macro := range macroConfig {
		macroMap[macro.Keycode] = IndexedMacro{Index: index, Macro: macro.Macro}
		macroArray[index] = macro.Macro
	}
	return macroMap, macroArray
}
