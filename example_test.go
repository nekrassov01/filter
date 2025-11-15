package filter

import (
	"fmt"
	"time"
)

type exampleTarget struct {
	Class          string
	Name           string
	HitPoint       float64
	SkillPoint     float64
	SpellPoint     float64
	LifePoint      int64
	Strength       int64
	Stamina        int64
	Dexterity      int64
	Magic          int64
	Speed          int64
	Sustainability time.Time
	Duration       time.Duration
}

func (o *exampleTarget) GetField(key string) (any, error) {
	switch key {
	case "Class":
		return o.Class, nil
	case "Name":
		return o.Name, nil
	case "HP":
		return o.HitPoint, nil
	case "SP":
		return o.SkillPoint, nil
	case "MP":
		return o.SpellPoint, nil
	case "LP":
		return o.LifePoint, nil
	case "STR":
		return o.Strength, nil
	case "STA":
		return o.Stamina, nil
	case "DEX":
		return o.Dexterity, nil
	case "MAG":
		return o.Magic, nil
	case "SPD":
		return o.Speed, nil
	case "SUS":
		return o.Sustainability, nil
	case "DUR":
		return o.Duration, nil
	default:
		return nil, fmt.Errorf("field not found: %q", key)
	}
}

var exampleInput = `SUS >= "2023-01-01T00:00:00Z"`

var exampleObject = exampleTarget{
	Class:          "軍師",
	Name:           "諸葛亮 孔明",
	HitPoint:       80,
	SkillPoint:     0,
	SpellPoint:     250,
	LifePoint:      5,
	Strength:       10,
	Stamina:        10,
	Dexterity:      10,
	Magic:          25,
	Speed:          25,
	Sustainability: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
	Duration:       time.Hour,
}

func Example() {
	expr, err := Parse(exampleInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	ok, err := expr.Eval(&exampleObject)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("matched:", ok)
	// Output:
	// matched: true
}
