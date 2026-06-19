package game

import "math"

type gridCell struct {
	X int
	Y int
}

type SpatialIndex struct {
	cellSize float64
	buckets  map[gridCell][]int
}

func NewSpatialIndex(cellSize float64) SpatialIndex {
	return SpatialIndex{
		cellSize: cellSize,
		buckets:  map[gridCell][]int{},
	}
}

func (s *SpatialIndex) Rebuild(skeletons []Skeleton) {
	clear(s.buckets)
	for i := range skeletons {
		cell := s.cellFor(skeletons[i].Pos)
		s.buckets[cell] = append(s.buckets[cell], i)
	}
}

func (s SpatialIndex) FirstNear(pos Vec2, radius float64, skeletons []Skeleton, matches func(int) bool) int {
	radiusSq := radius * radius
	found := -1
	s.ForEachNear(pos, radius, skeletons, func(i int) bool {
		if DistanceSq(pos, skeletons[i].Pos) <= radiusSq && matches(i) {
			found = i
			return false
		}
		return true
	})
	return found
}

func (s SpatialIndex) ForEachNear(pos Vec2, radius float64, skeletons []Skeleton, body func(int) bool) {
	minCell := s.cellFor(Vec2{X: pos.X - radius, Y: pos.Y - radius})
	maxCell := s.cellFor(Vec2{X: pos.X + radius, Y: pos.Y + radius})
	for y := minCell.Y; y <= maxCell.Y; y++ {
		for x := minCell.X; x <= maxCell.X; x++ {
			for _, idx := range s.buckets[gridCell{X: x, Y: y}] {
				if !body(idx) {
					return
				}
			}
		}
	}
}

func (s SpatialIndex) ForEachRect(min, max Vec2, body func(int) bool) {
	minCell := s.cellFor(min)
	maxCell := s.cellFor(max)
	for y := minCell.Y; y <= maxCell.Y; y++ {
		for x := minCell.X; x <= maxCell.X; x++ {
			for _, idx := range s.buckets[gridCell{X: x, Y: y}] {
				if !body(idx) {
					return
				}
			}
		}
	}
}

func (s SpatialIndex) cellFor(pos Vec2) gridCell {
	return gridCell{
		X: int(math.Floor(pos.X / s.cellSize)),
		Y: int(math.Floor(pos.Y / s.cellSize)),
	}
}
