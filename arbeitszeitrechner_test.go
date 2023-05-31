package arbeitszeitrechner_test

import (
	// "bytes"
	// "strings"
	"strings"
	"testing"
	"time"

	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

func TestNew(t *testing.T) {
	t.Parallel()

	fakeInput := strings.NewReader("6:14")

	zp := azr.New()
	zp.SetCurrentTime(time.Date(2020, 11, 7, 13, 35, 0, 0, time.Local))
	zp.SetInput(fakeInput)
	err := zp.SetBeginn()
	if err != nil {
		t.Fatal(err)
	}

	want := time.Date(2020, 11, 7, 6, 14, 0, 0, time.Local)
	got := zp.Beginn()
	if want != got {
		t.Errorf("Beginn: want %v, got %v", want, got)
	}
}

// func TestOutputTableToWriter(t *testing.T) {
// 	t.Parallel()

// 	var zp azr.Zeitpunkt
// 	fakeTerminal := &bytes.Buffer{}
// 	err := zp.SetBeginn("7:45")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	zp.Tabelle(fakeTerminal)

// 	wantBeginn := "7:45  "
// 	wantStandard := "16:03  "
// 	wantMax := "18:30  "
// 	got := fakeTerminal.String()
// 	if !(strings.Contains(got, wantBeginn) &&
// 		strings.Contains(got, wantStandard) &&
// 		strings.Contains(got, wantMax)) {
// 		t.Errorf(got)
// 	}
// }
