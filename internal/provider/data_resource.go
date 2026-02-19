// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &DataResource{}
var _ resource.ResourceWithImportState = &DataResource{}

func NewDataResource() resource.Resource {
	return &DataResource{}
}

type DataResource struct{}

type DataResourceModel struct {
	Id              types.String `tfsdk:"id"`
	InputWo         types.String `tfsdk:"input_wo"`
	InputWoVersion  types.Number `tfsdk:"input_wo_version"`
	TriggersReplace types.List   `tfsdk:"triggers_replace"`
	Output          types.String `tfsdk:"output"`
}

func (r *DataResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data"
}

func (r *DataResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A resource that mimics terraform_data with write-only input support.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"input_wo": schema.StringAttribute{
				MarkdownDescription: "Write-only input that accepts ephemeral values. The value is stored in `output`.",
				Optional:            true,
				Sensitive:           true,
			},
			"input_wo_version": schema.NumberAttribute{
				MarkdownDescription: "When changed, triggers re-capture of the input_wo value.",
				Optional:            true,
			},
			"triggers_replace": schema.ListAttribute{
				MarkdownDescription: "Forces replacement when value changes.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"output": schema.StringAttribute{
				MarkdownDescription: "The captured value from input_wo.",
				Computed:            true,
			},
		},
	}
}

func (r *DataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *DataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DataResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(generateRandomID())

	if !data.InputWo.IsNull() {
		data.Output = data.InputWo
	}

	tflog.Trace(ctx, "created a wo_data resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DataResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DataResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.InputWo.IsNull() {
		data.Output = data.InputWo
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DataResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "writeonly_data_" + string(b)
}
