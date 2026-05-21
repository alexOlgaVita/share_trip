package tracker

import (
	"fmt"
	"slices"
	"strings"
)

type Item struct {
	ID   string
	Name string
}

func (i Item) toString() string {

	return fmt.Sprintf("%s\t%s", i.ID, i.Name)
}

type Tracker struct {
	Items []Item
}

func NewTracker() *Tracker {
	return &Tracker{}
}

func (t *Tracker) AddItem(item Item) error {
	_, ok := t.indexOf(item.ID)
	if ok {
		return ErrAlreadyExists
	}
	t.Items = append(t.Items, item)
	return nil
}

func (t *Tracker) GetItems() []Item {
	return t.Items
}

func (t *Tracker) DeleteItem(name string) {
	for i := 0; i < len(t.Items); i++ {
		if t.Items[i].Name == name {
			res := t.Items[i].toString()
			t.Items = slices.Delete(t.Items, i, i+1)
			fmt.Printf("Item '%s' was deleted:\n", res)
			return
		}
	}
	fmt.Println("There is no item with this name")
}

func (t *Tracker) UpdateItem(item Item) error {
	index, ok := t.indexOf(item.ID)
	if !ok {
		return ErrNotFound
	}
	t.Items[index] = item
	return nil
}

func (t *Tracker) indexOf(id string) (int, bool) {
	for i, item := range t.Items {
		if item.ID == id {
			return i, true
		}
	}
	return -1, false
}

func (t *Tracker) FindItem(name string) {
	res := make([]string, 0)
	for i := 0; i < len(t.Items); i++ {
		if strings.Contains(t.Items[i].Name, name) {
			res = append(res, t.Items[i].toString())
		}
	}
	if len(res) == 0 {
		fmt.Printf("There is no item containing text %s:\n", name)
		return
	}
	fmt.Printf("These items containing text %s: were found:\n", name)
	fmt.Println(strings.Join(res, ",\n"))
}
