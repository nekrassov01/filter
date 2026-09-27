package examples

import (
	"net/netip"
	"time"

	"github.com/nekrassov01/filter"
)

// Stats contains the input fields used by the examples and benchmarks.
type Stats struct {
	Class      string
	Name       string
	BirthDate  time.Time
	BirthLimit time.Time
	ATBGauge   time.Duration
	MinGauge   time.Duration
	HitPoint   float64
	MagicPoint float64
	LifePoint  int
	Magic      int
	Speed      int
}

// SampleStats returns the input used by the examples and benchmarks.
func SampleStats() Stats {
	return Stats{
		Class:      "軍師",
		Name:       "諸葛亮 孔明",
		BirthDate:  time.Date(181, 7, 23, 0, 0, 0, 0, time.UTC),
		BirthLimit: time.Date(190, 1, 1, 0, 0, 0, 0, time.UTC),
		ATBGauge:   30 * time.Second,
		MinGauge:   20 * time.Second,
		HitPoint:   80,
		MagicPoint: 250,
		LifePoint:  5,
		Magic:      10,
		Speed:      25,
	}
}

// Resolve returns the value bound to the identifier.
func (o *Stats) Resolve(name string) (filter.Value, bool) {
	switch name {
	case "CLASS", "Class":
		return filter.String(o.Class), true
	case "NAME", "Name":
		return filter.String(o.Name), true
	case "BIRTH", "Birth", "BirthDate":
		return filter.Time(o.BirthDate), true
	case "BirthLimit":
		return filter.Time(o.BirthLimit), true
	case "ATB", "Atb", "ATBGauge":
		return filter.Duration(o.ATBGauge), true
	case "MinGauge":
		return filter.Duration(o.MinGauge), true
	case "HP", "Hp", "HitPoint":
		return filter.Float64(o.HitPoint), true
	case "MP", "Mp", "MagicPoint":
		return filter.Float64(o.MagicPoint), true
	case "LP", "Lp", "LifePoint":
		return filter.Int(o.LifePoint), true
	case "MAG", "Mag", "Magic":
		return filter.Int(o.Magic), true
	case "SPD", "Spd", "Speed":
		return filter.Int(o.Speed), true
	default:
		return filter.Value{}, false
	}
}

// Event has a single time value for the time literal examples.
type Event struct {
	At time.Time
}

// SampleTimeLayouts returns the event used by the time layout example.
func SampleTimeLayouts() Event {
	return Event{
		At: time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
	}
}

// Resolve returns the value bound to the identifier.
func (o *Event) Resolve(name string) (filter.Value, bool) {
	if name == "At" {
		return filter.Time(o.At), true
	}
	return filter.Value{}, false
}

// LogLine represents one application log entry.
type LogLine struct {
	Time    time.Time
	Level   string
	Status  int
	Latency time.Duration
	Path    string
}

// SampleLogLines returns the log entries used by the example.
func SampleLogLines() []LogLine {
	return []LogLine{
		{Time: time.Now(), Level: "info", Status: 200, Latency: 12 * time.Millisecond, Path: "/api/users"},
		{Time: time.Now(), Level: "error", Status: 200, Latency: 8 * time.Millisecond, Path: "/api/orders"},
		{Time: time.Now(), Level: "warn", Status: 503, Latency: 900 * time.Millisecond, Path: "/api/search"},
		{Time: time.Now(), Level: "warn", Status: 503, Latency: 900 * time.Millisecond, Path: "/health"},
	}
}

// Resolve returns the value bound to the identifier.
func (o *LogLine) Resolve(name string) (filter.Value, bool) {
	switch name {
	case "time":
		return filter.Time(o.Time), true
	case "level":
		return filter.String(o.Level), true
	case "status":
		return filter.Int(o.Status), true
	case "latency":
		return filter.Duration(o.Latency), true
	case "path":
		return filter.String(o.Path), true
	default:
		return filter.Value{}, false
	}
}

// LiteralValues contains the values used by the unquoted literal example.
type LiteralValues struct {
	Count     int
	At        time.Time
	Latency   time.Duration
	IPv4      netip.Addr
	IPv6      netip.Addr
	LinkLocal netip.Addr
	Loopback  netip.Addr
	Enabled   bool
}

// SampleLiterals returns the input used by the unquoted literal example.
func SampleLiterals() LiteralValues {
	return LiteralValues{
		Count:     42,
		At:        time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		Latency:   1500 * time.Millisecond,
		IPv4:      netip.MustParseAddr("192.0.2.1"),
		IPv6:      netip.MustParseAddr("2001:db8::1"),
		LinkLocal: netip.MustParseAddr("fe80::1%eth0"),
		Loopback:  netip.MustParseAddr("::1"),
		Enabled:   true,
	}
}

// Resolve returns the value bound to the identifier.
func (o *LiteralValues) Resolve(name string) (filter.Value, bool) {
	switch name {
	case "Count":
		return filter.Int(o.Count), true
	case "At":
		return filter.Time(o.At), true
	case "Latency":
		return filter.Duration(o.Latency), true
	case "IPv4":
		return filter.Addr(o.IPv4), true
	case "IPv6":
		return filter.Addr(o.IPv6), true
	case "LinkLocal":
		return filter.Addr(o.LinkLocal), true
	case "Loopback":
		return filter.Addr(o.Loopback), true
	case "Enabled":
		return filter.Bool(o.Enabled), true
	default:
		return filter.Value{}, false
	}
}
