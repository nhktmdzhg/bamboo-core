/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Nguyễn Hoàng Kỳ  <nhktmdzhg@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENSE>.
 */

package bamboo

import (
	"testing"
)

func TestRebuildFromText_SimpleASCII(t *testing.T) {
	composition := RebuildCompositionFromText("goo", true)
	result := Flatten(composition, VietnameseMode)
	if result != "goo" {
		t.Errorf("RebuildFromText('goo') = %q, want %q", result, "goo")
	}
}

func TestRebuildFromText_SimpleVietnamese(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"chào", "chào"},
		{"việt", "việt"},
		{"google", "google"},
		{"đường", "đường"},
		{"người", "người"},
		{"as", "as"},
		{"được", "được"},
		{"những", "những"},
		{"ước", "ước"},
		{"ươi", "ươi"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			composition := RebuildCompositionFromText(tt.input, true)
			result := Flatten(composition, VietnameseMode)
			if result != tt.want {
				t.Errorf("RebuildFromText(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestRebuildFromText_UpperCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Việt", "Việt"},
		{"OO", "OO"},
		{"DD", "DD"},
		{"Nội", "Nội"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			composition := RebuildCompositionFromText(tt.input, true)
			result := Flatten(composition, VietnameseMode)
			if result != tt.want {
				t.Errorf("RebuildFromText(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

func TestRebuildFromText_ThenProcessKey(t *testing.T) {
	im := ParseInputMethod(GetInputMethodDefinitions(), "Telex")
	engine := NewEngine(im, EstdFlags)

	engine.RebuildEngineFromText("go")
	result := engine.GetProcessedString(VietnameseMode)
	if result != "go" {
		t.Errorf("After rebuild 'go', GetProcessedString = %q, want %q", result, "go")
	}

	engine.ProcessKey('s', VietnameseMode)
	result = engine.GetProcessedString(VietnameseMode)
	if result != "gó" {
		t.Errorf("After rebuild 'go' + ProcessKey('s'), got %q, want %q", result, "gó")
	}
}

func TestRebuildFromText_CompareWithProcessString(t *testing.T) {
	im := ParseInputMethod(GetInputMethodDefinitions(), "Telex")
	engine := NewEngine(im, EstdFlags)

	engine.ProcessString("goo", VietnameseMode)
	buggyResult := engine.GetProcessedString(VietnameseMode)

	engine.Reset()
	engine.RebuildEngineFromText("goo")
	correctResult := engine.GetProcessedString(VietnameseMode)

	if buggyResult == "goo" {
		t.Log("ProcessString('goo') unexpectedly correct - Telex behavior may have changed")
	}
	if correctResult != "goo" {
		t.Errorf("RebuildFromText('goo') = %q, want 'goo'", correctResult)
	}
	t.Logf("ProcessString('goo') = %q, RebuildFromText('goo') = %q", buggyResult, correctResult)
}
