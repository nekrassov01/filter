package examples

import (
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

// Resolve returns the value bound to the identifier.
func (o *Stats) Resolve(name string) (any, bool) {
	switch name {
	case "CLASS", "Class":
		return o.Class, true
	case "NAME", "Name":
		return o.Name, true
	case "BIRTH", "Birth", "BirthDate":
		return o.Birth, true
	case "ATB", "Atb", "ActiveTimeBattleGauge":
		return o.ATBGauge, true
	case "HP", "Hp", "HitPoint":
		return o.HitPoint, true
	case "SP", "Sp", "SkillPoint":
		return o.SkillPoint, true
	case "MP", "Mp", "MagicPoint", "SpellPoint":
		return o.SpellPoint, true
	case "LP", "Lp", "LifePoint":
		return o.LifePoint, true
	case "STR", "Str", "Strength":
		return o.Strength, true
	case "STA", "Sta", "Stamina":
		return o.Stamina, true
	case "DEX", "Dex", "Dexterity":
		return o.Dexterity, true
	case "MAG", "Mag", "Magic":
		return o.Magic, true
	case "SPD", "Spd", "Speed":
		return o.Speed, true
	default:
		return nil, false
	}
}

// Event has a single time value for the time literal examples.
type Event struct {
	At time.Time
}

// Resolve returns the value bound to the identifier.
func (o *Event) Resolve(name string) (any, bool) {
	if name == "At" {
		return o.At, true
	}
	return nil, false
}
