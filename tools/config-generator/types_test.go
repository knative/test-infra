package main

import "testing"

func TestJobDetailMap(t *testing.T) {
	j := NewJobDetailMap()

	local := []string{"continuous", "nightly"}

	for _, t := range local {
		j.Add("serving", t)
	}

	for i := range local {
		if j["serving"][i] == local[i] {
			t.Logf("Entry %d matched", i)
		} else {
			t.Errorf("Entry %d did not match: %q != %q", i, j["serving"][i], local[i])
		}
	}
}
