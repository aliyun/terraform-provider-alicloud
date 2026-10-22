package sdkv2

import (
	"context"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/provider/intercept"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// hasTopLevelTags reports whether r exposes the tag surface the filter knows how to rewrite: a
// top level tags map. Nested tag blocks belong to their own resource and stay untouched.
func hasTopLevelTags(r *schema.Resource) bool {
	s, ok := r.Schema["tags"]
	return ok && s != nil && s.Type == schema.TypeMap
}

// wrapReadWithTagFilter wraps only the Read form the SDK will actually dispatch to, keeping the
// interceptor chain around it and filtering the returned tags afterwards.
func wrapReadWithTagFilter(r *schema.Resource, name string, chain []intercept.Interceptor) {
	switch {
	case r.Read != nil:
		r.Read = wrapReadOp(name, chain, r.Read)
	case r.ReadWithoutTimeout != nil:
		r.ReadWithoutTimeout = wrapReadOpContext(name, chain, r.ReadWithoutTimeout)
	case r.ReadContext != nil:
		r.ReadContext = wrapReadOpContext(name, chain, r.ReadContext)
	}
}

// wrapCreateWithTagFilter wraps the Create form the SDK dispatches to and filters the tags of the
// state Apply returns. A resource usually finishes Create by calling its own Read, which bypasses the
// wrapped Read, so the returned state would otherwise carry the ignored keys again until the next
// refresh.
func wrapCreateWithTagFilter(r *schema.Resource, name string, chain []intercept.Interceptor) {
	switch {
	case r.Create != nil:
		r.Create = wrapMutationWithTagFilter(name, intercept.OpCreate, chain, r.Create)
	case r.CreateWithoutTimeout != nil:
		r.CreateWithoutTimeout = wrapMutationWithTagFilterContext(name, intercept.OpCreate, chain, r.CreateWithoutTimeout)
	case r.CreateContext != nil:
		r.CreateContext = wrapMutationWithTagFilterContext(name, intercept.OpCreate, chain, r.CreateContext)
	}
}

// wrapUpdateWithTagFilter is the Update counterpart of wrapCreateWithTagFilter.
func wrapUpdateWithTagFilter(r *schema.Resource, name string, chain []intercept.Interceptor) {
	switch {
	case r.Update != nil:
		r.Update = wrapMutationWithTagFilter(name, intercept.OpUpdate, chain, r.Update)
	case r.UpdateWithoutTimeout != nil:
		r.UpdateWithoutTimeout = wrapMutationWithTagFilterContext(name, intercept.OpUpdate, chain, r.UpdateWithoutTimeout)
	case r.UpdateContext != nil:
		r.UpdateContext = wrapMutationWithTagFilterContext(name, intercept.OpUpdate, chain, r.UpdateContext)
	}
}

func wrapMutationWithTagFilter(name string, op intercept.Op, chain []intercept.Interceptor, inner func(*schema.ResourceData, interface{}) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		err := intercept.Execute(context.Background(), chain, intercept.Call{Name: name, Op: op, Meta: meta}, func() error {
			return inner(d, meta)
		})
		if err != nil {
			return err
		}
		return filterIgnoredTags(d, meta)
	}
}

func wrapMutationWithTagFilterContext(name string, op intercept.Op, chain []intercept.Interceptor, inner func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := intercept.Around(ctx, chain, intercept.Call{Name: name, Op: op, Meta: meta},
			func() diag.Diagnostics { return inner(ctx, d, meta) }, bridge)
		if diags.HasError() {
			return diags
		}
		if err := filterIgnoredTags(d, meta); err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		return diags
	}
}

func wrapReadOp(name string, chain []intercept.Interceptor, inner func(*schema.ResourceData, interface{}) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		err := intercept.Execute(context.Background(), chain, intercept.Call{Name: name, Op: intercept.OpRead, Meta: meta}, func() error {
			return inner(d, meta)
		})
		if err != nil {
			return err
		}
		return filterIgnoredTags(d, meta)
	}
}

func wrapReadOpContext(name string, chain []intercept.Interceptor, inner func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics) func(context.Context, *schema.ResourceData, interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		diags := intercept.Around(ctx, chain, intercept.Call{Name: name, Op: intercept.OpRead, Meta: meta},
			func() diag.Diagnostics { return inner(ctx, d, meta) }, bridge)
		if diags.HasError() {
			return diags
		}
		if err := filterIgnoredTags(d, meta); err != nil {
			return append(diags, diag.FromErr(err)...)
		}
		return diags
	}
}

// filterIgnoredTags drops the keys the provider block ignores, plus the system prefixes every tag
// helper already hides, from the state a successful Read produced. It runs even when the
// interceptor chain is empty, because filtering belongs to the provider configuration and not to
// any interceptor.
func filterIgnoredTags(d *schema.ResourceData, meta interface{}) error {
	if d == nil {
		return nil
	}
	client, ok := meta.(*connectivity.AliyunClient)
	if !ok || client == nil || client.IgnoreTags.Empty() {
		return nil
	}
	current, ok := d.Get("tags").(map[string]interface{})
	if !ok || len(current) == 0 {
		return nil
	}
	filtered := make(map[string]interface{}, len(current))
	dropped := false
	for key, value := range current {
		if client.IgnoreTags.Match(key) || tags.TagIgnored(key, value) {
			dropped = true
			continue
		}
		filtered[key] = value
	}
	if !dropped {
		return nil
	}
	return d.Set("tags", filtered)
}
