package tea

import (
	"image/color"
	"reflect"
	"testing"
)

func TestParseSequence_Events(t *testing.T) {
	input := []byte("\x1b\x1b[Ztest\x00\x1b]10;rgb:1234/1234/1234\x07\x1b[27;2;27~\x1b[?1049;2$y")
	want := []Msg{
		KeyPressMsg{Code: KeyTab, Mod: ModShift | ModAlt},
		KeyPressMsg{Code: 't', Text: "t"},
		KeyPressMsg{Code: 'e', Text: "e"},
		KeyPressMsg{Code: 's', Text: "s"},
		KeyPressMsg{Code: 't', Text: "t"},
		KeyPressMsg{Code: KeySpace, Mod: ModCtrl},
		ForegroundColorMsg{color.RGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xff}},
		KeyPressMsg{Code: KeyEscape, Mod: ModShift},
		modeReportMsg{Mode: 1049, Value: 2},
	}

	var p inputParser
	for i := 0; len(input) != 0; i++ {
		if i >= len(want) {
			t.Fatalf("reached end of want events")
		}
		n, got := p.parseSequence(input)
		if !reflect.DeepEqual(got, want[i]) {
			t.Errorf("got %#v (%T), want %#v (%T)", got, got, want[i], want[i])
		}
		input = input[n:]
	}
}

func BenchmarkParseSequence(b *testing.B) {
	var p inputParser
	input := []byte("\x1b\x1b[Ztest\x00\x1b]10;1234/1234/1234\x07\x1b[27;2;27~")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.parseSequence(input)
	}
}
