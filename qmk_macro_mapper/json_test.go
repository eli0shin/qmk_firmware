package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFlattenLayersUsesExplicitLayerOrderAndPadsUnusedRows(t *testing.T) {
	source := map[string]SplitLayer{
		"1": splitLayer([][]string{{"L1L"}}, [][]string{{"L1R"}}),
		"0": splitLayer([][]string{{"L0L"}}, [][]string{{"L0R"}}),
	}

	got, err := flattenLayers(source, 2)
	if err != nil {
		t.Fatal(err)
	}

	want := [][]string{
		{"L0L", "L0R", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
		{"L1L", "L1R", "KC_NO", "KC_NO", "KC_NO", "KC_NO"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("flattenLayers() = %#v, want %#v", got, want)
	}
}

func TestFlattenLayerInterleavesHandsByRow(t *testing.T) {
	layer := splitLayer(
		[][]string{{"L00", "L01"}, {"L10"}},
		[][]string{{"R00", "R01"}, {"R10"}},
	)

	got, err := flattenLayer(layer, 0)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"L00", "L01", "R00", "R01", "L10", "R10"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("flattenLayer() = %#v, want %#v", got, want)
	}
}

func TestFlattenLayersRejectsMissingLayerNumber(t *testing.T) {
	_, err := flattenLayers(map[string]SplitLayer{
		"0": splitLayer([][]string{{"L"}}, [][]string{{"R"}}),
		"2": splitLayer([][]string{{"L"}}, [][]string{{"R"}}),
	}, 0)
	if err == nil || !strings.Contains(err.Error(), "layer 1 is missing") {
		t.Fatalf("flattenLayers() error = %v, want missing layer error", err)
	}
}

func TestFlattenLayerRejectsDifferentHandWidths(t *testing.T) {
	_, err := flattenLayer(splitLayer(
		[][]string{{"L0", "L1"}},
		[][]string{{"R0"}},
	), 0)
	if err == nil || !strings.Contains(err.Error(), "same non-zero width") {
		t.Fatalf("flattenLayer() error = %v, want row width error", err)
	}
}

func splitLayer(left, right [][]string) SplitLayer {
	return SplitLayer{
		Left:  Hand{Rows: left},
		Right: Hand{Rows: right},
	}
}
