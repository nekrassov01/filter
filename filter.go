package filter

// Resolver resolves an identifier in an expression to its Value.
// The second result reports whether the identifier is known; unknown
// identifiers are reported by Eval as an error with the position.
type Resolver interface {
	Resolve(name string) (Value, bool)
}

// Expr represents a parsed expression.
type Expr struct {
	expr
}

// Parse parses input into an Expr that can be evaluated repeatedly.
func Parse(input string) (*Expr, error) {
	e, err := parse(input)
	if err != nil {
		return nil, err
	}
	return &Expr{expr: e}, nil
}

// MustParse is like Parse but panics if the input cannot be parsed.
func MustParse(input string) *Expr {
	e, err := parse(input)
	if err != nil {
		panic(err)
	}
	return &Expr{expr: e}
}

// Eval evaluates the expression against the values provided by r.
// An Expr can be evaluated by multiple goroutines at the same time.
func (e *Expr) Eval(r Resolver) (bool, error) {
	return eval(&e.expr, r)
}
