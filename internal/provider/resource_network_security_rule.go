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

var _ resource.Resource = &NetworkSecurityRuleResource{}
var _ resource.ResourceWithConfigure = &NetworkSecurityRuleResource{}

type NetworkSecurityRuleResource struct {
	subscriptionID string
}

type NetworkSecurityRuleModel struct {
	ID                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	ResourceGroupName                  types.String `tfsdk:"resource_group_name"`
	NetworkSecurityGroupName           types.String `tfsdk:"network_security_group_name"`
	Description                        types.String `tfsdk:"description"`
	Protocol                           types.String `tfsdk:"protocol"`
	SourcePortRange                    types.String `tfsdk:"source_port_range"`
	SourcePortRanges                   types.List   `tfsdk:"source_port_ranges"`
	DestinationPortRange               types.String `tfsdk:"destination_port_range"`
	DestinationPortRanges              types.List   `tfsdk:"destination_port_ranges"`
	SourceAddressPrefix                types.String `tfsdk:"source_address_prefix"`
	SourceAddressPrefixes              types.List   `tfsdk:"source_address_prefixes"`
	DestinationAddressPrefix           types.String `tfsdk:"destination_address_prefix"`
	DestinationAddressPrefixes         types.List   `tfsdk:"destination_address_prefixes"`
	Access                             types.String `tfsdk:"access"`
	Priority                           types.Int64  `tfsdk:"priority"`
	Direction                          types.String `tfsdk:"direction"`
}

func NewNetworkSecurityRuleResource() resource.Resource {
	return &NetworkSecurityRuleResource{}
}

func (r *NetworkSecurityRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_security_rule"
}

func (r *NetworkSecurityRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A standalone Network Security Group rule. Use either inline `security_rule` blocks on `azuresim_network_security_group` or this resource — not both.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "Rule ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the rule.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Resource Group name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"network_security_group_name": schema.StringAttribute{
				Required: true, Description: "Parent NSG name.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"description":                 schema.StringAttribute{Optional: true, Description: "Description for the rule."},
			"protocol":                    schema.StringAttribute{Required: true, Description: "`Tcp`, `Udp`, `Icmp`, `Esp`, `Ah`, or `*`."},
			"source_port_range":           schema.StringAttribute{Optional: true, Description: "Source port or range. Mutually exclusive with `source_port_ranges`."},
			"source_port_ranges":          schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "List of source ports."},
			"destination_port_range":      schema.StringAttribute{Optional: true, Description: "Destination port or range."},
			"destination_port_ranges":     schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "List of destination ports."},
			"source_address_prefix":       schema.StringAttribute{Optional: true, Description: "Source CIDR or service tag."},
			"source_address_prefixes":     schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "List of source CIDRs."},
			"destination_address_prefix":  schema.StringAttribute{Optional: true, Description: "Destination CIDR or service tag."},
			"destination_address_prefixes": schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "List of destination CIDRs."},
			"access":                      schema.StringAttribute{Required: true, Description: "`Allow` or `Deny`."},
			"priority":                    schema.Int64Attribute{Required: true, Description: "Priority (100-4096)."},
			"direction":                   schema.StringAttribute{Required: true, Description: "`Inbound` or `Outbound`."},
		},
	}
}

func (r *NetworkSecurityRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *NetworkSecurityRuleResource) applyComputed(plan *NetworkSecurityRuleModel) {
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkSecurityGroups/%s/securityRules/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(),
		plan.NetworkSecurityGroupName.ValueString(), plan.Name.ValueString(),
	))
}

func (r *NetworkSecurityRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkSecurityRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkSecurityRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkSecurityRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkSecurityRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkSecurityRuleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkSecurityRuleResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
