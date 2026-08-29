package benchmarks

import (
	"strings"
	"testing"
	"time"

	"github.com/nekrassov01/filter"
	"github.com/nekrassov01/filter/examples"
)

var (
	statsASCII = examples.Stats{
		Class:      "Knight",
		Name:       "William Marshal",
		Birth:      time.Date(1146, 1, 1, 0, 0, 0, 0, time.UTC),
		ATBGauge:   time.Second * 30,
		HitPoint:   120,
		SkillPoint: 40,
		SpellPoint: 0,
		LifePoint:  7,
		Strength:   30,
		Stamina:    28,
		Dexterity:  18,
		Magic:      2,
		Speed:      15,
	}
	statsUnicode = examples.Stats{
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
)

var (
	simpleASCII     = `Class == "Knight"`
	heavyASCII      = `Class == "Knight" && Name =~ '^(William|Richard|Geoffrey)' && Name != "" && (BirthDate < '1200-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s') && (HitPoint > "50" && SkillPoint > 30 && LifePoint != 0) && (Strength >= 20 || !(Speed < 20))`
	repeatedASCII   = repeatInput(heavyASCII, 30)
	simpleUnicode   = `Class == "軍師"`
	heavyUnicode    = `Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (BirthDate < '0190-01-01T00:00:00Z' && ActiveTimeBattleGauge >= '20s') && (HitPoint > "50" && MagicPoint > 100 && LifePoint != 0) && (Magic >= 20 || !(Speed < 20))`
	repeatedUnicode = repeatInput(heavyUnicode, 30)
)

func BenchmarkParseASCIISimple(b *testing.B) {
	benchmarkParse(b, simpleASCII)
}

func BenchmarkEvalASCIISimple(b *testing.B) {
	benchmarkEval(b, simpleASCII, &statsASCII)
}

func BenchmarkParseASCIIHeavy(b *testing.B) {
	benchmarkParse(b, heavyASCII)
}

func BenchmarkEvalASCIIHeavy(b *testing.B) {
	benchmarkEval(b, heavyASCII, &statsASCII)
}

func BenchmarkParseASCIIRepeated(b *testing.B) {
	benchmarkParse(b, repeatedASCII)
}

func BenchmarkEvalASCIIRepeated(b *testing.B) {
	benchmarkEval(b, repeatedASCII, &statsASCII)
}

func BenchmarkParseUnicodeSimple(b *testing.B) {
	benchmarkParse(b, simpleUnicode)
}

func BenchmarkEvalUnicodeSimple(b *testing.B) {
	benchmarkEval(b, simpleUnicode, &statsUnicode)
}

func BenchmarkParseUnicodeHeavy(b *testing.B) {
	benchmarkParse(b, heavyUnicode)
}

func BenchmarkEvalUnicodeHeavy(b *testing.B) {
	benchmarkEval(b, heavyUnicode, &statsUnicode)
}

func BenchmarkParseUnicodeRepeated(b *testing.B) {
	benchmarkParse(b, repeatedUnicode)
}

func BenchmarkEvalUnicodeRepeated(b *testing.B) {
	benchmarkEval(b, repeatedUnicode, &statsUnicode)
}

func benchmarkParse(b *testing.B, input string) {
	for b.Loop() {
		if _, err := filter.Parse(input); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkEval(b *testing.B, input string, r filter.Resolver) {
	expr, err := filter.Parse(input)
	if err != nil {
		b.Fatal(err)
	}
	for b.Loop() {
		if ok, err := expr.Eval(r); !ok || err != nil {
			b.Fatal(ok, err)
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
