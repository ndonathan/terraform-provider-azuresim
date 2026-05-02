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

var _ resource.Resource = &ApplicationGatewayResource{}
var _ resource.ResourceWithConfigure = &ApplicationGatewayResource{}

type ApplicationGatewayResource struct{ subscriptionID string }

type ApplicationGatewayModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	Location          types.String `tfsdk:"location"`
	Zones             types.List   `tfsdk:"zones"`
	EnableHTTP2       types.Bool   `tfsdk:"enable_http2"`
	FipsEnabled       types.Bool   `tfsdk:"fips_enabled"`
	FirewallPolicyID  types.String `tfsdk:"firewall_policy_id"`
	SKU               types.List   `tfsdk:"sku"`
	GatewayIPConfig   types.List   `tfsdk:"gateway_ip_configuration"`
	FrontendPort      types.List   `tfsdk:"frontend_port"`
	FrontendIPConfig  types.List   `tfsdk:"frontend_ip_configuration"`
	BackendAddressPool types.List  `tfsdk:"backend_address_pool"`
	BackendHTTPSettings types.List `tfsdk:"backend_http_settings"`
	HTTPListener      types.List   `tfsdk:"http_listener"`
	RequestRoutingRule types.List  `tfsdk:"request_routing_rule"`
	Tags              types.Map    `tfsdk:"tags"`
}

func NewApplicationGatewayResource() resource.Resource { return &ApplicationGatewayResource{} }

func (r *ApplicationGatewayResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_gateway"
}

func (r *ApplicationGatewayResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	nameOnly := func(desc string) schema.NestedBlockObject {
		return schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{Required: true, Description: desc},
			},
		}
	}

	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Application Gateway. Schema captures the major top-level blocks; detailed block attributes are intentionally minimal in this simulator.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Application Gateway ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Gateway name."},
			"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
			"location":           schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
			"zones":              schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Availability zones."},
			"enable_http2":       schema.BoolAttribute{Optional: true, Description: "Enable HTTP/2."},
			"fips_enabled":       schema.BoolAttribute{Optional: true, Description: "Enable FIPS mode."},
			"firewall_policy_id": schema.StringAttribute{Optional: true, Description: "WAF policy ID."},
			"tags":               schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"sku": schema.ListNestedBlock{
				Description: "SKU.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":     schema.StringAttribute{Required: true, Description: "SKU name (e.g. `Standard_v2`, `WAF_v2`)."},
						"tier":     schema.StringAttribute{Required: true, Description: "Tier."},
						"capacity": schema.Int64Attribute{Optional: true, Description: "Fixed capacity."},
					},
				},
			},
			"gateway_ip_configuration": schema.ListNestedBlock{
				Description: "Gateway IP configuration (subnet binding).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":      schema.StringAttribute{Required: true, Description: "Configuration name."},
						"subnet_id": schema.StringAttribute{Required: true, Description: "Subnet ID."},
					},
				},
			},
			"frontend_port": schema.ListNestedBlock{
				Description: "Frontend port.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{Required: true, Description: "Port name."},
						"port": schema.Int64Attribute{Required: true, Description: "Port number."},
					},
				},
			},
			"frontend_ip_configuration": schema.ListNestedBlock{
				Description: "Frontend IP configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name":                          schema.StringAttribute{Required: true, Description: "Configuration name."},
						"public_ip_address_id":          schema.StringAttribute{Optional: true, Description: "Public IP ID."},
						"subnet_id":                     schema.StringAttribute{Optional: true, Description: "Subnet ID for private LB."},
						"private_ip_address":            schema.StringAttribute{Optional: true, Description: "Static private IP."},
						"private_ip_address_allocation": schema.StringAttribute{Optional: true, Description: "`Static` or `Dynamic`."},
					},
				},
			},
			"backend_address_pool":  schema.ListNestedBlock{Description: "Backend pool.", NestedObject: nameOnly("Backend pool name.")},
			"backend_http_settings": schema.ListNestedBlock{Description: "Backend HTTP settings.", NestedObject: nameOnly("Settings name.")},
			"http_listener":         schema.ListNestedBlock{Description: "HTTP listener.", NestedObject: nameOnly("Listener name.")},
			"request_routing_rule":  schema.ListNestedBlock{Description: "Request routing rule.", NestedObject: nameOnly("Rule name.")},
		},
	}
}

func (r *ApplicationGatewayResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ApplicationGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ApplicationGatewayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApplicationGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ApplicationGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ApplicationGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationGatewayModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ApplicationGatewayModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ApplicationGatewayResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
