/*
 * Bamboo - A Vietnamese Input method editor
 * Copyright (C) Nguyễn Hoàng Kỳ  <nhktmdzhg@gmail.com>
 *
 * This software is licensed under the MIT license. For more information,
 * see <https://github.com/BambooEngine/bamboo-core/blob/master/LICENSE>.
 */

package bamboo

import (
	"unicode"
)

// RebuildCompositionFromText creates a composition (list of Transformations) directly
// from a Vietnamese Unicode string, bypassing all Input Method rules.
func RebuildCompositionFromText(text string, stdStyle bool) []*Transformation {
	var composition []*Transformation

	for _, ch := range text {
		lowerCh := unicode.ToLower(ch)
		isUpperCase := unicode.IsUpper(ch)

		// Decompose the character into root + mark + tone
		tone := FindToneFromChar(lowerCh)
		mark, hasMark := FindMarkFromChar(lowerCh)

		// Get the root character (no tone, no mark)
		rootChar := lowerCh
		if tone != ToneNone {
			rootChar = AddToneToChar(rootChar, 0)
		}
		if hasMark && mark != MarkNone {
			rootChar = AddMarkToChar(rootChar, 0)
		}

		// Create the base Appending transformation for the root character
		appendTrans := &Transformation{
			IsUpperCase: isUpperCase,
			Rule: Rule{
				Key:        rootChar,
				EffectOn:   rootChar,
				EffectType: Appending,
				Result:     rootChar,
			},
		}
		composition = append(composition, appendTrans)

		if hasMark && mark != MarkNone {
			markTrans := &Transformation{
				Target: appendTrans,
				Rule: Rule{
					Key:        0,
					EffectType: MarkTransformation,
					Effect:     uint8(mark),
					EffectOn:   rootChar,
					Result:     AddMarkToTonelessChar(rootChar, uint8(mark)),
				},
			}
			composition = append(composition, markTrans)
		}
	}

	// Apply tones
	var lastTone Tone = ToneNone
	for _, ch := range text {
		lowerCh := unicode.ToLower(ch)
		t := FindToneFromChar(lowerCh)
		if t != ToneNone {
			lastTone = t
		}
	}

	if lastTone != ToneNone {
		toneTarget := findToneTarget(composition, stdStyle)
		if toneTarget != nil {
			toneTrans := &Transformation{
				Target: toneTarget,
				Rule: Rule{
					Key:        0,
					EffectType: ToneTransformation,
					Effect:     uint8(lastTone),
				},
			}
			composition = append(composition, toneTrans)
		}
	}

	return composition
}

// RebuildEngineFromText resets the engine and rebuilds its internal composition
// state from the given Vietnamese Unicode text, bypassing all IM rules.
func (e *BambooEngine) RebuildEngineFromText(text string) {
	e.Reset()
	e.composition = RebuildCompositionFromText(text, e.flags&EstdToneStyle != 0)
}
