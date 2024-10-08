package project

import (
	"context"
	"time"
)

// Debug the project.
func (cl *Project) Debug(ctx context.Context) interface{} {
	cld := map[string]interface{}{}

	nctx, c := context.WithTimeout(context.Background(), time.Second*30)
	defer c()

	rs, ds, _ := cl.matchRequirements(nctx)

	cl_v := "valid"
	if err := cl.Validate(ctx); err != nil {
		cl_v = err.Error()
	}
	cld["project"] = struct {
		Valid string
	}{
		Valid: cl_v,
	}

	cld_components := map[string]interface{}{}
	for id, c := range cl.Components {
		c_v := "valid"
		if err := c.Validate(ctx); err != nil {
			c_v = err.Error()
		}

		cld_components[id] = struct {
			Name  string      `json:"name"`
			Debug interface{} `json:"debug"`
			Valid string      `json:"valid"`
		}{
			Name:  c.Name(),
			Debug: c.Debug(),
			Valid: c_v,
		}
	}
	cld["components"] = cld_components

	cld_requirements := map[string]interface{}{}
	for _, r := range rs {
		cld_r_d := "<un-met>"
		if d := r.Matched(ctx); d != nil {
			cld_r_d = d.Id()
		}

		cld_requirements[r.Id()] = struct {
			Id          string `json:"id"`
			Description string `json:"description"`
			Dependency  string `json:"dependency"`
		}{
			Id:          r.Id(),
			Description: r.Describe(),
			Dependency:  cld_r_d,
		}
	}
	cld["requirements"] = cld_requirements

	cld_dependencies := map[string]interface{}{}
	for _, d := range ds {
		d_v := "valid"
		if err := d.Validate(ctx); err != nil {
			d_v = err.Error()
		}

		cld_dependencies[d.Id()] = struct {
			Id          string `json:"id"`
			Description string `json:"description"`
			Valid       string `json:"valid"`
		}{
			Id:          d.Id(),
			Description: d.Describe(),
			Valid:       d_v,
		}

	}
	cld["dependencies"] = cld_dependencies

	return cld
}
