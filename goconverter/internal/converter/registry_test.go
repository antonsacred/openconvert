package converter

import (
	"reflect"
	"slices"
	"testing"

	"github.com/h2non/bimg"
)

func TestConversionTargetsBySource(t *testing.T) {
	output := ConversionTargetsBySource()
	expected := expectedConversionTargetsBySource()

	if !reflect.DeepEqual(normalizeTargetsMap(output), normalizeTargetsMap(expected)) {
		t.Fatalf("unexpected conversion matrix.\nexpected: %v\ngot: %v", normalizeTargetsMap(expected), normalizeTargetsMap(output))
	}
}

func TestFindConverter(t *testing.T) {
	expected := expectedConversionTargetsBySource()

	for source, targets := range expected {
		for _, target := range targets {
			_, ok := FindConverter(source, target)
			if !ok {
				t.Fatalf("expected converter for %s -> %s", source, target)
			}
		}
	}

	for _, format := range []string{"avif", "gif", "heif", "jpeg", "png", "tiff", "webp"} {
		if !bimg.IsTypeNameSupported(format) {
			continue
		}
		if !bimg.IsTypeNameSupportedSave(format) {
			continue
		}

		_, ok := FindConverter(format, format)
		if ok {
			t.Fatalf("did not expect converter for same-format conversion %s -> %s", format, format)
		}
	}
}

func expectedConversionTargetsBySource() map[string][]string {
	loadSupportedSources := []string{"avif", "gif", "heif", "jpeg", "png", "tiff", "webp", "magick", "pdf", "svg"}
	saveSupportedTargets := []string{"avif", "gif", "heif", "jpeg", "png", "tiff", "webp"}

	expected := make(map[string][]string)
	for _, source := range loadSupportedSources {
		if !bimg.IsTypeNameSupported(source) {
			continue
		}

		for _, target := range saveSupportedTargets {
			if source == target {
				continue
			}
			if !bimg.IsTypeNameSupportedSave(target) {
				continue
			}

			expected[source] = append(expected[source], target)
		}
	}

	return expected
}

func normalizeTargetsMap(input map[string][]string) map[string][]string {
	output := make(map[string][]string, len(input))
	for source, targets := range input {
		cloned := slices.Clone(targets)
		slices.Sort(cloned)
		output[source] = cloned
	}

	return output
}
