package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &StorageAccountResource{}
var _ resource.ResourceWithConfigure = &StorageAccountResource{}

type StorageAccountResource struct {
	subscriptionID string
}

type StorageAccountModel struct {
	ID                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	ResourceGroupName      types.String `tfsdk:"resource_group_name"`
	Location               types.String `tfsdk:"location"`
	AccountTier            types.String `tfsdk:"account_tier"`
	AccountReplicationType types.String `tfsdk:"account_replication_type"`
	AccountKind            types.String `tfsdk:"account_kind"`
	PrimaryBlobEndpoint    types.String `tfsdk:"primary_blob_endpoint"`
	PrimaryAccessKey       types.String `tfsdk:"primary_access_key"`
	Tags                   types.Map    `tfsdk:"tags"`
}

func NewStorageAccountResource() resource.Resource {
	return &StorageAccountResource{}
}

func (r *StorageAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_account"
}

func (r *StorageAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Storage Account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The simulated Storage Account ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Storage Account (must be globally unique in real Azure, 3-24 lowercase alphanumeric).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_group_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Resource Group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"location": schema.StringAttribute{
				Required:    true,
				Description: "The Azure region.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_tier": schema.StringAttribute{
				Required:    true,
				Description: "The tier of the storage account (Standard or Premium).",
			},
			"account_replication_type": schema.StringAttribute{
				Required:    true,
				Description: "The replication type (LRS, GRS, RAGRS, ZRS).",
			},
			"account_kind": schema.StringAttribute{
				Optional:    true,
				Description: "The kind of storage account (StorageV2, BlobStorage, etc.). Defaults to StorageV2.",
			},
			"primary_blob_endpoint": schema.StringAttribute{
				Computed:    true,
				Description: "The simulated primary blob endpoint URL.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"primary_access_key": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "A simulated primary access key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to assign to the resource.",
			},
		},
	}
}

func (r *StorageAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *StorageAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan StorageAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()

	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.PrimaryBlobEndpoint = types.StringValue(fmt.Sprintf("https://%s.blob.core.windows.net/", name))
	plan.PrimaryAccessKey = types.StringValue("c2ltdWxhdGVkLWFjY2Vzcy1rZXktMDAwMDAwMDAwMDAwMDAwMA==")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state StorageAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *StorageAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan StorageAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state StorageAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.PrimaryBlobEndpoint = state.PrimaryBlobEndpoint
	plan.PrimaryAccessKey = state.PrimaryAccessKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *StorageAccountResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: simulated resource
}
