package benchmarks

import (
	"strings"
	"testing"
	"time"

	"github.com/nekrassov01/filter"
	"github.com/nekrassov01/filter/examples"
)

var stats = examples.Stats{
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

var simple = `Class == "軍師"`

var heavy = `
	Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (
		BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s'
	) && (
		HitPoint > "50" && MagicPoint > 100 && LifePoint != 0
	) && (
		Magic >= 20 || !(Speed < 20)
	)
`

func BenchmarkParseSimple(b *testing.B) {
	for b.Loop() {
		if _, err := filter.Parse(simple); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalSimple(b *testing.B) {
	expr, err := filter.Parse(simple)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := expr.Eval(&stats); !ok || err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseHeavy(b *testing.B) {
	for b.Loop() {
		if _, err := filter.Parse(heavy); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalHeavy(b *testing.B) {
	expr, err := filter.Parse(heavy)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := expr.Eval(&stats); !ok || err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseRepeated(b *testing.B) {
	input := repeatInput(heavy, 30)
	for b.Loop() {
		if _, err := filter.Parse(input); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvalRepeated(b *testing.B) {
	input := repeatInput(heavy, 30)
	expr, err := filter.Parse(input)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := expr.Eval(&stats); !ok || err != nil {
			b.Fatal(err)
		}
	}
}

func repeatInput(input string, n int) string {
	if n <= 0 {
		return input
	}
	var sb strings.Builder
	sb.Grow(len(input) + n*(len(input)+2))
	sb.WriteString(input)
	for range n {
		sb.WriteString("&&")
		sb.WriteString(input)
	}
	return sb.String()
}
