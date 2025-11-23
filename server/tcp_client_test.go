package server

import "testing"

func TestHostName(t *testing.T) {
	got, err := hostname()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(got)

}
