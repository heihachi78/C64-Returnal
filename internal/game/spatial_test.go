package game

import "testing"

func TestSpatialIndexFindsNearbyCandidate(t *testing.T) {
	index := NewSpatialIndex(32)
	skeletons := []Skeleton{
		{ID: 1, Pos: Vec2{X: 100, Y: 100}},
		{ID: 2, Pos: Vec2{X: 8, Y: 10}},
	}
	index.Rebuild(skeletons)

	got := index.FirstNear(Vec2{X: 0, Y: 0}, 16, skeletons, func(i int) bool {
		return skeletons[i].ID == 2
	})
	if got != 1 {
		t.Fatalf("FirstNear = %d, want 1", got)
	}
}

func TestSpatialIndexRectIteration(t *testing.T) {
	index := NewSpatialIndex(32)
	skeletons := []Skeleton{
		{ID: 1, Pos: Vec2{X: -10, Y: -10}},
		{ID: 2, Pos: Vec2{X: 96, Y: 96}},
		{ID: 3, Pos: Vec2{X: 4, Y: 7}},
	}
	index.Rebuild(skeletons)

	seen := map[int]bool{}
	index.ForEachRect(Vec2{X: -16, Y: -16}, Vec2{X: 16, Y: 16}, func(i int) bool {
		seen[skeletons[i].ID] = true
		return true
	})

	if !seen[1] || !seen[3] {
		t.Fatalf("rect did not include expected skeletons: %#v", seen)
	}
	if seen[2] {
		t.Fatalf("rect included far skeleton: %#v", seen)
	}
}
