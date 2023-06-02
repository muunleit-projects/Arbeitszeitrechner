package arbeitszeitrechner_test

import (
	"bytes"
	"testing"
	"time"

	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

func TestNew(t *testing.T) {
	t.Parallel()

	fakeOutput := &bytes.Buffer{}

	az := azr.Zeitpunkt{
		CurrentTime: time.Date(2020, 7, 23, 13, 6, 0, 0, time.Local),
		Output:      fakeOutput,
	}

	az.Tabelle("8:12")

	want := "Beginn                 08:12  Thu 23.07.2020" + "\n" +
		"Standard-Tag           16:30  Thu 23.07.2020    3h24m0s" + "\n" +
		"maximale Arbeitszeit   18:57  Thu 23.07.2020    5h51m0s" + "\n"
	got := fakeOutput.String()
	if want != got {
		t.Errorf("Tabelle: \nwant \n%v got \n%v", want, got)
	}
}
