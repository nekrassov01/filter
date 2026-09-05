package benchmarks

import (
	"testing"

	"cel.dev/cel-go/cel"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/nekrassov01/filter"
	"github.com/nekrassov01/filter/examples"
)

const (
	filterExpression = `Class == "軍師" && Name =~ '^(諸葛亮|龐統|法正)' && Name != "" && (BirthDate < 0190-01-01T00:00:00Z && ATBGauge >= 20s) && (HitPoint > 50 && MagicPoint > 100 && LifePoint != 0) && Speed >= 20`
	exprExpression   = `Class == "軍師" && Name matches '^(諸葛亮|龐統|法正)' && Name != "" && (BirthDate < BirthLimit && ATBGauge >= MinGauge) && (HitPoint > 50 && MagicPoint > 100 && LifePoint != 0) && Speed >= 20`
	celExpression    = `Class == "軍師" && Name.matches('^(諸葛亮|龐統|法正)') && Name != "" && (BirthDate < timestamp("0190-01-01T00:00:00Z") && ATBGauge >= duration("20s")) && (HitPoint > 50.0 && MagicPoint > 100.0 && LifePoint != 0) && Speed >= 20`
)

var (
	// Both inputs evaluate every predicate. Only the final Speed comparison
	// changes, producing one match and one miss.
	matchInput = benchmarkInput(25)
	missInput  = benchmarkInput(15)
	celMatch   = celActivation(matchInput)
	celMiss    = celActivation(missInput)
)

func TestExpressionsAgree(t *testing.T) {
	filterProgram, err := filter.Parse(filterExpression)
	if err != nil {
		t.Fatal(err)
	}
	exprProgram, err := expr.Compile(exprExpression, expr.Env(examples.Stats{}), expr.AsBool())
	if err != nil {
		t.Fatal(err)
	}
	celProgram := compileCEL(t, newCELEnv(t), newCELOptimizer(t))

	tests := []struct {
		name  string
		input *examples.Stats
		cel   map[string]any
		want  bool
	}{
		{name: "Match", input: &matchInput, cel: celMatch, want: true},
		{name: "Miss", input: &missInput, cel: celMiss, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filterResult, err := filterProgram.Eval(test.input)
			if err != nil {
				t.Fatal(err)
			}
			var machine vm.VM
			exprResult, err := machine.Run(exprProgram, *test.input)
			if err != nil {
				t.Fatal(err)
			}
			celResult, _, err := celProgram.Eval(test.cel)
			if err != nil {
				t.Fatal(err)
			}
			if filterResult != test.want || exprResult != test.want || celResult.Value() != test.want {
				t.Fatalf("filter=%v expr=%v CEL=%v, want %v", filterResult, exprResult, celResult.Value(), test.want)
			}
		})
	}
}

// Prepare measures the public work needed to produce a reusable evaluator.
// Reusable environment and option configuration is created before measurement.

func BenchmarkPrepareFilter(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := filter.Parse(filterExpression); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPrepareExpr(b *testing.B) {
	options := []expr.Option{expr.Env(examples.Stats{}), expr.AsBool()}
	b.ReportAllocs()
	for b.Loop() {
		if _, err := expr.Compile(exprExpression, options...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPrepareCEL(b *testing.B) {
	environment := newCELEnv(b)
	optimizer := newCELOptimizer(b)
	b.ReportAllocs()
	for b.Loop() {
		compileCEL(b, environment, optimizer)
	}
}

// Eval prepares each expression once, then measures only repeated evaluation.

func BenchmarkEvalFilter(b *testing.B) {
	program, err := filter.Parse(filterExpression)
	if err != nil {
		b.Fatal(err)
	}
	b.Run("Match", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := program.Eval(&matchInput)
			if err != nil || !result {
				b.Fatal(result, err)
			}
		}
	})
	b.Run("Miss", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, err := program.Eval(&missInput)
			if err != nil || result {
				b.Fatal(result, err)
			}
		}
	})
}

func BenchmarkEvalExpr(b *testing.B) {
	program, err := expr.Compile(exprExpression, expr.Env(examples.Stats{}), expr.AsBool())
	if err != nil {
		b.Fatal(err)
	}
	b.Run("Match", func(b *testing.B) {
		var machine vm.VM
		b.ReportAllocs()
		for b.Loop() {
			result, err := machine.Run(program, matchInput)
			if err != nil || result != true {
				b.Fatal(result, err)
			}
		}
	})
	b.Run("Miss", func(b *testing.B) {
		var machine vm.VM
		b.ReportAllocs()
		for b.Loop() {
			result, err := machine.Run(program, missInput)
			if err != nil || result != false {
				b.Fatal(result, err)
			}
		}
	})
}

func BenchmarkEvalCEL(b *testing.B) {
	program := compileCEL(b, newCELEnv(b), newCELOptimizer(b))
	b.Run("Match", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, _, err := program.Eval(celMatch)
			if err != nil || result.Value() != true {
				b.Fatal(result, err)
			}
		}
	})
	b.Run("Miss", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			result, _, err := program.Eval(celMiss)
			if err != nil || result.Value() != false {
				b.Fatal(result, err)
			}
		}
	})
}

func benchmarkInput(speed int) examples.Stats {
	input := examples.SampleStats()
	input.Speed = speed
	return input
}

func newCELEnv(tb testing.TB) *cel.Env {
	tb.Helper()
	environment, err := cel.NewEnv(
		cel.Variable("Class", cel.StringType),
		cel.Variable("Name", cel.StringType),
		cel.Variable("BirthDate", cel.TimestampType),
		cel.Variable("ATBGauge", cel.DurationType),
		cel.Variable("HitPoint", cel.DoubleType),
		cel.Variable("MagicPoint", cel.DoubleType),
		cel.Variable("LifePoint", cel.IntType),
		cel.Variable("Speed", cel.IntType),
	)
	if err != nil {
		tb.Fatal(err)
	}
	return environment
}

func newCELOptimizer(tb testing.TB) *cel.StaticOptimizer {
	tb.Helper()
	folding, err := cel.NewConstantFoldingOptimizer()
	if err != nil {
		tb.Fatal(err)
	}
	optimizer, err := cel.NewStaticOptimizer(folding)
	if err != nil {
		tb.Fatal(err)
	}
	return optimizer
}

func compileCEL(tb testing.TB, environment *cel.Env, optimizer *cel.StaticOptimizer) cel.Program {
	tb.Helper()
	ast, issues := environment.Compile(celExpression)
	if issues.Err() != nil {
		tb.Fatal(issues.Err())
	}
	ast, issues = optimizer.Optimize(environment, ast)
	if issues.Err() != nil {
		tb.Fatal(issues.Err())
	}
	program, err := environment.Program(ast, cel.EvalOptions(cel.OptOptimize))
	if err != nil {
		tb.Fatal(err)
	}
	return program
}

func celActivation(input examples.Stats) map[string]any {
	return map[string]any{
		"Class":      input.Class,
		"Name":       input.Name,
		"BirthDate":  input.BirthDate,
		"ATBGauge":   input.ATBGauge,
		"HitPoint":   input.HitPoint,
		"MagicPoint": input.MagicPoint,
		"LifePoint":  input.LifePoint,
		"Speed":      input.Speed,
	}
}
