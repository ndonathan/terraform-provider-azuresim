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

var _ resource.Resource = &NetworkSecurityGroupResource{}
var _ resource.ResourceWithConfigure = &NetworkSecurityGroupResource{}

type NetworkSecurityGroupResource struct {
	subscriptionID string
}

type NetworkSecurityGroupModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	Location          types.String `tfsdk:"location"`
	SecurityRule      types.List   `tfsdk:"security_rule"`
	Tags              types.Map    `tfsdk:"tags"`
}

type SecurityRuleModel struct {
	Name                                 types.String `tfsdk:"name"`
	Description                          types.String `tfsdk:"description"`
	Protocol                             types.String `tfsdk:"protocol"`
	SourcePortRange                      types.String `tfsdk:"source_port_range"`
	SourcePortRanges                     types.List   `tfsdk:"source_port_ranges"`
	DestinationPortRange                 types.String `tfsdk:"destination_port_range"`
	DestinationPortRanges                types.List   `tfsdk:"destination_port_ranges"`
	SourceAddressPrefix                  types.String `tfsdk:"source_address_prefix"`
	SourceAddressPrefixes                types.List   `tfsdk:"source_address_prefixes"`
	DestinationAddressPrefix             types.String `tfsdk:"destination_address_prefix"`
	DestinationAddressPrefixes           types.List   `tfsdk:"destination_address_prefixes"`
	Access                               types.String `tfsdk:"access"`
	Priority                             types.Int64  `tfsdk:"priority"`
	Direction                            types.String `tfsdk:"direction"`
}

func NewNetworkSecurityGroupResource() resource.Resource {
	return &NetworkSecurityGroupResource{}
}

func (r *NetworkSecurityGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_security_group"
}

func (r *NetworkSecurityGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Network Security Group.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The simulated Network Security Group ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Network Security Group.",
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
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to assign to the resource.",
			},
		},
		Blocks: map[string]schema.Block{
			"security_rule": schema.ListNestedBlock{
				Description: "An inline security rule.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of the security rule.",
						},
						"description": schema.StringAttribute{
							Optional:    true,
							Description: "A description for the rule.",
						},
						"protocol": schema.StringAttribute{
							Required:    true,
							Description: "Network protocol (Tcp, Udp, Icmp, Esp, Ah, *).",
						},
						"source_port_range": schema.StringAttribute{
							Optional:    true,
							Description: "Source port or range. Use `*` for any. Mutually exclusive with `source_port_ranges`.",
						},
						"source_port_ranges": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "List of source ports or ranges. Mutually exclusive with `source_port_range`.",
						},
						"destination_port_range": schema.StringAttribute{
							Optional:    true,
							Description: "Destination port or range. Mutually exclusive with `destination_port_ranges`.",
						},
						"destination_port_ranges": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "List of destination ports or ranges. Mutually exclusive with `destination_port_range`.",
						},
						"source_address_prefix": schema.StringAttribute{
							Optional:    true,
							Description: "CIDR or service tag (e.g. `*`, `VirtualNetwork`, `Internet`). Mutually exclusive with `source_address_prefixes`.",
						},
						"source_address_prefixes": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "List of source CIDRs. Mutually exclusive with `source_address_prefix`.",
						},
						"destination_address_prefix": schema.StringAttribute{
							Optional:    true,
							Description: "CIDR or service tag. Mutually exclusive with `destination_address_prefixes`.",
						},
						"destination_address_prefixes": schema.ListAttribute{
							Optional:    true,
							ElementType: types.StringType,
							Description: "List of destination CIDRs. Mutually exclusive with `destination_address_prefix`.",
						},
						"access": schema.StringAttribute{
							Required:    true,
							Description: "`Allow` or `Deny`.",
						},
						"priority": schema.Int64Attribute{
							Required:    true,
							Description: "Rule priority (100-4096). Lower numbers evaluate first.",
						},
						"direction": schema.StringAttribute{
							Required:    true,
							Description: "`Inbound` or `Outbound`.",
						},
					},
				},
			},
		},
	}
}

func (r *NetworkSecurityGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *NetworkSecurityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkSecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkSecurityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkSecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkSecurityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkSecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state NetworkSecurityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkSecurityGroupResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: simulated resource
}
