package game

type Player struct {
	Pos           Vec2
	Facing        float64
	Moving        bool
	MoveDir       Vec2
	AnimTimer     float64
	AnimFrame     int
	HitFlash      float64
	DeathTimer    float64
	DeathRotation float64
}
type Skeleton struct {
	ID        int
	Pos       Vec2
	Kind      SkeletonKind
	HP        int
	Reward    int
	Facing    float64
	HitFlash  float64
	AnimFrame int
}
type Fireball struct {
	Pos               Vec2
	TargetID          int
	Velocity          Vec2
	TimeWithoutTarget float64
	AnimFrame         int
}
type OrbitalOrb struct {
	Pos                  Vec2
	Active               bool
	MissingOrbitProgress float64
	AnimFrame            int
}
type MeteorProjectile struct {
	Pos       Vec2
	Start     Vec2
	Impact    Vec2
	Age       float64
	AnimFrame int
}
type Coin struct {
	Pos    Vec2
	Amount int
	Level  int
	Phase  float64
}
type Chest struct {
	Pos  Vec2
	Tier ChestTier
}
type Effect struct {
	Kind        EffectKind
	Start       Vec2
	End         Vec2
	Pos         Vec2
	Points      []Vec2
	InnerPoints []Vec2
	Frame       int
	Facing      float64
	Radius      float64
	TTL         float64
	MaxTTL      float64
}
