package game

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExperienceLogRecordsAwardedXP(t *testing.T) {
	path := filepath.Join(t.TempDir(), "xp.log")
	g := NewWithExperienceLog(path)
	g.skeleton = []Skeleton{{ID: 42, Kind: SkeletonRed, HP: 1, Reward: SkeletonRed.ExperienceReward()}}

	g.destroySkeleton(0, AttackFireball)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read xp log: %v", err)
	}
	line := string(data)
	for _, want := range []string{
		"attack=fireball",
		"deaths=1",
		"xp=1",
		"enemy=red",
		"enemy_id=42",
		"level=1->2",
		"xp_bar=0/1->0/2",
		"level_ups=1",
		"total_kills=1",
	} {
		if !strings.Contains(line, want) {
			t.Fatalf("xp log %q missing %q", line, want)
		}
	}
}
