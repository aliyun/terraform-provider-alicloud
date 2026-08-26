// Package sdkv2 adapts the interceptor chain to SDK v2 resources and data sources.
//
// The three CRUD forms are wrapped in the SDK's own dispatch order — plain,
// *WithoutTimeout, *Context — and only the highest-priority one that is set,
// since the SDK never reaches the others. A plain form is never promoted to
// *Context: that would silently change timeout semantics.
//
// Everything outside CRUD reaches the SDK unchanged because the wrapper starts
// from a full struct copy and reassigns only the CRUD fields.
package sdkv2

import (
	"context"
	"errors"
	"strings"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func WrapResource(name string, r *schema.Resource, chain []intercept.Interceptor) *schema.Resource {
	if r == nil || len(chain) == 0 {
		return r
	}
	cp := *r
	wrapCreate(&cp, name, chain)
	wrapRead(&cp, name, chain)
	wrapUpdate(&cp, name, chain)
	wrapDelete(&cp, name, chain)
	return &cp
}

func WrapDataSource(name string, ds *schema.Resource, chain []intercept.Interceptor) *schema.Resource {
	if ds == nil || len(chain) == 0 {
		return ds
	}
	cp := *ds
	wrapRead(&cp, name, chain)
	return &cp
}

func wrapCreate(r *schema.Resource, name string, chain []intercept.Interceptor) {
	if r.Create != nil {
		r.Create = wrapOp(name, intercept.OpCreate, chain, r.Create)
	} else if r.CreateWithoutTimeout != nil {
		r.CreateWithoutTimeout = wrapOpContext(name, intercept.OpCreate, chain, r.CreateWithoutTimeout)
	} else if r.CreateContext != nil {
		r.CreateContext = wrapOpContext(name, intercept.OpCreate, chain, r.CreateContext)
	}
}

func wrapRead(r *schema.Resource, name string, chain []intercept.Interceptor) {
	if r.Read != nil {
		r.Read = wrapOp(name, intercept.OpRead, chain, r.Read)
	} else if r.ReadWithoutTimeout != nil {
		r.ReadWithoutTimeout = wrapOpContext(name, intercept.OpRead, chain, r.ReadWithoutTimeout)
	} else if r.ReadContext != nil {
		r.ReadContext = wrapOpContext(name, intercept.OpRead, chain, r.ReadContext)
	}
}

func wrapUpdate(r *schema.Resource, name string, chain []intercept.Interceptor) {
	if r.Update != nil {
		r.Update = wrapOp(name, intercept.OpUpdate, chain, r.Update)
	} else if r.UpdateWithoutTimeout != nil {
		r.UpdateWithoutTimeout = wrapOpContext(name, intercept.OpUpdate, chain, r.UpdateWithoutTimeout)
	} else if r.UpdateContext != nil {
		r.UpdateContext = wrapOpContext(name, intercept.OpUpdate, chain, r.UpdateContext)
	}
}

func wrapDelete(r *schema.Resource, name string, chain []intercept.Interceptor) {
	if r.Delete != nil {
		r.Delete = wrapOp(name, intercept.OpDelete, chain, r.Delete)
	} else if r.DeleteWithoutTimeout != nil {
		r.DeleteWithoutTimeout = wrapOpContext(name, intercept.OpDelete, chain, r.DeleteWithoutTimeout)
	} else if r.DeleteContext != nil {
		r.DeleteContext = wrapOpContext(name, intercept.OpDelete, chain, r.DeleteContext)
	}
}

func wrapOp(name string, op intercept.Op, chain []intercept.Interceptor, inner func(*schema.ResourceData, interface{}) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		return intercept.Execute(context.Background(), chain, intercept.Call{Name: name, Op: op, Meta: meta}, func() error {
			return inner(d, meta)
		})
	}
}

func wrapOpContext(name string, op intercept.Op, chain []intercept.Interceptor, inner func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		return intercept.Around(ctx, chain, intercept.Call{Name: name, Op: op, Meta: meta},
			func() diag.Diagnostics { return inner(ctx, d, meta) }, bridge)
	}
}

var bridge = intercept.DiagBridge[diag.Diagnostics]{
	ToError:   diagsToErr,
	WithError: appendDiagError,
}

func diagsToErr(diags diag.Diagnostics) error {
	var errs diag.Diagnostics
	for _, d := range diags {
		if d.Severity == diag.Error {
			errs = append(errs, d)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return &diagsErr{diags: errs}
}

func appendDiagError(diags diag.Diagnostics, err error) diag.Diagnostics {
	if err == nil {
		return diags
	}
	// A hook that only carried the originals forward has nothing of its own to say.
	var de *diagsErr
	if errors.As(err, &de) && err.Error() == de.Error() {
		return diags
	}
	return append(diags, diag.FromErr(err)...)
}

// diagsErr carries every error-severity diagnostic across the hop through the
// interceptor contract's single error.
type diagsErr struct {
	diags diag.Diagnostics
}

func (e *diagsErr) Error() string {
	msgs := make([]string, 0, len(e.diags))
	for _, d := range e.diags {
		if d.Detail != "" {
			msgs = append(msgs, d.Summary+": "+d.Detail)
			continue
		}
		msgs = append(msgs, d.Summary)
	}
	return strings.Join(msgs, "; ")
}
