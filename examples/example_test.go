package examples

import (
	"fmt"
	"time"

	"github.com/nekrassov01/filter"
)

func Example_basic() {
	stats := &Stats{
		Class:      "軍師",
		Name:       "諸葛亮 孔明",
		Birth:      time.Date(181, 7, 23, 0, 0, 0, 0, time.UTC),
		ATBGauge:   time.Second * 30,
		HitPoint:   80,
		SkillPoint: 0,
		SpellPoint: 250,
		LifePoint:  5,
		Strength:   10,
		Stamina:    10,
		Dexterity:  10,
		Magic:      25,
		Speed:      25,
	}
	inputs := []string{
		`Class == "軍師"`,
		`Name =~ '^(諸葛亮|龐統|法正)'`,
		`Name != ""`,
		`BirthDate < '0190-01-01T00:00:00Z'`,
		`ActiveTimeBattleGauge >= '20s'`,
		`HitPoint > "50"`,
		`MagicPoint > 100`,
		`LifePoint != 0`,
		`Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != ""`,
		`BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s'`,
		`HitPoint > "50" && MagicPoint > 100 && LifePoint != 0`,
		`Magic >= 20 || !(Speed < 20)`,
		`Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s') && (HitPoint > "50" && MagicPoint > 100 && LifePoint != 0) && (Magic >= 20 || !(Speed < 20))`,
		`Class == "君主"`,
	}
	for _, input := range inputs {
		expr, err := filter.Parse(input)
		if err != nil {
			fmt.Printf("%s: %v\n", input, err)
			continue
		}
		ok, err := expr.Eval(stats)
		if err != nil {
			fmt.Printf("%s: %v\n", input, err)
			continue
		}
		fmt.Printf("%s: %v\n", input, ok)
	}
	// Output:
	// Class == "軍師": true
	// Name =~ '^(諸葛亮|龐統|法正)': true
	// Name != "": true
	// BirthDate < '0190-01-01T00:00:00Z': true
	// ActiveTimeBattleGauge >= '20s': true
	// HitPoint > "50": true
	// MagicPoint > 100: true
	// LifePoint != 0: true
	// Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "": true
	// BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s': true
	// HitPoint > "50" && MagicPoint > 100 && LifePoint != 0: true
	// Magic >= 20 || !(Speed < 20): true
	// Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s') && (HitPoint > "50" && MagicPoint > 100 && LifePoint != 0) && (Magic >= 20 || !(Speed < 20)): true
	// Class == "君主": false
}

func Example_timeLayouts() {
	event := &Event{
		At: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
	}
	inputs := []string{
		// Forms without spaces can be written bare.
		`At == 2025-01-01T09:00:00Z`,
		`At == 2025-01-01T18:00:00+09:00`,
		`At == 2025-01-01T09:00:00`,
		`At >= 2025-01-01`,
		`At == 1735722000`,
		// Forms with spaces must be quoted.
		`At == '2025-01-01 09:00:00'`,
		`At == 'Wed, 01 Jan 2025 09:00:00 UTC'`,
		`At == 'Wed, 01 Jan 2025 18:00:00 +0900'`,
		`At == 'Wednesday, 01-Jan-25 09:00:00 GMT'`,
		`At == '01 Jan 25 09:00 UTC'`,
		`At == '01 Jan 25 18:00 +0900'`,
		// A zone abbreviation other than UTC or GMT is rejected.
		`At == 'Wed, 01 Jan 2025 04:00:00 EST'`,
	}
	for _, input := range inputs {
		expr, err := filter.Parse(input)
		if err != nil {
			fmt.Printf("%s: %v\n", input, err)
			continue
		}
		ok, err := expr.Eval(event)
		if err != nil {
			fmt.Printf("%s: %v\n", input, err)
			continue
		}
		fmt.Printf("%s: %v\n", input, ok)
	}
	// Output:
	// At == 2025-01-01T09:00:00Z: true
	// At == 2025-01-01T18:00:00+09:00: true
	// At == 2025-01-01T09:00:00: true
	// At >= 2025-01-01: true
	// At == 1735722000: true
	// At == '2025-01-01 09:00:00': true
	// At == 'Wed, 01 Jan 2025 09:00:00 UTC': true
	// At == 'Wed, 01 Jan 2025 18:00:00 +0900': true
	// At == 'Wednesday, 01-Jan-25 09:00:00 GMT': true
	// At == '01 Jan 25 09:00 UTC': true
	// At == '01 Jan 25 18:00 +0900': true
	// At == 'Wed, 01 Jan 2025 04:00:00 EST': eval error at 1:7: invalid time "Wed, 01 Jan 2025 04:00:00 EST"
}
