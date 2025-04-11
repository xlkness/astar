package astar

import "container/heap"

type minFHeap []*calcCell

type sortedList struct {
	minHeap *minFHeap
}

func newSortedList() *sortedList {
	l := new(sortedList)
	minHeap := make(minFHeap, 0)
	l.minHeap = &minHeap
	heap.Init(l.minHeap)
	return l
}

func (list *sortedList) addCell(n *calcCell) {
	heap.Push(list.minHeap, n)
}

func (list *sortedList) hasCell(n *calcCell) bool {
	return n.sortIndex >= 0
}

func (list *sortedList) delCell(n *calcCell) {
	if n.sortIndex < 0 {
		return
	}
	heap.Remove(list.minHeap, n.sortIndex)
}

func (list *sortedList) getMinFCell() *calcCell {
	n := heap.Pop(list.minHeap)
	if n != nil {
		return n.(*calcCell)
	}
	return nil
}

func (list *sortedList) isEmpty() bool {
	return list.minHeap.Len() <= 0
}

func (mh minFHeap) Len() int {
	return len(mh)
}

func (mh minFHeap) Less(i, j int) bool {
	return mh[i].g+mh[i].h < mh[j].g+mh[j].h
}

func (mh minFHeap) Swap(i, j int) {
	mh[i], mh[j] = mh[j], mh[i]
	mh[i].sortIndex = i
	mh[j].sortIndex = j
}

func (mh *minFHeap) Push(x interface{}) {
	no := x.(*calcCell)
	no.sortIndex = len(*mh)
	*mh = append(*mh, no)
}

func (mh *minFHeap) Pop() interface{} {
	old := *mh
	n := len(old)
	no := old[n-1]
	no.sortIndex = -1
	*mh = old[0 : n-1]
	return no
}
