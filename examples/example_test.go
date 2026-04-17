package examples

import (
	"fmt"
	"time"

	"github.com/nekrassov01/filter"
)

var stats = Stats{
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

var desired = `
	Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (
		BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s'
	) && (
		HitPoint > "50" && MagicPoint > 100 && LifePoint != 0
	) && (
		Magic >= 20 || !(Speed < 20)
	)
`

func Example() {
	expr, err := filter.Parse(desired)
	if err != nil {
		fmt.Println(err)
		return
	}
	ok, err := expr.Eval(&stats)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("matched:", ok)
	// Output:
	// matched: true
}
