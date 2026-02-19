package provider

import (
	"context"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
}

func (r *DataResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
}

func (r *DataResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

}

func (r *DataResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *DataResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *DataResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

func (r *DataResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}

func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return "writeonly_data_" + string(b)
}
