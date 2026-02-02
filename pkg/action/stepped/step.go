/*
Package stepped a Phase that is a sequence of step calls
*/
package stepped

import "context"

type Steps []Step

// Add a Step to the set of Steps.
func (ss *Steps) Add(sa Step) {
	*ss = append(*ss, sa)
}

// Merge two Steps together.
func (ss *Steps) Merge(ssa Steps) {
	*ss = append(*ss, ssa...)
}

// Step in a stepped Phase.
type Step interface {
	ID() string
	Run(context.Context) error
}
