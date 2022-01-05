package data

import (
	"fmt"
	"strconv"
	"strings"
)

type node struct {
	data uint64
	next *node
}

var (
	displayNodeSeparator = "{%d} - "
	displayNodeNewLine   = "\n"
)

//******************************
// getters and setters
//******************************
func (nodeInfo *node) getData() uint64 {
	return nodeInfo.data
}

//******************************
// public methods
//******************************
func (nodeInfo *node) String() string {
	var sb strings.Builder
	var nodeTemp = nodeInfo
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
func (nodeInfo *node) Count() int {
	var result = 1
	var currentNode = nodeInfo
	for currentNode.next != nil {
		currentNode = currentNode.next
		result++
	}
	return result
}
func (nodeInfo *node) NodeByIndex(index int) *node {
	var currentNode = nodeInfo
	for ok := true; ok; ok = currentNode.next != nil {
		if index == 0 {
			return currentNode
		}
		currentNode = currentNode.next
		index--
	}
	return nil
}

func (nodeInfo *node) Dispose() *node {

}

//******************************
// private methods
//******************************
// bitPosition [0, 63]
func (nodeInfo *node) setValue(bitPosition uint8, value bool) {
	if bitPosition < 64 {
		nodeInfo.writeData(bitPosition, value)
	} else {
		if nodeInfo.next == nil {
			// no need to allow here (by default == false)
			if value == false {
				return
			}
			nodeInfo.next = new(node)
		}
		// recursive call!!
		nodeInfo.next.setValue(bitPosition-64, value)
	}
}
func (nodeInfo *node) getValue(bitPosition uint8) bool {
	if bitPosition < 64 {
		return nodeInfo.readData(bitPosition)
	} else {
		if nodeInfo.next == nil {
			return false
		} else {
			// recursive call!!
			return nodeInfo.next.getValue(bitPosition - 64)
		}
	}
}

func (nodeInfo *node) writeData(bitPosition uint8, value bool) {
	var mask uint64 = 0
	if bitPosition < 64 {
		mask = 1
		mask <<= bitPosition
		if value == true {
			nodeInfo.data |= mask
		} else {
			nodeInfo.data &= ^mask
		}
	}
}

func (nodeInfo *node) readData(bitPosition uint8) bool {
	return ((nodeInfo.data >> bitPosition) & 1) > 0
}
