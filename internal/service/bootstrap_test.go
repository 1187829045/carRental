package service

import (
	"testing"
	"time"
)

func TestBuildMockOperationalSeeds(t *testing.T) {
	now := time.Date(2026, 5, 3, 12, 0, 0, 0, time.Local)
	seeds := buildMockOperationalSeeds(now, 3, 1)
	if len(seeds) != mockCheckSeedCount {
		t.Fatalf("expected %d seeds, got %d", mockCheckSeedCount, len(seeds))
	}

	seenChecks := make(map[string]struct{}, len(seeds))
	seenRents := make(map[string]struct{}, len(seeds))
	seenCars := make(map[string]struct{}, len(seeds))
	seenIdentities := make(map[string]struct{}, len(seeds))
	zhangsanCount := 0

	threeMonthsAgo := now.AddDate(0, 0, -89)
	for i, seed := range seeds {
		if _, ok := seenChecks[seed.Check.CheckID]; ok {
			t.Fatalf("duplicate check id: %s", seed.Check.CheckID)
		}
		seenChecks[seed.Check.CheckID] = struct{}{}
		if _, ok := seenRents[seed.Rent.RentID]; ok {
			t.Fatalf("duplicate rent id: %s", seed.Rent.RentID)
		}
		seenRents[seed.Rent.RentID] = struct{}{}
		if _, ok := seenCars[seed.Car.CarNumber]; ok {
			t.Fatalf("duplicate car number: %s", seed.Car.CarNumber)
		}
		seenCars[seed.Car.CarNumber] = struct{}{}
		if _, ok := seenIdentities[seed.Customer.Identity]; ok {
			t.Fatalf("duplicate identity: %s", seed.Customer.Identity)
		}
		seenIdentities[seed.Customer.Identity] = struct{}{}

		if seed.Check.CheckDate.Before(threeMonthsAgo) || seed.Check.CheckDate.After(now) {
			t.Fatalf("seed %d check date out of range: %s", i, seed.Check.CheckDate)
		}
		if seed.Rent.ReturnDate == nil {
			t.Fatalf("seed %d missing return date", i)
		}
		if seed.Rent.Identity != seed.Customer.Identity {
			t.Fatalf("seed %d customer/rent identity mismatch", i)
		}
		if seed.Check.RentID != seed.Rent.RentID {
			t.Fatalf("seed %d check/rent id mismatch", i)
		}
		if i < mockZhangsanSeedCount {
			if seed.Check.OperId != 3 || seed.Car.OperId != 3 || seed.LegacyOperName != "张三" {
				t.Fatalf("seed %d expected zhangsan ownership", i)
			}
			zhangsanCount++
		}
	}

	if zhangsanCount != mockZhangsanSeedCount {
		t.Fatalf("expected %d zhangsan seeds, got %d", mockZhangsanSeedCount, zhangsanCount)
	}
}
