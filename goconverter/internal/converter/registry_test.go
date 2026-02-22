package converter

import (
	"reflect"
	"slices"
	"testing"
)

func TestConversionTargetsBySource(t *testing.T) {
	output := ConversionTargetsBySource()

	expected := map[string][]string{
		"png":  {"jpg", "webp"},
		"jpg":  {"png", "webp"},
		"webp": {"jpg", "png"},
	}

	if len(output) != len(expected) {
		t.Fatalf("expected %d source formats, got %d", len(expected), len(output))
	}

	for source, expectedTargets := range expected {
		targets, ok := output[source]
		if !ok {
			t.Fatalf("expected output to include key %q, got %v", source, output)
		}

		slices.Sort(targets)
		slices.Sort(expectedTargets)
		if !reflect.DeepEqual(targets, expectedTargets) {
			t.Fatalf("expected output for %s to be %v, got %v", source, expectedTargets, targets)
		}
	}
}
