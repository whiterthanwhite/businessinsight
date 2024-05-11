package helper

import "testing"

func TestExportToCSV(t *testing.T) {
	values := [][]string{
		{"value 1", "value 2", "value 3"},
		{"value 4", "value 5", "value 6"},
	}

	err := exportToCSV("testfile", "test data exported", values)
	if err != nil {
		t.Fatal(err.Error())
	}
}
