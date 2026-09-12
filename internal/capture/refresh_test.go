package capture

import (
	"testing"
	"time"
)

func TestRefreshAfterForegroundSettles(t *testing.T) {
	var plan RefreshPlan
	now := time.Unix(100, 0)
	calls := 0
	attach := func() bool {
		calls++
		return true
	}
	plan.Tick(10, now, true, attach)
	plan.Tick(10, now.Add(100*time.Millisecond), true, attach)
	if calls != 0 {
		t.Fatal("refreshed before foreground settled")
	}
	plan.Tick(10, now.Add(time.Second), true, attach)
	plan.Tick(10, now.Add(2*time.Second), true, attach)
	if calls != 1 {
		t.Fatalf("refresh must run once per foreground change, got %d", calls)
	}
	plan.Tick(20, now.Add(3*time.Second), true, attach)
	plan.Tick(20, now.Add(4*time.Second), true, attach)
	if calls != 2 {
		t.Fatal("did not refresh after entering another window")
	}
}

func TestRefreshWaitsForReleasedKeysAndEnabledMapping(t *testing.T) {
	var plan RefreshPlan
	now := time.Unix(100, 0)
	calls := 0
	attach := func() bool {
		calls++
		return true
	}
	plan.Tick(10, now, false, attach)
	plan.Tick(10, now.Add(time.Second), false, attach)
	if calls != 0 {
		t.Fatal("refreshed while keys were held or mapping was paused")
	}
	plan.Tick(10, now.Add(2*time.Second), true, attach)
	if calls != 1 {
		t.Fatal("pending refresh was lost")
	}
}

func TestManualRefreshAllowsTimeToReturnToTargetWindow(t *testing.T) {
	var plan RefreshPlan
	now := time.Unix(100, 0)
	calls := 0
	attach := func() bool {
		calls++
		return true
	}
	plan.Request(now.Add(2 * time.Second))
	plan.Tick(10, now, true, attach)
	plan.Tick(20, now.Add(time.Second), true, attach)
	if calls != 0 {
		t.Fatal("foreground change bypassed manual delay")
	}
	plan.Tick(20, now.Add(3*time.Second), true, attach)
	if calls != 1 {
		t.Fatal("manual request did not run")
	}
}

func TestRefreshFailureDoesNotBusyLoop(t *testing.T) {
	var plan RefreshPlan
	now := time.Unix(100, 0)
	calls := 0
	attach := func() bool {
		calls++
		return false
	}
	plan.Tick(10, now, true, attach)
	plan.Tick(10, now.Add(time.Second), true, attach)
	plan.Tick(10, now.Add(2*time.Second), true, attach)
	if calls != 1 {
		t.Fatal("failed installation retried without a new request")
	}
	plan.Request(now.Add(3 * time.Second))
	plan.Tick(10, now.Add(4*time.Second), true, attach)
	if calls != 2 {
		t.Fatal("manual retry did not run")
	}
}

func TestNullForegroundDoesNotRefresh(t *testing.T) {
	var plan RefreshPlan
	now := time.Unix(100, 0)
	plan.Tick(0, now, true, func() bool {
		t.Fatal("null foreground must not install hooks")
		return false
	})
}
