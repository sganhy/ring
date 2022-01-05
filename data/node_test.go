package data

import (
	"testing"
)

//test SetField(), GetField()
func Test__Node__setValue(t *testing.T) {
	nodeData := new(node)

	// test value from 0..255
	for i := 0; i < 256; i++ {
		// test default value
		if nodeData.getValue(uint8(i)) != false {
			t.Errorf("Node.getValue() ==> getValue(%d) should be equal false", i)
		}
		nodeData.setValue(uint8(i), true)
		if nodeData.getValue(uint8(i)) != true {
			t.Errorf("Node.getValue() ==> getValue(%d) should be equal true", i)
		}
		nodeData.setValue(uint8(i), false)
		if nodeData.getValue(uint8(i)) != false {
			t.Errorf("Node.getValue() ==> getValue(%d) should be equal false", i)
		}
		nodeData.setValue(uint8(255), false)
		if nodeData.getValue(uint8(255)) != false {
			t.Errorf("Node.getValue() ==> getValue(%d) should be equal false", 0)
		}
		nodeData.setValue(uint8(0), false)
		if nodeData.getValue(uint8(i)) != false {
			t.Errorf("Node.getValue() ==> getValue(0) should be equal false")
		}
		if nodeData.Count() != (i>>6)+1 {
			t.Errorf("Node.Count() ==> Count(%d) should be equal %d instead of %d",
				i, (i>>6)+1, nodeData.Count())
		}
	}
}
