package arbeitszeitrechner_test

import (
	"bytes"
	"testing"
	"time"

	azr "github.com/muunleit-projects/Arbeitszeitrechner"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		start string
		want  string
	}{
		{
			name:  "Standard case",
			start: "8:12",
			want: "Beginn                 08:12  Thu 23.07.2020" + "\n" +
				"Standard-Tag           16:30  Thu 23.07.2020    3h24m0s" + "\n" +
				"maximale Arbeitszeit   18:57  Thu 23.07.2020    5h51m0s" + "\n",
		},
		{
			name:  "Early start",
			start: "6:00",
			want: "Beginn                 06:00  Thu 23.07.2020" + "\n" +
				"Standard-Tag           14:18  Thu 23.07.2020    1h12m0s" + "\n" +
				"maximale Arbeitszeit   16:45  Thu 23.07.2020    3h39m0s" + "\n",
		},
		{
			name:  "Late start",
			start: "10:00",
			want: "Beginn                 10:00  Thu 23.07.2020" + "\n" +
				"Standard-Tag           18:18  Thu 23.07.2020    5h12m0s" + "\n" +
				"maximale Arbeitszeit   20:45  Thu 23.07.2020    7h39m0s" + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeOutput := &bytes.Buffer{}

			az, err := azr.NewArbeitszeitrechner(
				azr.Now(time.Date(2020, 7, 23, 13, 6, 0, 0, time.Local)),
				azr.Output(fakeOutput),
			)
			if err != nil {
				t.Fatal(err)
			}

			az.Tabelle(tt.start)

			got := fakeOutput.String()

			if tt.want != got {
				t.Errorf("Tabelle: \nwant \n%v got \n%v", tt.want, got)
			}
		})
	}
}
