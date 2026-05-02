package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &StorageContainerResource{}
var _ resource.ResourceWithConfigure = &StorageContainerResource{}

type StorageContainerResource struct {
	subscriptionID string
}

type StorageContainerModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	StorageAccountName   types.String `tfsdk:"storage_account_name"`
	ContainerAccessType  types.String `tfsdk:"container_access_type"`
	Metadata             types.Map    `tfsdk:"metadata"`
	HasImmutabilityPolicy types.Bool  `tfsdk:"has_immutability_policy"`
	HasLegalHold         types.Bool   `tfsdk:"has_legal_hold"`
}

func NewStorageContainerResource() resource.Resource {
	return &StorageContainerResource{}
}

func (r *StorageContainerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_container"
}

func (r *StorageContainerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Storage Container.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "Data-plane URL of the container.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Container name (lowercase, 3-63 chars).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"storage_account_name": schema.StringAttribute{
				Required: true, Description: "Parent Storage Account name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"container_access_type": schema.StringAttribute{
				Optional: true, Computed: true,
				Description: "`blob`, `container`, or `private`. Defaults to `private`.",
			},
			"metadata": schema.MapAttribute{
				Optional: true, ElementType: types.StringType,
				Description: "Container metadata.",
			},
			"has_immutability_policy": schema.BoolAttribute{
				Computed: true, Description: "Whether an immutability policy is set. Always `false` in this simulator.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
			"has_legal_hold": schema.BoolAttribute{
				Computed: true, Description: "Whether a legal hold is set. Always `false` in this simulator.",
				PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func (r *StorageContainerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *StorageContainerResource) applyComputed(plan *StorageContainerModel) {
	plan.ID = types.StringValue(fmt.Sprintf(
		"https://%s.blob.core.windows.net/%s",
		plan.StorageAccountName.ValueString(), plan.Name.ValueString(),
	))
	if plan.ContainerAccessType.IsNull() || plan.ContainerAccessType.IsUnknown() {
		plan.ContainerAccessType = types.StringValue("private")
	}
	plan.HasImmutabilityPolicy = types.BoolValue(false)
	plan.HasLegalHold = types.BoolValue(false)
}

func (r *StorageContainerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageContainerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageContainerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageContainerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageContainerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageContainerModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state StorageContainerModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageContainerResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
