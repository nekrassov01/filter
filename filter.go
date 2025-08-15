package filter

// Target represents an entity that can provide field values for evaluation.
type Target interface {
	GetField(key string) (any, error)
}

// Expr represents an expression that can be evaluated against a Target.
type Expr interface {
	// Eval evaluates the expression for the given Target.
	Eval(t Target) (bool, error)
}

// Parse parses a string expression into an Expr.
func Parse(input string) (Expr, error) {
	if input == "" {
		return nil, parseError("empty input")
	}
	parser := &parser{
		lexer:  newLexer(input),
		nodes:  make([]node, 0, 32),
		idents: make(map[string]struct{}, 8),
	}
	root, err := parser.parseExpr()
	if err != nil {
		return nil, err
	}
	if parser.peek().typ != tokenEOF {
		return nil, parseError("unexpected token after parsing: %s", parser.peek().val)
	}
	return &expr{
		parser: parser,
		root:   root,
	}, nil
}
