package account

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAccountIdJsonUnmarshal(t *testing.T) {
	input := []byte(`{"A": "-1:5461696b6920697320666f722065726572796f6e6520696e20776562332e2e2e"}`)
	var a struct{ A AccountId }
	err := json.Unmarshal(input, &a)
	if err != nil {
		t.Fatal(err)
	}
	if a.A.Workchain != -1 {
		t.Fatal("invalid workchain")
	}
	for i, b := range []byte{84, 97, 105, 107, 105, 32, 105, 115, 32, 102, 111, 114, 32, 101, 114, 101, 114, 121, 111, 110, 101, 32, 105, 110, 32, 119, 101, 98, 51, 46, 46, 46} {
		if a.A.Address[i] != b {
			t.Fatal("invalid address")
		}
	}
}

func TestAccountIdString(t *testing.T) {
	addr := AccountId{
		Workchain: 0,
		Address:   [32]byte{84, 97, 105, 107, 105, 32, 105, 115, 32, 102, 111, 114, 32, 101, 114, 101, 114, 121, 111, 110, 101, 32, 105, 110, 32, 119, 101, 98, 51, 46, 46, 46},
	}
	if addr.String() != "0:5461696b6920697320666f722065726572796f6e6520696e20776562332e2e2e" {
		t.Fatal("AccountId.String() incorrect")
	}

	res, _ := addr.MarshalJSON()

	fmt.Printf("%x\n", res)
}
