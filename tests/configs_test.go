package tests

import (
	"encoder/configs"
	"fmt"
	"testing"
)

func TestGetConfigWithExtension(t *testing.T) {
	type configTest struct {
		input    []string
		expected *configs.Formats
	}

	tests := []configTest{
		{input: []string{"mp4", "ld"}, expected: &configs.Formats{Name: "mp4_ld", Extension: "mpv"}},
		{input: []string{"mp4", "uhd"}, expected: &configs.Formats{Name: "mp4_uhd", Extension: "mpv"}},
		{input: []string{"avi", "uld"}, expected: &configs.Formats{Name: "avi_uld", Extension: "avi"}},
	}

	for _, tc := range tests {
		testName := fmt.Sprintf("%v_%v", tc.input[0], tc.input[1])
		t.Run(testName, func(t *testing.T) {
			actual := configs.GetConfigWithExtension(tc.input[0], tc.input[1])
			if actual.Extension != tc.expected.Extension {
				t.Errorf("expected '%s' got '%s'", tc.expected.Extension, actual.Extension)
			}

			if actual.Name != testName {
				t.Fatalf("expected '%s' but got '%v'", testName, actual.Name)
			}
		})
	}
}
