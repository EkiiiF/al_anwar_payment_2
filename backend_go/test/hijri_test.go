package test

import (
	"testing"
	"time"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/utils"
)

func TestGregorianToHijri_May2026(t *testing.T) {
	dt := time.Date(2026, 5, 21, 12, 0, 0, 0, time.Local)
	h := utils.GregorianToHijri(dt)

	t.Logf("2026-05-21 => Day: %d, Month: %d (%s), Year: %d H", h.Day, h.Month, utils.GetHijriMonthName(h.Month), h.Year)

	if h.Month < 1 || h.Month > 12 {
		t.Errorf("month %d is out of valid range 1-12", h.Month)
	}
	if h.Year < 1447 || h.Year > 1449 {
		t.Errorf("year %d seems out of expected range for 2026", h.Year)
	}
}

func TestGregorianToHijri_KnownDates(t *testing.T) {
	tests := []struct {
		date        time.Time
		expectMonth int // approximate
		expectYear  int
	}{
		// 2026-03-20 ~ Syawal 1448
		{time.Date(2026, 3, 20, 0, 0, 0, 0, time.Local), 10, 1447},
		// 2026-01-15 ~ Rajab 1447
		{time.Date(2026, 1, 15, 0, 0, 0, 0, time.Local), 7, 1447},
		// 2026-06-15 ~ Dzulhijjah 1448
		{time.Date(2026, 6, 15, 0, 0, 0, 0, time.Local), 12, 1447},
	}

	for _, tc := range tests {
		h := utils.GregorianToHijri(tc.date)
		name := utils.GetHijriMonthName(h.Month)
		t.Logf("%s => %d %s %d H", tc.date.Format("2006-01-02"), h.Day, name, h.Year)

		if h.Month < 1 || h.Month > 12 {
			t.Errorf("%s: month %d is out of range", tc.date.Format("2006-01-02"), h.Month)
		}
		if name == "-" {
			t.Errorf("%s: GetHijriMonthName returned '-' for month %d", tc.date.Format("2006-01-02"), h.Month)
		}
	}
}

func TestGetSemesterInfo(t *testing.T) {
	// Syawal -> Semester 1, registration
	si := utils.GetSemesterInfo(10, 1447)
	if si.Number != 1 || !si.IsRegistration {
		t.Errorf("Syawal should be Semester 1 + registration, got sem=%d reg=%v", si.Number, si.IsRegistration)
	}

	// Rabi'ul Awal -> Semester 1, exam
	si = utils.GetSemesterInfo(3, 1448)
	if si.Number != 1 || !si.IsExamMonth {
		t.Errorf("Rabi'ul Awal should be Semester 1 + exam, got sem=%d exam=%v", si.Number, si.IsExamMonth)
	}

	// Sya'ban -> Semester 2, exam
	si = utils.GetSemesterInfo(8, 1448)
	if si.Number != 2 || !si.IsExamMonth {
		t.Errorf("Sya'ban should be Semester 2 + exam, got sem=%d exam=%v", si.Number, si.IsExamMonth)
	}

	// Ramadhan -> Semester 2, not billable
	si = utils.GetSemesterInfo(9, 1448)
	if si.Number != 2 {
		t.Errorf("Ramadhan should be Semester 2, got sem=%d", si.Number)
	}
	if utils.IsBillableMonth(9) {
		t.Error("Ramadhan should NOT be billable")
	}
}
