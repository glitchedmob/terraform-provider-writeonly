package provider

import (
	"context"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DataResource{}

func NewDataResource() resource.Resource {
	return &DataResource{}
}

type DataResource struct{}

type DataResourceModel struct {
	Id              types.String `tfsdk:"id"`
	InputWo         types.String `tfsdk:"input_wo"`
	InputWoVersion  types.Int64  `tfsdk:"input_wo_version"`
	TriggersReplace types.List   `tfsdk:"triggers_replace"`
	Output          types.String `tfsdk:"output"`
}

func (r *DataResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data"
}

func (r *DataResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"input_wo": schema.StringAttribute{
				Optional:  true,
				WriteOnly: true,
			},
			"input_wo_version": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"triggers_replace": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"output": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *DataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *DataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DataResourceModel
	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Id.IsNull() {
		plan.Id = types.StringValue(generateRandomID())
	}

	if !plan.InputWo.IsNull() {
		plan.Output = plan.InputWo
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DataResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *DataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DataResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DataResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	versionChanged := !plan.InputWoVersion.Equal(state.InputWoVersion)

	if versionChanged {
		var config DataResourceModel
		diags = req.Config.Get(ctx, &config)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		plan.Output = config.InputWo
	} else {
		plan.Output = state.Output
	}

	plan.Id = state.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *DataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	req.State.RemoveResource(ctx)
}

func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "writeonly_data_" + string(b)
}
