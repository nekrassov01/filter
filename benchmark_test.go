package filter

import (
	"fmt"
	"testing"
)

type benchTarget struct {
	Class      string
	Name       string
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

func (o *benchTarget) GetField(key string) (any, error) {
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
	default:
		return nil, fmt.Errorf("field not found: %q", key)
	}
}

var benchInput = `Class == "軍師"
&&
Name =~ '孔明'
&&
Name != ""
&&
(
	HP > "50"
	&&
	MP > 100
	&&
	LP != 0
)
&&
(
	MAG >= 20
	||
	!
	(
		SPD < 20
	)
)`

func BenchmarkParse(b *testing.B) {
	for b.Loop() {
		_, err := Parse(benchInput)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEval(b *testing.B) {
	target := benchTarget{
		Class:      "軍師",
		Name:       "諸葛亮 孔明",
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
	e, err := Parse(benchInput)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := e.Eval(&target); !ok || err != nil {
			b.Fatal(err)
		}
	}
}
