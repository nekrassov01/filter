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

// Eval evaluates the expression against the values provided by r.
func (e *Expr) Eval(r Resolver) (bool, error) {
	return eval(&e.expr, r)
}
