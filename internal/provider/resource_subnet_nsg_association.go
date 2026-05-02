package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &SubnetNSGAssociationResource{}
var _ resource.ResourceWithConfigure = &SubnetNSGAssociationResource{}

type SubnetNSGAssociationResource struct {
	subscriptionID string
}

type SubnetNSGAssociationModel struct {
	ID                       types.String `tfsdk:"id"`
	SubnetID                 types.String `tfsdk:"subnet_id"`
	NetworkSecurityGroupID   types.String `tfsdk:"network_security_group_id"`
}

func NewSubnetNSGAssociationResource() resource.Resource {
	return &SubnetNSGAssociationResource{}
}

func (r *SubnetNSGAssociationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet_network_security_group_association"
}

func (r *SubnetNSGAssociationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Associates a Network Security Group with a Subnet. Mirrors the AzureRM convention of using the Subnet ID as the resource ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Subnet ID (used as the association's identifier).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"subnet_id": schema.StringAttribute{
				Required: true, Description: "The Subnet ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"network_security_group_id": schema.StringAttribute{
				Required: true, Description: "The Network Security Group ID.",
			},
		},
	}
}

func (r *SubnetNSGAssociationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *SubnetNSGAssociationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SubnetNSGAssociationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = plan.SubnetID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetNSGAssociationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SubnetNSGAssociationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SubnetNSGAssociationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SubnetNSGAssociationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = plan.SubnetID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SubnetNSGAssociationResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
