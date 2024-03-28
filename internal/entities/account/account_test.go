package account

import (
	"fmt"
	"testing"
)

func TestAccount(t *testing.T) {
	type test struct {
		source        string
		expectedError bool
	}

	tests := []test{
		{
			source:        `{"id":0,"name":"BOG (GEL)","currency_code":"GEL"}`,
			expectedError: true,
		},
		{
			source:        `[{"id":0,"name":"BOG (GEL)","currency_code":"GEL"}]`,
			expectedError: false,
		},
		{
			source:        `[{"id":0,"name":"BOG (GEL)","currency_code":"GEL"},{"id":1,"name":"BOG (USD)","currency_code":"USD"}]`,
			expectedError: false,
		},
		{
			source:        `[]`,
			expectedError: false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint("test", i), func(t *testing.T) {
			dataJSON := []byte(tt.source)
			accounts, err := ParseJSON(dataJSON)
			if tt.expectedError {
				if err == nil {
					t.Fatal("Expected error")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				t.Log(accounts)
			}
		})
	}
}
