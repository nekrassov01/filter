// The same predicates run through filter, expr, and CEL, each in its fastest
// form over the examples.Stats fields: expr reads the time and duration
// bounds from variables, since its date and duration calls are evaluated on
// every run; CEL folds constants and precompiles regular expressions.

package benchmarks

import (
	"testing"
	"time"

	"cel.dev/cel-go/cel"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/nekrassov01/filter"
	"github.com/nekrassov01/filter/examples"
)

var (
	filterASCII = examples.Stats{
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
	filterUnicode = examples.Stats{
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
	exprASCII = exprEnv{
		Stats:      filterASCII,
		BirthLimit: time.Date(1200, 1, 1, 0, 0, 0, 0, time.UTC),
		MinGauge:   20 * time.Second,
	}
	exprUnicode = exprEnv{
		Stats:      filterUnicode,
		BirthLimit: time.Date(190, 1, 1, 0, 0, 0, 0, time.UTC),
		MinGauge:   20 * time.Second,
	}
	celASCII = map[string]any{
		"Class":      filterASCII.Class,
		"Name":       filterASCII.Name,
		"BirthDate":  filterASCII.Birth,
		"ATBGauge":   filterASCII.ATBGauge,
		"HitPoint":   filterASCII.HitPoint,
		"SkillPoint": filterASCII.SkillPoint,
		"MagicPoint": filterASCII.SpellPoint,
		"LifePoint":  filterASCII.LifePoint,
		"Strength":   filterASCII.Strength,
		"Magic":      filterASCII.Magic,
		"Speed":      filterASCII.Speed,
	}
	celUnicode = map[string]any{
		"Class":      filterUnicode.Class,
		"Name":       filterUnicode.Name,
		"BirthDate":  filterUnicode.Birth,
		"ATBGauge":   filterUnicode.ATBGauge,
		"HitPoint":   filterUnicode.HitPoint,
		"SkillPoint": filterUnicode.SkillPoint,
		"MagicPoint": filterUnicode.SpellPoint,
		"LifePoint":  filterUnicode.LifePoint,
		"Strength":   filterUnicode.Strength,
		"Magic":      filterUnicode.Magic,
		"Speed":      filterUnicode.Speed,
	}
)

var (
	filterSimpleASCII   = `Class == "Knight"`
	filterHeavyASCII    = `Class == "Knight" && Name =~ '^(William|Richard|Geoffrey)' && Name != "" && (BirthDate < 1200-01-01T00:00:00Z && ATBGauge >= 20s) && (HitPoint > 50 && SkillPoint > 30 && LifePoint != 0) && (Strength >= 20 || !(Speed < 20))`
	filterSimpleUnicode = `Class == "軍師"`
	filterHeavyUnicode  = `Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (BirthDate < 0190-01-01T00:00:00Z && ATBGauge >= 20s) && (HitPoint > 50 && MagicPoint > 100 && LifePoint != 0) && (Magic >= 20 || !(Speed < 20))`
	exprSimpleASCII     = `Class == "Knight"`
	exprHeavyASCII      = `Class == "Knight" && Name matches '^(William|Richard|Geoffrey)' && Name != "" && (Birth < BirthLimit && ATBGauge >= MinGauge) && (HitPoint > 50 && SkillPoint > 30 && LifePoint != 0) && (Strength >= 20 || !(Speed < 20))`
	exprSimpleUnicode   = `Class == "軍師"`
	exprHeavyUnicode    = `Class == "軍師" && Name matches '^(諸葛亮|龐統|法正)' && Name != "" && (Birth < BirthLimit && ATBGauge >= MinGauge) && (HitPoint > 50 && SpellPoint > 100 && LifePoint != 0) && (Magic >= 20 || !(Speed < 20))`
	celSimpleASCII      = `Class == "Knight"`
	celHeavyASCII       = `Class == "Knight" && Name.matches('^(William|Richard|Geoffrey)') && Name != "" && (BirthDate < timestamp("1200-01-01T00:00:00Z") && ATBGauge >= duration("20s")) && (HitPoint > 50.0 && SkillPoint > 30.0 && LifePoint != 0) && (Strength >= 20 || !(Speed < 20))`
	celSimpleUnicode    = `Class == "軍師"`
	celHeavyUnicode     = `Class == "軍師" && Name.matches('^(諸葛亮|龐統|法正)') && Name != "" && (BirthDate < timestamp("0190-01-01T00:00:00Z") && ATBGauge >= duration("20s")) && (HitPoint > 50.0 && MagicPoint > 100.0 && LifePoint != 0) && (Magic >= 20 || !(Speed < 20))`
)

func BenchmarkParseASCIISimpleFilter(b *testing.B) {
	benchmarkParseFilter(b, filterSimpleASCII)
}

func BenchmarkEvalASCIISimpleFilter(b *testing.B) {
	benchmarkEvalFilter(b, filterSimpleASCII, &filterASCII)
}

func BenchmarkParseASCIIHeavyFilter(b *testing.B) {
	benchmarkParseFilter(b, filterHeavyASCII)
}

func BenchmarkEvalASCIIHeavyFilter(b *testing.B) {
	benchmarkEvalFilter(b, filterHeavyASCII, &filterASCII)
}

func BenchmarkParseUnicodeSimpleFilter(b *testing.B) {
	benchmarkParseFilter(b, filterSimpleUnicode)
}

func BenchmarkEvalUnicodeSimpleFilter(b *testing.B) {
	benchmarkEvalFilter(b, filterSimpleUnicode, &filterUnicode)
}

func BenchmarkParseUnicodeHeavyFilter(b *testing.B) {
	benchmarkParseFilter(b, filterHeavyUnicode)
}

func BenchmarkEvalUnicodeHeavyFilter(b *testing.B) {
	benchmarkEvalFilter(b, filterHeavyUnicode, &filterUnicode)
}

func BenchmarkParseASCIISimpleExpr(b *testing.B) {
	benchmarkParseExpr(b, exprSimpleASCII)
}

func BenchmarkEvalASCIISimpleExpr(b *testing.B) {
	benchmarkEvalExpr(b, exprSimpleASCII, exprASCII)
}

func BenchmarkParseASCIIHeavyExpr(b *testing.B) {
	benchmarkParseExpr(b, exprHeavyASCII)
}

func BenchmarkEvalASCIIHeavyExpr(b *testing.B) {
	benchmarkEvalExpr(b, exprHeavyASCII, exprASCII)
}

func BenchmarkParseUnicodeSimpleExpr(b *testing.B) {
	benchmarkParseExpr(b, exprSimpleUnicode)
}

func BenchmarkEvalUnicodeSimpleExpr(b *testing.B) {
	benchmarkEvalExpr(b, exprSimpleUnicode, exprUnicode)
}

func BenchmarkParseUnicodeHeavyExpr(b *testing.B) {
	benchmarkParseExpr(b, exprHeavyUnicode)
}

func BenchmarkEvalUnicodeHeavyExpr(b *testing.B) {
	benchmarkEvalExpr(b, exprHeavyUnicode, exprUnicode)
}

func BenchmarkParseASCIISimpleCEL(b *testing.B) {
	benchmarkParseCEL(b, celSimpleASCII)
}

func BenchmarkEvalASCIISimpleCEL(b *testing.B) {
	benchmarkEvalCEL(b, celSimpleASCII, celASCII)
}

func BenchmarkParseASCIIHeavyCEL(b *testing.B) {
	benchmarkParseCEL(b, celHeavyASCII)
}

func BenchmarkEvalASCIIHeavyCEL(b *testing.B) {
	benchmarkEvalCEL(b, celHeavyASCII, celASCII)
}

func BenchmarkParseUnicodeSimpleCEL(b *testing.B) {
	benchmarkParseCEL(b, celSimpleUnicode)
}

func BenchmarkEvalUnicodeSimpleCEL(b *testing.B) {
	benchmarkEvalCEL(b, celSimpleUnicode, celUnicode)
}

func BenchmarkParseUnicodeHeavyCEL(b *testing.B) {
	benchmarkParseCEL(b, celHeavyUnicode)
}

func BenchmarkEvalUnicodeHeavyCEL(b *testing.B) {
	benchmarkEvalCEL(b, celHeavyUnicode, celUnicode)
}

func benchmarkParseFilter(b *testing.B, input string) {
	for b.Loop() {
		if _, err := filter.Parse(input); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkEvalFilter(b *testing.B, input string, r filter.Resolver) {
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

func benchmarkParseExpr(b *testing.B, input string) {
	for b.Loop() {
		if _, err := expr.Compile(input, expr.Env(exprEnv{}), expr.AsBool()); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkEvalExpr(b *testing.B, input string, e exprEnv) {
	program, err := expr.Compile(input, expr.Env(exprEnv{}), expr.AsBool())
	if err != nil {
		b.Fatal(err)
	}
	var v vm.VM
	for b.Loop() {
		if out, err := v.Run(program, e); err != nil || out != true {
			b.Fatal(out, err)
		}
	}
}

func benchmarkParseCEL(b *testing.B, input string) {
	e := celEnv(b)
	for b.Loop() {
		celProgram(b, e, input)
	}
}

func benchmarkEvalCEL(b *testing.B, input string, activation map[string]any) {
	program := celProgram(b, celEnv(b), input)
	for b.Loop() {
		if out, _, err := program.Eval(activation); err != nil || out.Value() != true {
			b.Fatal(out, err)
		}
	}
}

// exprEnv adds the time and duration bounds to Stats as variables, the
// cheapest form expr offers for values it has no literal for.
type exprEnv struct {
	examples.Stats

	BirthLimit time.Time
	MinGauge   time.Duration
}

func celEnv(b *testing.B) *cel.Env {
	e, err := cel.NewEnv(
		cel.Variable("Class", cel.StringType),
		cel.Variable("Name", cel.StringType),
		cel.Variable("BirthDate", cel.TimestampType),
		cel.Variable("ATBGauge", cel.DurationType),
		cel.Variable("HitPoint", cel.DoubleType),
		cel.Variable("SkillPoint", cel.DoubleType),
		cel.Variable("MagicPoint", cel.DoubleType),
		cel.Variable("LifePoint", cel.IntType),
		cel.Variable("Strength", cel.IntType),
		cel.Variable("Magic", cel.IntType),
		cel.Variable("Speed", cel.IntType),
	)
	if err != nil {
		b.Fatal(err)
	}
	return e
}

func celProgram(b *testing.B, e *cel.Env, input string) cel.Program {
	ast, iss := e.Compile(input)
	if iss.Err() != nil {
		b.Fatal(iss.Err())
	}
	folding, err := cel.NewConstantFoldingOptimizer()
	if err != nil {
		b.Fatal(err)
	}
	optimizer, err := cel.NewStaticOptimizer(folding)
	if err != nil {
		b.Fatal(err)
	}
	ast, iss = optimizer.Optimize(e, ast)
	if iss.Err() != nil {
		b.Fatal(iss.Err())
	}
	program, err := e.Program(ast, cel.EvalOptions(cel.OptOptimize))
	if err != nil {
		b.Fatal(err)
	}
	return program
}
