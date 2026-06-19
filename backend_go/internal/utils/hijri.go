package utils

import (
	"fmt"
	"math"
	"time"
)

type HijriDate struct {
	Day   int
	Month int
	Year  int
}

var HijriMonthNames = []string{
	"",
	"Muharram",
	"Safar",
	"Rabi'ul Awal",
	"Rabi'ul Akhir",
	"Jumadil Awal",
	"Jumadil Akhir",
	"Rajab",
	"Sya'ban",
	"Ramadhan",
	"Syawal",
	"Dzulqa'dah",
	"Dzulhijjah",
}

const (
	HijriSyawal      = 10
	HijriRobiulAwal  = 3
	HijriRobiulAkhir = 4
	HijriSyaban      = 8
	HijriRamadhan    = 9
)

type SemesterInfo struct {
	Number         int
	Name           string
	HijriMonth     int
	HijriYear      int
	IsExamMonth    bool
	IsRegistration bool
}

func gregorianToJD(y, m, d int) float64 {
	if m <= 2 {
		y--
		m += 12
	}
	a := math.Floor(float64(y) / 100.0)
	b := 2.0 - a + math.Floor(a/4.0)
	return math.Floor(365.25*float64(y+4716)) + math.Floor(30.6001*float64(m+1)) + float64(d) + b - 1524.5
}

func GregorianToHijri(t time.Time) HijriDate {
	jd := gregorianToJD(t.Year(), int(t.Month()), t.Day())
	jd = math.Floor(jd) + 0.5
	days := int(jd - 1948439.5)
	cycles := days / 10631
	remaining := days % 10631

	year := cycles * 30
	dayCount := 0
	for i := 1; i <= 30; i++ {
		daysInYear := 354
		if isHijriLeapYear(i) {
			daysInYear = 355
		}
		if dayCount+daysInYear > remaining {
			break
		}
		dayCount += daysInYear
		year++
	}
	remaining -= dayCount
	year++

	month := 0
	for m := 1; m <= 12; m++ {
		daysInMonth := 30
		if m%2 == 0 {
			daysInMonth = 29
		}
		if m == 12 && isHijriLeapYear(((year-1)%30)+1) {
			daysInMonth = 30
		}
		if remaining < daysInMonth {
			month = m
			break
		}
		remaining -= daysInMonth
	}

	day := remaining + 1
	if month < 1 {
		month = 1
	}
	if month > 12 {
		month = 12
	}

	return HijriDate{Day: day, Month: month, Year: year}
}

func isHijriLeapYear(yearInCycle int) bool {
	leapYears := []int{2, 5, 7, 10, 13, 16, 18, 21, 24, 26, 29}
	for _, ly := range leapYears {
		if yearInCycle == ly {
			return true
		}
	}
	return false
}

func GetHijriMonthName(month int) string {
	if month < 1 || month > 12 {
		return "-"
	}
	return HijriMonthNames[month]
}

func GetCurrentHijriDate() HijriDate {
	wib := time.FixedZone("WIB", 7*3600)
	return GregorianToHijri(time.Now().In(wib))
}
func GetSemesterInfo(hijriMonth, hijriYear int) SemesterInfo {
	info := SemesterInfo{
		HijriMonth: hijriMonth,
		HijriYear:  hijriYear,
	}

	switch {
	case hijriMonth >= HijriSyawal || hijriMonth <= HijriRobiulAwal:
		info.Number = 1
		info.Name = "Semester 1 (Syawal - Rabi'ul Awal)"
	case hijriMonth >= HijriRobiulAkhir && hijriMonth <= HijriSyaban:
		info.Number = 2
		info.Name = "Semester 2 (Rabi'ul Akhir - Sya'ban)"
	default:
		info.Number = 2
		info.Name = "Semester 2 (Rabi'ul Akhir - Sya'ban)"
	}
	info.IsExamMonth = (hijriMonth == HijriRobiulAwal || hijriMonth == HijriSyaban)
	info.IsRegistration = (hijriMonth == HijriSyawal)

	return info
}

func IsBillableMonth(hijriMonth int) bool {
	return true
}
func GetAcademicYearLabel(hijriMonth, hijriYear int) string {
	startYear := hijriYear
	if hijriMonth < HijriSyawal {
		startYear = hijriYear - 1
	}
	return fmt.Sprintf("%d/%d H", startYear, startYear+1)
}

func GetSemesterMonths(semester int) []int {
	if semester == 1 {
		return []int{10, 11, 12, 1, 2, 3} 
	}
	return []int{4, 5, 6, 7, 8} 
}

func HijriToGregorian(y, m, d int) time.Time {
	jd := float64(d) + math.Ceil(29.5*float64(m-1)) + float64(y-1)*354.0 + math.Floor(float64(11*y+3)/30.0) + 1948439.5

	jd = jd + 0.5
	z := math.Floor(jd)
	f := jd - z
	var a float64
	if z < 2299161 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1.0 + alpha - math.Floor(alpha/4.0)
	}
	b := a + 1524.0
	c := math.Floor((b - 122.1) / 365.25)
	dVal := math.Floor(365.25 * c)
	e := math.Floor((b - dVal) / 30.6001)

	day := b - dVal - math.Floor(30.6001*e) + f
	var month float64
	if e < 14 {
		month = e - 1
	} else {
		month = e - 13
	}
	var year float64
	if month > 2 {
		year = c - 4716
	} else {
		year = c - 4715
	}

	wib := time.FixedZone("WIB", 7*3600)
	return time.Date(int(year), time.Month(int(month)), int(day), 0, 0, 0, 0, wib)
}
