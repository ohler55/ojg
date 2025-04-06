// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import "github.com/ohler55/ojg"

// Plan is an assembly plan that can be described by a JSON document or a SEN
// document. The format is much like LISP but with brackets instead of
// parenthesis. A plan is evaluated by evaluating the plan function which is
// usually an 'asm' function. The plan operates on a data map which is the
// root during evaluation. The source data is in the $.src and the expected
// assembled output should be in $.asm.
type Plan struct {
	Fn
}

// NewPlan creates new place from a simplified (JSON) encoding of the
// instance.
func NewPlan(plan []any) *Plan {
	if len(plan) == 0 {
		return nil
	}
	p := Plan{}
	if name, _ := plan[0].(string); 0 < len(name) {
		if name == "asm" {
			p.Fn = asmFn
		} else if af := NewFn(name); af != nil {
			p.Fn = *af
		}
		p.Args = plan[1:]
	}
	if p.Eval == nil {
		p.Fn = asmFn
		p.Args = plan
	}
	p.compile()

	return &p
}

// Execute a plan.
func (p *Plan) Execute(root map[string]any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ojg.NewError(r)
		}
	}()
	p.Eval(root, root, p.Args...)

	return
}
