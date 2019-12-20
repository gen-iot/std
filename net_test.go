package std

import (
	"encoding/json"
	"testing"
)

func TestListIf(t *testing.T) {
	netIfs, err := NetworkList(SkipDefault)
	AssertError(err, "list network")
	jstr, err := json.MarshalIndent(netIfs, KJsonIndentPrefix, KJsonIndent)
	AssertError(err, "marshal indent")
	t.Log(string(jstr))
	for _, net := range netIfs {
		t.Log("Name=", net.Name)
		t.Log("Mac=", net.Mac)
		for _, addr := range net.Address {
			t.Log("IP=", addr.IP)
			t.Log("Mask=", addr.Mask)
		}
		t.Log()
	}
}
