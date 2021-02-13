package data

type List struct {
	sorted   bool
	data     []*Record
	itemType string
}

//******************************
// public methods
//******************************
func (list *List) Count() int {
	if list.data == nil {
		return 0
	}
	return len(list.data)
}

func (list *List) ItemType() string {
	return list.itemType
}

func (list *List) Sorted() bool {
	return list.sorted
}

func (list *List) ItemByIndex(index int) *Record {
	if list.data != nil && index >= 0 && index < len(list.data) {
		return list.data[index]
	}
	return nil
}

//func (list *List) AppendItem(item *Record) {
func (list *List) AppendItem(item interface{}) {
	if list.data == nil {
		list.data = make([]*Record, 0, 2)
	}
	switch v := item.(type) {
	case int:
		break
	case *Record:
		list.appendItem(v)
	case Record:
		list.appendItem(&v)
	default:
	}
}

func (list *List) appendItem(item *Record) {
	if list.data == nil {
		list.data = make([]*Record, 0, 2)
	}
	list.data = append(list.data, item)
}
