package aaa

type TestAAA[Vars, ActResult any] struct{}

type TestArrangement[Vars, ActResult any] struct {
	vars *Vars
}

type TestAction[Vars, ActResult any] struct {
	vars   *Vars
	result *ActResult
}

type ArrangeFunc[Vars any] func() *Vars

type ActFunc[Vars, ActResult any] func(*Vars) *ActResult

type AssertFunc[Vars, ActResult any] func(*Vars, *ActResult)

func New[Vars, ActResult any]() *TestAAA[Vars, ActResult] {
	return &TestAAA[Vars, ActResult]{}
}

func (t *TestAAA[Vars, ActResult]) Arrange(f ArrangeFunc[Vars]) *TestArrangement[Vars, ActResult] {
	return &TestArrangement[Vars, ActResult]{
		vars: f(),
	}
}

func (t *TestArrangement[Vars, ActResult]) Act(f ActFunc[Vars, ActResult]) *TestAction[Vars, ActResult] {
	return &TestAction[Vars, ActResult]{
		vars:   t.vars,
		result: f(t.vars),
	}
}

func (t *TestAction[Vars, ActResult]) Assert(f AssertFunc[Vars, ActResult]) {
	f(t.vars, t.result)
}
