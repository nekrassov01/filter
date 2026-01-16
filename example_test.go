package filter

import (
	"fmt"
	"time"
)

type personnel struct {
	Class      string
	Name       string
	Birth      time.Time
	ATBGauge   time.Duration
	HitPoint   float64
	SkillPoint float64
	SpellPoint float64
	LifePoint  int64
	Strength   int64
	Stamina    int64
	Dexterity  int64
	Magic      int64
	Speed      int64
}

func (o *personnel) GetField(key string) (any, error) {
	switch key {
	case "CLASS", "Class":
		return o.Class, nil
	case "NAME", "Name":
		return o.Name, nil
	case "BIRTH", "Birth", "BirthDate":
		return o.Birth, nil
	case "ATB", "Atb", "ActiveTimeBattleGauge":
		return o.ATBGauge, nil
	case "HP", "Hp", "HitPoint":
		return o.HitPoint, nil
	case "SP", "Sp", "SkillPoint":
		return o.SkillPoint, nil
	case "MP", "Mp", "MagicPoint", "SpellPoint":
		return o.SpellPoint, nil
	case "LP", "Lp", "LifePoint":
		return o.LifePoint, nil
	case "STR", "Str", "Strength":
		return o.Strength, nil
	case "STA", "Sta", "Stamina":
		return o.Stamina, nil
	case "DEX", "Dex", "Dexterity":
		return o.Dexterity, nil
	case "MAG", "Mag", "Magic":
		return o.Magic, nil
	case "SPD", "Spd", "Speed":
		return o.Speed, nil
	default:
		return nil, fmt.Errorf("field not found: %q", key)
	}
}

var talent = personnel{
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

func Example() {
	var desired = `
Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (
	BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s'
) && (
	HitPoint > "50" && MagicPoint > 100 && LifePoint != 0
) && (
	Magic >= 20 || !(Speed < 20)
)`

	expr, err := Parse(desired)
	if err != nil {
		fmt.Println(err)
		return
	}
	ok, err := expr.Eval(&talent)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("matched:", ok)
	// Output:
	// matched: true
}
