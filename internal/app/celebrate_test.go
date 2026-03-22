package app

import "testing"

func TestSpawnPerfectCelebration(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	if len(c.particles) != 12 {
		t.Errorf("perfect tier particles = %d, want 12", len(c.particles))
	}
	if c.nextState != stateDrillComplete {
		t.Error("nextState should be stateDrillComplete")
	}
	if !c.active() {
		t.Error("celebration should be active after spawn")
	}
}

func TestSpawnLevelUpCelebration(t *testing.T) {
	c := spawnCelebration(tierLevelUp, stateLevelComplete, "LEVEL 3 COMPLETE!")
	if len(c.particles) != 20 {
		t.Errorf("level-up tier particles = %d, want 20", len(c.particles))
	}
}

func TestTickReducesLifetime(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	initialLife := c.particles[0].lifetime
	c.tick()
	if c.particles[0].lifetime != initialLife-1 {
		t.Errorf("lifetime after tick = %d, want %d", c.particles[0].lifetime, initialLife-1)
	}
}

func TestTickRemovesDeadParticles(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	for i := range c.particles {
		c.particles[i].lifetime = 1
	}
	c.tick()
	if len(c.particles) != 0 {
		t.Errorf("particles after expiry = %d, want 0", len(c.particles))
	}
	if c.active() {
		t.Error("celebration should be inactive after all particles expire")
	}
}

func TestCelebrationRender(t *testing.T) {
	c := spawnCelebration(tierPerfect, stateDrillComplete, "PERFECT!")
	result := c.render(80)
	if result == "" {
		t.Error("render should produce output")
	}
}

func TestPassTierNoParticles(t *testing.T) {
	c := spawnCelebration(tierPass, statePhaseComplete, "Phase complete!")
	if len(c.particles) != 0 {
		t.Errorf("pass tier particles = %d, want 0", len(c.particles))
	}
}
