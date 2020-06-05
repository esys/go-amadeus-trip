package amadeus

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func readTestData(file string, t *testing.T) []byte {
	f, err := os.Open(file)
	if err != nil {
		t.Errorf("failed to load testdata %s: %v", file, err)
		return nil
	}
	byt, err := ioutil.ReadAll(f)
	if err != nil {
		t.Errorf("failed to load testdata %s: %v", file, err)
		return nil
	}
	return byt
}

func TestUnmarshall(t *testing.T) {
	tests := []struct {
		name    string
		payload []byte
		want    interface{}
		wantErr bool
	}{
		{
			name:    "read air product JSON",
			payload: []byte(readTestData("testdata/air.json", t)),
			want:    &resultResponseData{},
			wantErr: false,
		},
		{
			name:    "read hotel product JSON",
			payload: []byte(readTestData("testdata/hotel.json", t)),
			want:    &resultResponseData{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal(tt.payload, tt.want); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
