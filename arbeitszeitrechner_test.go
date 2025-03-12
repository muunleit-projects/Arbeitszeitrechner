package arbeitszeitrechner_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

// TestNew tests the NewArbeitszeitrechner function.
func TestNew(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		explanation    string
		input          string
		OutputExpected string
		errExpected    error
	}{
		{
			explanation: "Standard case",
			input:       "8:12",
			OutputExpected: "Beginn                 08:12  Thu 23.07.2020" + "\n" +
				"Standard-Tag           16:30  Thu 23.07.2020    3h24m0s" + "\n" +
				"maximale Arbeitszeit   18:57  Thu 23.07.2020    5h51m0s" + "\n",
			errExpected: nil,
		},
		{
			explanation: "Early start",
			input:       "6:00",
			OutputExpected: "Beginn                 06:00  Thu 23.07.2020" + "\n" +
				"Standard-Tag           14:18  Thu 23.07.2020    1h12m0s" + "\n" +
				"maximale Arbeitszeit   16:45  Thu 23.07.2020    3h39m0s" + "\n",
			errExpected: nil,
		},
		{
			explanation: "Late start",
			input:       "10:00",
			OutputExpected: "Beginn                 10:00  Thu 23.07.2020" + "\n" +
				"Standard-Tag           18:18  Thu 23.07.2020    5h12m0s" + "\n" +
				"maximale Arbeitszeit   20:45  Thu 23.07.2020    7h39m0s" + "\n",
			errExpected: nil,
		},
		{
			explanation: "Night start",
			input:       "22:12",
			OutputExpected: "Beginn                 22:12  Wed 22.07.2020" + "\n" +
				"Standard-Tag           06:30  Thu 23.07.2020" + "\n" +
				"maximale Arbeitszeit   08:57  Thu 23.07.2020" + "\n",
			errExpected: nil,
		},
	} {
		t.Run(fmt.Sprintf("%s [%s]", tt.explanation, tt.input), func(t *testing.T) {
			fakeOutput := &bytes.Buffer{}
			az, err := azr.NewArbeitszeitrechner(
				azr.Now(time.Date(2020, 7, 23, 13, 6, 0, 0, time.Local)),
				azr.Output(fakeOutput),
			)

			if got, want := err, tt.errExpected; got != want {
				t.Fatalf("err=%v, want=%v", got, want)
			}

			az.Tabelle(tt.input)

			if got, want := fakeOutput.String(), tt.OutputExpected; got != want {
				t.Errorf("\ngot=\n%q \nwant=\n%q ", got, want)
			}
		})
	}
}
