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

var _ resource.Resource = &APIManagementResource{}
var _ resource.ResourceWithConfigure = &APIManagementResource{}

type APIManagementResource struct{ subscriptionID string }

type APIManagementModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ResourceGroupName   types.String `tfsdk:"resource_group_name"`
	Location            types.String `tfsdk:"location"`
	PublisherName       types.String `tfsdk:"publisher_name"`
	PublisherEmail      types.String `tfsdk:"publisher_email"`
	SKUName             types.String `tfsdk:"sku_name"`
	Zones               types.List   `tfsdk:"zones"`
	VirtualNetworkType  types.String `tfsdk:"virtual_network_type"`
	GatewayURL          types.String `tfsdk:"gateway_url"`
	PortalURL           types.String `tfsdk:"portal_url"`
	ManagementAPIURL    types.String `tfsdk:"management_api_url"`
	DeveloperPortalURL  types.String `tfsdk:"developer_portal_url"`
	PublicIPAddresses   types.List   `tfsdk:"public_ip_addresses"`
	Tags                types.Map    `tfsdk:"tags"`
}

func NewAPIManagementResource() resource.Resource { return &APIManagementResource{} }

func (r *APIManagementResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_management"
}

func (r *APIManagementResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure API Management instance.",
		Attributes: map[string]schema.Attribute{
			"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "APIM ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "APIM name (globally unique)."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":             schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"publisher_name":       schema.StringAttribute{Required: true, Description: "Publisher (organization) name."},
			"publisher_email":      schema.StringAttribute{Required: true, Description: "Publisher email."},
			"sku_name":             schema.StringAttribute{Required: true, Description: "SKU (e.g. `Developer_1`, `Standard_2`, `Premium_4`, `Consumption_0`, `BasicV2_1`, `StandardV2_1`)."},
			"zones":                schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Availability zones."},
			"virtual_network_type": schema.StringAttribute{Optional: true, Description: "`None`, `External`, or `Internal`."},
			"gateway_url": schema.StringAttribute{
				Computed: true, Description: "Simulated gateway URL.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"portal_url": schema.StringAttribute{
				Computed: true, Description: "Simulated legacy portal URL.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"management_api_url": schema.StringAttribute{
				Computed: true, Description: "Simulated management API URL.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"developer_portal_url": schema.StringAttribute{
				Computed: true, Description: "Simulated developer portal URL.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"public_ip_addresses": schema.ListAttribute{
				Computed: true, ElementType: types.StringType, Description: "Simulated public IPs.",
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
	}
}

func (r *APIManagementResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *APIManagementResource) applyComputed(plan *APIManagementModel) {
	name := plan.Name.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), name,
	))
	plan.GatewayURL = types.StringValue(fmt.Sprintf("https://%s.azure-api.net", name))
	plan.PortalURL = types.StringValue(fmt.Sprintf("https://%s.portal.azure-api.net", name))
	plan.ManagementAPIURL = types.StringValue(fmt.Sprintf("https://%s.management.azure-api.net", name))
	plan.DeveloperPortalURL = types.StringValue(fmt.Sprintf("https://%s.developer.azure-api.net", name))
	plan.PublicIPAddresses = types.ListValueMust(types.StringType, nil)
}

func (r *APIManagementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan APIManagementModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *APIManagementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state APIManagementModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *APIManagementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan APIManagementModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state APIManagementModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.GatewayURL = state.GatewayURL
	plan.PortalURL = state.PortalURL
	plan.ManagementAPIURL = state.ManagementAPIURL
	plan.DeveloperPortalURL = state.DeveloperPortalURL
	plan.PublicIPAddresses = state.PublicIPAddresses
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *APIManagementResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
