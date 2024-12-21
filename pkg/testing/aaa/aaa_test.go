package aaa

func ExampleNew() {
	type (
		Vars struct {
			a, b int
		}
		Result struct {
			sum int
		}
	)

	aaa := New[Vars, Result]()

	aaa.Arrange(func() *Vars {
		return &Vars{
			a: 1,
			b: 2,
		}
	}).Act(func(v *Vars) *Result {
		sum := add(v.a, v.b)
		return &Result{
			sum: sum,
		}
	}).Assert(func(_ *Vars, r *Result) {
		if r.sum != 3 {
			panic("sum is not 3")
		}
	})
	// Output:
}

func add(a, b int) int { return a + b }
