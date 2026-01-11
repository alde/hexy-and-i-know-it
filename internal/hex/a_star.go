package hex

import (
	"container/heap"
	"slices"
)

type Hex struct {
	Q int64
	R int64
}

type PathNode struct {
	hex    Hex
	fScore int
	index  int
}

type PriorityQueue []*PathNode

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].fScore < pq[j].fScore
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(pnode any) {
	n := len(*pq)
	*pq = append(*pq, pnode.(*PathNode))
	(*pq)[n].index = n
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func HexDistance(a, b Hex) int64 {
	dq := a.Q - b.Q
	dr := a.R - b.R
	return max(abs(dq)+abs(dr)+abs(dq+dr)) / 2
}

func abs[T int64 | int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func GetNeighbors(hex Hex) []Hex {
	return []Hex{
		{hex.Q + 1, hex.R},
		{hex.Q - 1, hex.R},
		{hex.Q, hex.R + 1},
		{hex.Q, hex.R - 1},
		{hex.Q + 1, hex.R - 1},
		{hex.Q - 1, hex.R + 1},
	}
}

func FindPath(start, goal Hex, isWalkable func(Hex) bool) []Hex {
	openSet := &PriorityQueue{}
	closedSet := make(map[Hex]bool)
	cameFrom := make(map[Hex]Hex)
	gScore := make(map[Hex]int)
	heap.Init(openSet)

	heap.Push(openSet, &PathNode{hex: start, fScore: 0})

	for openSet.Len() > 0 {
		current := heap.Pop(openSet).(*PathNode).hex

		if current == goal {
			return reconstructPath(cameFrom, current)
		}

		closedSet[current] = true

		for _, neighbor := range GetNeighbors(current) {
			ifWalkable := isWalkable(neighbor)
			_, inClosedSet := closedSet[neighbor]

			if !ifWalkable || inClosedSet {
				continue
			}

			tentativeGScore := gScore[current] + 1
			tentative, exists := gScore[neighbor]
			if !exists || tentativeGScore < tentative {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeGScore
				fScore := tentativeGScore + int(HexDistance(neighbor, goal))
				heap.Push(openSet, &PathNode{hex: neighbor, fScore: fScore})
			}
		}
	}
	return nil
}

func reconstructPath(cameFrom map[Hex]Hex, current Hex) []Hex {
	path := []Hex{current}
	for {
		prev, exists := cameFrom[current]
		if !exists {
			break
		}
		path = append(path, prev)
		current = prev
	}
	slices.Reverse(path)
	return path
}
