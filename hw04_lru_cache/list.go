package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	firstItem *ListItem
	lastItem  *ListItem
	len       int
}

func (l *list) Len() int {
	if l == nil {
		return 0
	}

	return l.len
}

func (l *list) Front() *ListItem {
	if l == nil {
		return nil
	}

	return l.firstItem
}

func (l *list) Back() *ListItem {
	if l == nil {
		return nil
	}

	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	listItem := ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.len == 0 {
		l.len++
		l.firstItem = &listItem
		l.lastItem = &listItem

		return &listItem
	}

	listItem.Next = l.firstItem
	l.firstItem.Prev = &listItem
	l.firstItem = &listItem

	l.len++

	return &listItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	listItem := ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.len == 0 {
		l.len++
		l.firstItem = &listItem
		l.lastItem = &listItem

		return &listItem
	}

	l.len++

	listItem.Prev = l.lastItem
	l.lastItem.Next = &listItem
	l.lastItem = &listItem

	return &listItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	if i == l.firstItem {
		l.firstItem = i.Next
	}

	if i == l.lastItem {
		l.lastItem = i.Prev
	}

	l.len--

	// для сборщика мусора
	i.Next = nil
	i.Prev = nil
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.firstItem {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.lastItem = i.Prev
	}

	i.Next = l.firstItem
	i.Prev = nil
	l.firstItem.Prev = i
	l.firstItem = i
}

func NewList() List {
	return new(list)
}
