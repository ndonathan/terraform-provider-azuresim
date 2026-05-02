package provider

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &PublicIPResource{}
var _ resource.ResourceWithConfigure = &PublicIPResource{}

type PublicIPResource struct {
	subscriptionID string
}

type PublicIPModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	ResourceGroupName    types.String `tfsdk:"resource_group_name"`
	Location             types.String `tfsdk:"location"`
	AllocationMethod     types.String `tfsdk:"allocation_method"`
	SKU                  types.String `tfsdk:"sku"`
	SKUTier              types.String `tfsdk:"sku_tier"`
	IPVersion            types.String `tfsdk:"ip_version"`
	DomainNameLabel      types.String `tfsdk:"domain_name_label"`
	IdleTimeoutInMinutes types.Int64  `tfsdk:"idle_timeout_in_minutes"`
	ReverseFQDN          types.String `tfsdk:"reverse_fqdn"`
	Zones                types.List   `tfsdk:"zones"`
	IPAddress            types.String `tfsdk:"ip_address"`
	FQDN                 types.String `tfsdk:"fqdn"`
	Tags                 types.Map    `tfsdk:"tags"`
}

func NewPublicIPResource() resource.Resource {
	return &PublicIPResource{}
}

func (r *PublicIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

func (r *PublicIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Public IP Address.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The simulated Public IP ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Public IP.",
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
			"allocation_method": schema.StringAttribute{
				Required:    true,
				Description: "Allocation method (`Static` or `Dynamic`).",
			},
			"sku": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "SKU (`Basic` or `Standard`). Defaults to `Basic`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"sku_tier": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "SKU tier (`Regional` or `Global`). Defaults to `Regional`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP version (`IPv4` or `IPv6`). Defaults to `IPv4`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain_name_label": schema.StringAttribute{
				Optional:    true,
				Description: "DNS label. When set, `fqdn` becomes `<label>.<location>.cloudapp.azure.com`.",
			},
			"idle_timeout_in_minutes": schema.Int64Attribute{
				Optional:    true,
				Description: "Idle timeout in minutes (4-30).",
			},
			"reverse_fqdn": schema.StringAttribute{
				Optional:    true,
				Description: "Reverse FQDN to assign.",
			},
			"zones": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Availability zones (e.g. `[\"1\", \"2\", \"3\"]`).",
			},
			"ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "The simulated IP address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"fqdn": schema.StringAttribute{
				Computed:    true,
				Description: "Fully qualified domain name. Empty unless `domain_name_label` is set.",
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

func (r *PublicIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

// simulatedIPv4 derives a deterministic dummy IP in the TEST-NET-3 range
// (203.0.113.0/24, RFC 5737) from the resource name + RG so collisions are unlikely.
func simulatedIPv4(name, resourceGroup string) string {
	sum := sha1.Sum([]byte(resourceGroup + "/" + name))
	return fmt.Sprintf("203.0.113.%d", sum[0])
}

func simulatedIPv6(name, resourceGroup string) string {
	sum := sha1.Sum([]byte(resourceGroup + "/" + name))
	return fmt.Sprintf("2001:db8::%x:%x", sum[0], sum[1])
}

func (r *PublicIPResource) applyComputed(plan *PublicIPModel) {
	if plan.SKU.IsNull() || plan.SKU.IsUnknown() {
		plan.SKU = types.StringValue("Basic")
	}
	if plan.SKUTier.IsNull() || plan.SKUTier.IsUnknown() {
		plan.SKUTier = types.StringValue("Regional")
	}
	if plan.IPVersion.IsNull() || plan.IPVersion.IsUnknown() {
		plan.IPVersion = types.StringValue("IPv4")
	}

	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()

	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/publicIPAddresses/%s",
		r.subscriptionID, rg, name,
	))

	if strings.EqualFold(plan.IPVersion.ValueString(), "IPv6") {
		plan.IPAddress = types.StringValue(simulatedIPv6(name, rg))
	} else {
		plan.IPAddress = types.StringValue(simulatedIPv4(name, rg))
	}

	if !plan.DomainNameLabel.IsNull() && !plan.DomainNameLabel.IsUnknown() && plan.DomainNameLabel.ValueString() != "" {
		plan.FQDN = types.StringValue(fmt.Sprintf(
			"%s.%s.cloudapp.azure.com",
			plan.DomainNameLabel.ValueString(), plan.Location.ValueString(),
		))
	} else {
		plan.FQDN = types.StringValue("")
	}
}

func (r *PublicIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PublicIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.applyComputed(&plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PublicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PublicIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PublicIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PublicIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state PublicIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	r.applyComputed(&plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PublicIPResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: simulated resource
}
