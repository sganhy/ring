package data

import (
	"testing"
)

//test SetField(), GetField()
func Test__Node__setValue(t *testing.T) {
	nodeData := new(bitchain)

	// test value from 0..255
	for i := 0; i < 256; i++ {
		// test default value
		if nodeData.GetValue(uint8(i)) != false {
			t.Errorf("Node.GetValue() ==> GetValue(%d) should be equal to false", i)
		}
		nodeData.SetValue(uint8(i), true)
		if nodeData.GetValue(uint8(i)) != true {
			t.Errorf("Node.GetValue() ==> GetValue(%d) should be equal true", i)
		}
		nodeData.SetValue(uint8(i), false)
		if nodeData.GetValue(uint8(i)) != false {
			t.Errorf("Node.GetValue() ==> GetValue(%d) should be equal false", i)
		}
		nodeData.SetValue(uint8(255), false)
		if nodeData.GetValue(uint8(255)) != false {
			t.Errorf("Node.GetValue() ==> GetValue(%d) should be equal false", 0)
		}
		nodeData.SetValue(uint8(0), false)
		if nodeData.GetValue(uint8(i)) != false {
			t.Errorf("Node.GetValue() ==> GetValue(0) should be equal false")
		}
		if nodeData.Count() != (i>>6)+1 {
			t.Errorf("Node.Count() ==> Count(%d) should be equal %d instead of %d",
				i, (i>>6)+1, nodeData.Count())
		}
	}
}

func Test__Node__setAll(t *testing.T) {
	nodeData := new(bitchain)
	nodeData.ResetAll(255, true)
	for i := 0; i < 256; i++ {
		if nodeData.GetValue(uint8(i)) != true {
			t.Errorf("Node.setAll(255, true) ==> getValue(%d) should be equal to true", i)
		}
	}
	nodeData.ResetAll(255, false)
	for i := 0; i < 256; i++ {
		if nodeData.GetValue(uint8(i)) != false {
			t.Errorf("Node.setAll(255, true) ==> getValue(%d) should be equal to false", i)
		}
	}
}

func Test__Node__CountSetBits(t *testing.T) {
	nodeData := new(bitchain)
	// test value from 0..255
	// increasing number of bits
	for i := 0; i < 256; i++ {
		if nodeData.CountSetBits() != i {
			t.Errorf("Node.CountSetBits() ==> {1} should be equal to %d instead of %d", i, nodeData.CountSetBits())
		}
		nodeData.SetValue(uint8(i), true)
		if nodeData.CountSetBits() != i+1 {
			t.Errorf("Node.CountSetBits() ==> {2} should be equal to %d instead of %d", i+1, nodeData.CountSetBits())
		}
	}
	// decreasing number of bits
	for i := 255; i >= 0; i-- {
		nodeData.SetValue(uint8(i), false)
		if nodeData.CountSetBits() != i {
			t.Errorf("Node.CountSetBits() ==> should be equal to %d instead of %d", i, nodeData.CountSetBits())
		}
	}
}

func Test__Node__NodeByIndex(t *testing.T) {
	nodeData := new(bitchain)
	nodeData.ResetAll(255, true)

	//t.Errorf("nodeData.Count()  %d", nodeData.Count())

	if nodeData.NodeByIndex(0) == nil {
		t.Errorf("Node.NodeByIndex() ==> NodeByIndex(%d) should be different than null", 0)
	}
	if nodeData.NodeByIndex(1) == nil {
		t.Errorf("Node.NodeByIndex() ==> NodeByIndex(%d) should be different than null", 1)
	}
	if nodeData.NodeByIndex(2) == nil {
		t.Errorf("Node.NodeByIndex() ==> NodeByIndex(%d) should be different than null", 2)
	}
	if nodeData.NodeByIndex(3) == nil {
		t.Errorf("Node.NodeByIndex() ==> NodeByIndex(%d) should be different than null", 3)
	}
	if nodeData.NodeByIndex(4) != nil {
		t.Errorf("Node.NodeByIndex() ==> NodeByIndex(%d) should be equal to null", 4)
	}
}
