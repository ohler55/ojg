// Copyright (c) 2021, Peter Ohler, All rights reserved.

package asm

import "fmt"

type Plan struct {
	Fn
}

func NewPlan(plan []interface{}) *Plan {
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
	if p.Fn.Eval == nil {
		p.Fn = asmFn
		p.Args = plan
	}
	p.compile()

	return &p
}

func (p *Plan) Execute(root map[string]interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if err, _ = r.(error); err == nil {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	p.Eval(root, root, p.Args...)

	return
}
