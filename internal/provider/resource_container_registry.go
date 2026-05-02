package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ContainerRegistryResource{}
var _ resource.ResourceWithConfigure = &ContainerRegistryResource{}

type ContainerRegistryResource struct {
	subscriptionID string
}

type ContainerRegistryModel struct {
	ID                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	ResourceGroupName           types.String `tfsdk:"resource_group_name"`
	Location                    types.String `tfsdk:"location"`
	SKU                         types.String `tfsdk:"sku"`
	AdminEnabled                types.Bool   `tfsdk:"admin_enabled"`
	PublicNetworkAccessEnabled  types.Bool   `tfsdk:"public_network_access_enabled"`
	ZoneRedundancyEnabled       types.Bool   `tfsdk:"zone_redundancy_enabled"`
	AnonymousPullEnabled        types.Bool   `tfsdk:"anonymous_pull_enabled"`
	DataEndpointEnabled         types.Bool   `tfsdk:"data_endpoint_enabled"`
	LoginServer                 types.String `tfsdk:"login_server"`
	AdminUsername               types.String `tfsdk:"admin_username"`
	AdminPassword               types.String `tfsdk:"admin_password"`
	Tags                        types.Map    `tfsdk:"tags"`
}

func NewContainerRegistryResource() resource.Resource { return &ContainerRegistryResource{} }

func (r *ContainerRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_registry"
}

func (r *ContainerRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Container Registry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Registry ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Registry name (5-50 alphanumeric, globally unique in real Azure).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Resource Group.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"location": schema.StringAttribute{
				Required: true, Description: "Azure region.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"sku":                            schema.StringAttribute{Required: true, Description: "`Basic`, `Standard`, or `Premium`."},
			"admin_enabled":                  schema.BoolAttribute{Optional: true, Description: "Enable admin user."},
			"public_network_access_enabled":  schema.BoolAttribute{Optional: true, Description: "Allow public network access."},
			"zone_redundancy_enabled":        schema.BoolAttribute{Optional: true, Description: "Enable zone redundancy (Premium only)."},
			"anonymous_pull_enabled":         schema.BoolAttribute{Optional: true, Description: "Allow anonymous pulls."},
			"data_endpoint_enabled":          schema.BoolAttribute{Optional: true, Description: "Enable dedicated data endpoints."},
			"login_server": schema.StringAttribute{
				Computed: true, Description: "Simulated login server (`<name>.azurecr.io`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"admin_username": schema.StringAttribute{
				Computed: true, Description: "Admin username (always equals `name` when admin is enabled).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"admin_password": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated admin password.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
	}
}

func (r *ContainerRegistryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ContainerRegistryResource) applyComputed(plan *ContainerRegistryModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerRegistry/registries/%s",
		r.subscriptionID, rg, name,
	))
	plan.LoginServer = types.StringValue(fmt.Sprintf("%s.azurecr.io", name))
	plan.AdminUsername = types.StringValue(name)
	sum := sha1.Sum([]byte("acr-admin/" + rg + "/" + name))
	plan.AdminPassword = types.StringValue(base64.StdEncoding.EncodeToString(sum[:]))
}

func (r *ContainerRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ContainerRegistryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ContainerRegistryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ContainerRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ContainerRegistryModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ContainerRegistryModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.LoginServer = state.LoginServer
	plan.AdminUsername = state.AdminUsername
	plan.AdminPassword = state.AdminPassword
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ContainerRegistryResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
