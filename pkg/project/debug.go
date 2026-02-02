package project

import (
	"context"
	"time"
)

// Debug the project.
func (p *Project) Debug(ctx context.Context) any {
	cld := map[string]any{}

	nctx, c := context.WithTimeout(context.Background(), time.Second*30)
	defer c()

	rs, ds, _ := p.matchRequirements(nctx)

	clv := "valid"
	if err := p.Validate(ctx); err != nil {
		clv = err.Error()
	}
	cld["project"] = struct {
		Valid string
	}{
		Valid: clv,
	}

	cldComponents := map[string]any{}
	for id, c := range p.Components {
		cv := "valid"
		if err := c.Validate(ctx); err != nil {
			cv = err.Error()
		}

		cldComponents[id] = struct {
			Name  string `json:"name"`
			Debug any    `json:"debug"`
			Valid string `json:"valid"`
		}{
			Name:  c.Name(),
			Debug: c.Debug(),
			Valid: cv,
		}
	}
	cld["components"] = cldComponents

	cldRequirements := map[string]any{}
	for _, r := range rs {
		cldRD := "<un-met>"
		if d := r.Matched(ctx); d != nil {
			cldRD = d.ID()
		}

		cldRequirements[r.ID()] = struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Dependency  string `json:"dependency"`
		}{
			ID:          r.ID(),
			Description: r.Describe(),
			Dependency:  cldRD,
		}
	}
	cld["requirements"] = cldRequirements

	cldDependencies := map[string]any{}
	for _, d := range ds {
		dv := "valid"
		if err := d.Validate(ctx); err != nil {
			dv = err.Error()
		}

		cldDependencies[d.ID()] = struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Valid       string `json:"valid"`
		}{
			ID:          d.ID(),
			Description: d.Describe(),
			Valid:       dv,
		}

	}
	cld["dependencies"] = cldDependencies

	return cld
}
