package queue

import "sync"

// bh represents a binary heap [https://en.wikipedia.org/wiki/Binary_heap]
type bh struct {
	m sync.RWMutex

	// items
	items []*Item

	// hC signals on head change
	hC chan struct{}
}

// BinaryHeap returns an instance of the binary heap
func BinaryHeap() *bh {
	return &bh{
		hC: make(chan struct{}, 1),
	}
}

// satisfyMinInvariant compares parent and child to verify min heap invariant
func satisfyMinInvariant(parent, child *Item) bool {
	return parent.When.Before(child.When)
}

// signal tries to send a dummy struct to the hC chan or drains it
func (h *bh) signal() {
	select {
	case h.hC <- struct{}{}:
	case <-h.hC:
	}
}

// bubbleUp
func (h *bh) bubbleUp() {
	for pos := len(h.items) - 1; pos > 0; {
		// parent position
		ppos := (pos - 1) / 2

		// compare current and parent
		if satisfyMinInvariant(h.items[ppos], h.items[pos]) {
			break
		}

		// swap current <=> parent
		h.items[pos], h.items[ppos] = h.items[ppos], h.items[pos]

		// update position
		pos = ppos
	}
}

// removeItem
func (h *bh) removeItem(pos int) {
	lastIndex := len(h.items) - 1

	// replace with bottom
	h.items[pos] = h.items[lastIndex]

	// remove last item
	h.items = h.items[:lastIndex]

	// bubble down
	for {
		l, r := (2*pos)+1, (2*pos)+2

		if l >= len(h.items) {
			break
		}

		// smallest child node
		child := l
		if r < len(h.items) && h.items[l].When.After(h.items[r].When) {
			child = r
		}

		// verify heap invariant
		if satisfyMinInvariant(h.items[pos], h.items[child]) {
			break
		}

		// swap current <=> parent
		h.items[pos], h.items[child] = h.items[child], h.items[pos]

		// update position
		pos = child
	}
}

// Size returns the amount of items enqueued
func (h *bh) Size() int {
	return len(h.items)
}

// Add enqueues a item
func (h *bh) Add(item *Item) {
	h.m.Lock()
	defer h.m.Unlock()

	h.items = append(h.items, item)
	oldHead := h.items[0]

	h.bubbleUp()

	// signal on head change
	if oldHead.Id != h.items[0].Id {
		h.signal()
	}
}

// Pop returns the next enqueued item
func (h *bh) Pop() *Item {
	h.m.Lock()
	defer h.m.Unlock()

	// no items
	if len(h.items) < 1 {
		return nil
	}

	head := h.items[0]

	h.removeItem(0)

	return head
}

// Remove a queued item with matching id
func (h *bh) Remove(id int64) {
	h.m.Lock()
	defer h.m.Unlock()

	pos := -1

	for i, item := range h.items {
		if id == item.Id {
			pos = i
			break
		}
	}

	// no items or not found
	if pos == -1 {
		return
	}

	// head will change
	headChange := pos == 0

	h.removeItem(pos)

	// signal on head change
	if headChange {
		h.signal()
	}
}

// Peek returns the next enqueued item without removing it from queue
func (h *bh) Peek() *Item {
	h.m.RLock()
	defer h.m.RUnlock()

	// no items
	if len(h.items) < 1 {
		return nil
	}

	return h.items[0]
}

// OnHeadChange returns a channel to listen to queue head change events
func (h *bh) OnHeadChange() <-chan struct{} {
	return h.hC
}
