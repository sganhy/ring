package data

import (
	"fmt"
	"strconv"
	"strings"
)

type bitchain struct {
	data uint64
	next *bitchain
}

const (
	nodeAllDataTo1 uint64 = 0xFFFFFFFFFFFFFFFF
	m1                    = 0x5555555555555555 //binary: 0101...
	m2                    = 0x3333333333333333 //binary: 00110011..
	m4                    = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...
	h01                   = 0x0101010101010101 //the sum of 256 to the power of 0,1,2,3...
)

var (
	displayNodeSeparator = "{%d} - "
	displayNodeNewLine   = "\n"
)

//******************************
// getters and setters
//******************************
func (node *bitchain) getData() uint64 {
	return node.data
}

//******************************
// public methods
//******************************
func (node *bitchain) String() string {
	var sb strings.Builder
	var nodeTemp = node
	var index = 0

	for ok := true; ok; ok = nodeTemp != nil {
		sb.WriteString(displayNodeNewLine)
		sb.WriteString(fmt.Sprintf(displayNodeSeparator, index))
		sb.WriteString(strconv.FormatInt(int64(nodeTemp.data), 2))
		nodeTemp = nodeTemp.next
		index++
	}

	return sb.String()
}
func (node *bitchain) Count() int {
	var result = 1
	var currentNode = node
	for currentNode.next != nil {
		currentNode = currentNode.next
		result++
	}
	return result
}
func (node *bitchain) NodeByIndex(index int) *bitchain {
	nodeTemp := node
	i := 0
	for ok := true; ok; ok = nodeTemp != nil {
		if i == index {
			return nodeTemp
		}
		nodeTemp = nodeTemp.next
		i++
	}
	return nil
}
func (node *bitchain) SetValue(bitPosition uint8, value bool) {
	currentNode := node
	for position := int(bitPosition); position >= 0; position -= 64 {
		if position < 64 {
			var mask uint64 = 1
			mask <<= position
			if value == true {
				currentNode.data |= mask
			} else {
				currentNode.data &= ^mask
			}
			break
		} else {
			if currentNode.next == nil {
				// no need to allow here (by default == false)
				if value == false {
					return
				}
				currentNode.next = new(bitchain)
			}
			currentNode = currentNode.next
		}
	}
}
func (node *bitchain) GetValue(bitPosition uint8) bool {
	currentNode := node
	for position := int(bitPosition); position >= 0; position -= 64 {
		if position < 64 {
			return ((currentNode.data >> position) & 1) > 0
		}
		if currentNode.next == nil {
			break
		}
		currentNode = currentNode.next
	}
	return false
}
func (node *bitchain) ResetAll(bitPosition uint8, value bool) {
	if value {
		currentNode := node
		position := int(bitPosition)
		for {
			currentNode.data = nodeAllDataTo1
			position -= 64
			if position < 0 {
				break
			}
			if currentNode.next == nil {
				currentNode.next = new(bitchain)
			}
			currentNode = currentNode.next
		}
	} else {
		node.clearAll(bitPosition)
	}
}

// CountSetBits Hamming weight approach: complexity O(1)
func (node *bitchain) CountSetBits() int {
	currentNode := node
	result := 0
	var x uint64 = 0
	for currentNode != nil {
		x = currentNode.data
		x -= (x >> 1) & m1
		x = (x & m2) + ((x >> 2) & m2)
		x = (x + (x >> 4)) & m4
		result += int((x * h01) >> 56)
		currentNode = currentNode.next
	}
	return result
}

//******************************
// private methods
//******************************
// bitPosition [0, 63]
func (node *bitchain) clearAll(bitPosition uint8) {
	currentNode := node
	position := int(bitPosition)
	for {
		currentNode.data = 0
		position -= 64
		if position < 0 || currentNode.next == nil {
			break
		}
		currentNode = currentNode.next
	}
}
