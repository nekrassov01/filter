package examples

import (
	"github.com/nekrassov01/filter"

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
func (o *Stats) Resolve(name string) (filter.Value, bool) {
	switch name {
	case "CLASS", "Class":
		return filter.String(o.Class), true
	case "NAME", "Name":
		return filter.String(o.Name), true
	case "BIRTH", "Birth", "BirthDate":
		return filter.Time(o.Birth), true
	case "ATB", "Atb", "ATBGauge":
		return filter.Duration(o.ATBGauge), true
	case "HP", "Hp", "HitPoint":
		return filter.Number(o.HitPoint), true
	case "SP", "Sp", "SkillPoint":
		return filter.Number(o.SkillPoint), true
	case "MP", "Mp", "MagicPoint", "SpellPoint":
		return filter.Number(o.SpellPoint), true
	case "LP", "Lp", "LifePoint":
		return filter.Number(float64(o.LifePoint)), true
	case "STR", "Str", "Strength":
		return filter.Number(float64(o.Strength)), true
	case "STA", "Sta", "Stamina":
		return filter.Number(float64(o.Stamina)), true
	case "DEX", "Dex", "Dexterity":
		return filter.Number(float64(o.Dexterity)), true
	case "MAG", "Mag", "Magic":
		return filter.Number(float64(o.Magic)), true
	case "SPD", "Spd", "Speed":
		return filter.Number(float64(o.Speed)), true
	default:
		return filter.Value{}, false
	}
}

// Event has a single time value for the time literal examples.
type Event struct {
	At time.Time
}

// Resolve returns the value bound to the identifier.
func (o *Event) Resolve(name string) (filter.Value, bool) {
	if name == "At" {
		return filter.Time(o.At), true
	}
	return filter.Value{}, false
}
