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
	CacheKey Key
	Value    interface{}
	Next     *ListItem
	Prev     *ListItem
}

type list struct {
	len   int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushAfter(prev *ListItem, item *ListItem) {
	item.Prev = prev

	if prev.Next == nil {
		item.Next = nil
		l.back = item
	} else {
		item.Next = prev.Next
		prev.Next.Prev = item
	}

	prev.Next = item
	l.len++
}

func (l *list) PushBefore(next *ListItem, item *ListItem) {
	item.Next = next

	if next.Prev == nil {
		item.Prev = nil
		l.front = item
	} else {
		item.Prev = next.Prev
		next.Prev.Next = item
	}

	next.Prev = item
	l.len++
}

func (l *list) PushItemFront(i *ListItem) {
	if l.front == nil {
		l.front = i
		l.back = i
		i.Prev, i.Next = nil, nil
		l.len++

		return
	}

	l.PushBefore(l.front, i)
}

func (l *list) PushFront(v any) *ListItem {
	item := &ListItem{Value: v}
	l.PushItemFront(item)

	return item
}

func (l *list) PushBack(v any) *ListItem {
	item := &ListItem{Value: v}

	if l.back == nil {
		l.PushItemFront(item)
	} else {
		l.PushAfter(l.back, item)
	}

	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	l.Remove(i)
	l.PushItemFront(i)
}
