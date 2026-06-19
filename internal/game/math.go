package game

import "math"

type Vec2 struct {
	X float64
	Y float64
}

func (v Vec2) Add(o Vec2) Vec2 {
	return Vec2{X: v.X + o.X, Y: v.Y + o.Y}
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{X: v.X - o.X, Y: v.Y - o.Y}
}

func (v Vec2) Mul(s float64) Vec2 {
	return Vec2{X: v.X * s, Y: v.Y * s}
}

func (v Vec2) LenSq() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Len() float64 {
	return math.Sqrt(v.LenSq())
}

func (v Vec2) Normalized() Vec2 {
	length := v.Len()
	if length == 0 {
		return Vec2{}
	}
	return Vec2{X: v.X / length, Y: v.Y / length}
}

func DistanceSq(a, b Vec2) float64 {
	return a.Sub(b).LenSq()
}

func Clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
