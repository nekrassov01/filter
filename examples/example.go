package examples

import (
	"fmt"
	"time"
)

// Stats represents the statistics of a character.
type Stats struct {
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

// GetField returns the value of the specified field.
func (o *Stats) GetField(key string) (any, error) {
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
