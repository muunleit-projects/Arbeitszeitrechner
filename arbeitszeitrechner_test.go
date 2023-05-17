package arbeitszeitrechner_test

import (
	"bytes"
	"strings"
	"testing"

	az "github.com/muunleit-projects/Arbeitszeitrechner"
)

func TestNew(t *testing.T) {
	t.Parallel()

	var zp az.Zeitpunkt
	err := zp.SetBeginn("6:14")
	if err != nil {
		t.Fatal(err)
	}

	want := "06:14"
	got := zp.Beginn()
	if !strings.HasPrefix(got, want) {
		t.Errorf("Beginn: want %q, got %q", want, got)
	}
}

func TestOutputTableToWriter(t *testing.T) {
	t.Parallel()

	var zp az.Zeitpunkt
	fakeTerminal := &bytes.Buffer{}
	err := zp.SetBeginn("7:45")
	if err != nil {
		t.Fatal(err)
	}
	zp.Tabelle(fakeTerminal)

	wantBeginn := "7:45  "
	wantStandard := "16:03  "
	wantMax := "18:30  "
	got := fakeTerminal.String()
	if !(strings.Contains(got, wantBeginn) &&
		strings.Contains(got, wantStandard) &&
		strings.Contains(got, wantMax)) {
		t.Errorf(got)
	}
}
