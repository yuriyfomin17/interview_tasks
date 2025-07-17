package lru_cache

type Node struct {
	prev  *Node
	next  *Node
	key   int
	value int
}
type DoubleLinkedList struct {
	head *Node
	tail *Node
}

func (dll *DoubleLinkedList) moveToFront(currNode *Node) {
	if currNode == dll.head {
		return
	}
	if dll.head == nil {
		dll.head = currNode
		dll.tail = currNode
		return
	}
	prevNode := currNode.prev
	nextNode := currNode.next
	currNode.prev = nil
	currNode.next = nil

	if prevNode != nil {
		prevNode.next = nextNode
	}
	if nextNode != nil {
		nextNode.prev = prevNode
	}
	if dll.tail == currNode {
		dll.tail = prevNode

	}
	if dll.tail != nil {
		dll.tail.next = nil
	}
	currNode.next = dll.head
	dll.head.prev = currNode
	currNode.prev = nil
	dll.head = currNode

}
func (dll *DoubleLinkedList) removeElement() *Node {
	if dll.tail == nil {
		return nil
	}
	if dll.tail == dll.head {
		resp := dll.tail
		dll.tail = nil
		dll.head = nil
		return resp
	}
	resp := dll.tail
	dll.tail = dll.tail.prev
	dll.tail.next = nil
	return resp
}

type LRUCache struct {
	elements  map[int]*Node
	capacity  int
	cacheList *DoubleLinkedList
}

// cap 2  -> put (1, 1)        list (1, 1)
// cap 2  -> put (2, 2)        list (2, 2) (1, 1)
// cap 2. -> get (1) -> 1      list (1, 1) (2, 2)
// cap 2. -> put (3, 3)        list (3, 3) (1, 1)
// cap 2. -> get (2) -> -1     list (3, 3) (1, 1)
// cap 2. -> put (4, 4)        list (4, 4) (3, 3)
// cap 2. -> get (1) -> -1     list (4, 4) (3, 3)
// cap 2. -> get (3) -> 3      list (3, 3),(4, 4)
// cap 2. -> get (4) -> 4      list (4, 4), (3, 3)

func Constructor(capacity int) LRUCache {
	return LRUCache{
		elements:  make(map[int]*Node),
		capacity:  capacity,
		cacheList: &DoubleLinkedList{},
	}
}

func (this *LRUCache) Get(key int) int {

	if valNode, ok := this.elements[key]; ok {
		this.cacheList.moveToFront(valNode)
		return valNode.value
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if valNode, ok := this.elements[key]; ok {
		this.cacheList.moveToFront(valNode)
		valNode.value = value
		return
	}
	if len(this.elements) >= this.capacity {

		lastNode := this.cacheList.removeElement()
		delete(this.elements, lastNode.key)
	}
	newNode := &Node{
		key:   key,
		value: value,
	}
	this.elements[key] = newNode
	this.cacheList.moveToFront(newNode)
}
