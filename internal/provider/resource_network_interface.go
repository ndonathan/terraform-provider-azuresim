package provider

import (
	"context"
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &NetworkInterfaceResource{}
var _ resource.ResourceWithConfigure = &NetworkInterfaceResource{}

type NetworkInterfaceResource struct {
	subscriptionID string
}

type NetworkInterfaceModel struct {
	ID                           types.String `tfsdk:"id"`
	Name                         types.String `tfsdk:"name"`
	ResourceGroupName            types.String `tfsdk:"resource_group_name"`
	Location                     types.String `tfsdk:"location"`
	IPConfiguration              types.List   `tfsdk:"ip_configuration"`
	DNSServers                   types.List   `tfsdk:"dns_servers"`
	InternalDNSNameLabel         types.String `tfsdk:"internal_dns_name_label"`
	AcceleratedNetworkingEnabled types.Bool   `tfsdk:"accelerated_networking_enabled"`
	IPForwardingEnabled          types.Bool   `tfsdk:"ip_forwarding_enabled"`
	EdgeZone                     types.String `tfsdk:"edge_zone"`
	MACAddress                   types.String `tfsdk:"mac_address"`
	PrivateIPAddress             types.String `tfsdk:"private_ip_address"`
	PrivateIPAddresses           types.List   `tfsdk:"private_ip_addresses"`
	AppliedDNSServers            types.List   `tfsdk:"applied_dns_servers"`
	InternalDomainNameSuffix     types.String `tfsdk:"internal_domain_name_suffix"`
	Tags                         types.Map    `tfsdk:"tags"`
}

type IPConfigurationModel struct {
	Name                       types.String `tfsdk:"name"`
	SubnetID                   types.String `tfsdk:"subnet_id"`
	PrivateIPAddressAllocation types.String `tfsdk:"private_ip_address_allocation"`
	PrivateIPAddress           types.String `tfsdk:"private_ip_address"`
	PrivateIPAddressVersion    types.String `tfsdk:"private_ip_address_version"`
	PublicIPAddressID          types.String `tfsdk:"public_ip_address_id"`
	Primary                    types.Bool   `tfsdk:"primary"`
}

var ipConfigurationAttrTypes = map[string]attr.Type{
	"name":                          types.StringType,
	"subnet_id":                     types.StringType,
	"private_ip_address_allocation": types.StringType,
	"private_ip_address":            types.StringType,
	"private_ip_address_version":    types.StringType,
	"public_ip_address_id":          types.StringType,
	"primary":                       types.BoolType,
}

func NewNetworkInterfaceResource() resource.Resource {
	return &NetworkInterfaceResource{}
}

func (r *NetworkInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_interface"
}

func (r *NetworkInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Network Interface.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The simulated Network Interface ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Network Interface.",
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
			"dns_servers": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "List of DNS server IP addresses.",
			},
			"internal_dns_name_label": schema.StringAttribute{
				Optional:    true,
				Description: "Relative DNS name for this NIC (within the VNet).",
			},
			"accelerated_networking_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether accelerated networking is enabled. Defaults to `false`.",
			},
			"ip_forwarding_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether IP forwarding is enabled. Defaults to `false`.",
			},
			"edge_zone": schema.StringAttribute{
				Optional:    true,
				Description: "The Edge Zone within the Azure region.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mac_address": schema.StringAttribute{
				Computed:    true,
				Description: "Simulated MAC address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "Primary private IP address (mirrors the first `ip_configuration`).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_ip_addresses": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "All private IP addresses, in order of `ip_configuration` blocks.",
			},
			"applied_dns_servers": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "DNS servers actually applied to the NIC (echoes `dns_servers`).",
			},
			"internal_domain_name_suffix": schema.StringAttribute{
				Computed:    true,
				Description: "Internal DNS suffix.",
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
		Blocks: map[string]schema.Block{
			"ip_configuration": schema.ListNestedBlock{
				Description: "One or more IP configurations. At least one is required.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Name of the IP configuration.",
						},
						"subnet_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of the subnet (required for IPv4).",
						},
						"private_ip_address_allocation": schema.StringAttribute{
							Required:    true,
							Description: "`Static` or `Dynamic`.",
						},
						"private_ip_address": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The private IP address. Required when `private_ip_address_allocation` is `Static`; computed when `Dynamic`.",
						},
						"private_ip_address_version": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "`IPv4` or `IPv6`. Defaults to `IPv4`.",
						},
						"public_ip_address_id": schema.StringAttribute{
							Optional:    true,
							Description: "ID of an associated Public IP.",
						},
						"primary": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "Whether this is the primary IP configuration.",
						},
					},
				},
			},
		},
	}
}

func (r *NetworkInterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func simulatedMAC(name, resourceGroup string) string {
	sum := sha1.Sum([]byte(resourceGroup + "/" + name))
	// Use 00-15-5D (Microsoft Hyper-V OUI) as prefix for realism.
	return fmt.Sprintf("00-15-5D-%02X-%02X-%02X", sum[0], sum[1], sum[2])
}

func simulatedPrivateIPv4(nicName, configName string, index int) string {
	sum := sha1.Sum([]byte(nicName + "/" + configName))
	// Avoid the first 4 reserved addresses in an Azure subnet by adding 4 + index.
	octet := int(sum[0])%250 + 4 + index
	if octet > 254 {
		octet = 254 - (index % 10)
	}
	return fmt.Sprintf("10.0.0.%d", octet)
}

func simulatedPrivateIPv6(nicName, configName string, index int) string {
	sum := sha1.Sum([]byte(nicName + "/" + configName))
	return fmt.Sprintf("fd00::%x:%x:%x", sum[0], sum[1], byte(index))
}

func (r *NetworkInterfaceResource) applyComputed(ctx context.Context, plan *NetworkInterfaceModel) error {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()

	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkInterfaces/%s",
		r.subscriptionID, rg, name,
	))
	plan.MACAddress = types.StringValue(simulatedMAC(name, rg))
	plan.InternalDomainNameSuffix = types.StringValue("internal.cloudapp.net")

	if plan.AcceleratedNetworkingEnabled.IsNull() || plan.AcceleratedNetworkingEnabled.IsUnknown() {
		plan.AcceleratedNetworkingEnabled = types.BoolValue(false)
	}
	if plan.IPForwardingEnabled.IsNull() || plan.IPForwardingEnabled.IsUnknown() {
		plan.IPForwardingEnabled = types.BoolValue(false)
	}

	// Mirror dns_servers into applied_dns_servers (or empty list if unset).
	if plan.DNSServers.IsNull() || plan.DNSServers.IsUnknown() {
		plan.AppliedDNSServers = types.ListValueMust(types.StringType, []attr.Value{})
	} else {
		plan.AppliedDNSServers = plan.DNSServers
	}

	// Resolve ip_configuration: assign computed private_ip_address / primary,
	// then surface the primary address up to the top level.
	var configs []IPConfigurationModel
	if !plan.IPConfiguration.IsNull() && !plan.IPConfiguration.IsUnknown() {
		diags := plan.IPConfiguration.ElementsAs(ctx, &configs, false)
		if diags.HasError() {
			return fmt.Errorf("failed to read ip_configuration: %v", diags.Errors())
		}
	}

	resolvedElems := make([]attr.Value, 0, len(configs))
	addresses := make([]attr.Value, 0, len(configs))
	var primaryAddr string

	for i, cfg := range configs {
		ipVersion := cfg.PrivateIPAddressVersion.ValueString()
		if ipVersion == "" {
			ipVersion = "IPv4"
		}

		// Determine the private IP for this config.
		var privateIP string
		if !cfg.PrivateIPAddress.IsNull() && !cfg.PrivateIPAddress.IsUnknown() && cfg.PrivateIPAddress.ValueString() != "" {
			privateIP = cfg.PrivateIPAddress.ValueString()
		} else if strings.EqualFold(ipVersion, "IPv6") {
			privateIP = simulatedPrivateIPv6(name, cfg.Name.ValueString(), i)
		} else {
			privateIP = simulatedPrivateIPv4(name, cfg.Name.ValueString(), i)
		}
		cfg.PrivateIPAddress = types.StringValue(privateIP)

		// Default `primary` to true for the first config, false for the rest.
		if cfg.Primary.IsNull() || cfg.Primary.IsUnknown() {
			cfg.Primary = types.BoolValue(i == 0)
		}

		// Default `private_ip_address_version` if unset.
		if cfg.PrivateIPAddressVersion.IsNull() || cfg.PrivateIPAddressVersion.IsUnknown() {
			cfg.PrivateIPAddressVersion = types.StringValue("IPv4")
		}

		// Default `public_ip_address_id` to empty string when unset (keeps state shape stable).
		if cfg.PublicIPAddressID.IsNull() || cfg.PublicIPAddressID.IsUnknown() {
			cfg.PublicIPAddressID = types.StringNull()
		}

		obj, diags := types.ObjectValue(ipConfigurationAttrTypes, map[string]attr.Value{
			"name":                          cfg.Name,
			"subnet_id":                     cfg.SubnetID,
			"private_ip_address_allocation": cfg.PrivateIPAddressAllocation,
			"private_ip_address":            cfg.PrivateIPAddress,
			"private_ip_address_version":    cfg.PrivateIPAddressVersion,
			"public_ip_address_id":          cfg.PublicIPAddressID,
			"primary":                       cfg.Primary,
		})
		if diags.HasError() {
			return fmt.Errorf("failed to build ip_configuration object: %v", diags.Errors())
		}
		resolvedElems = append(resolvedElems, obj)
		addresses = append(addresses, types.StringValue(privateIP))

		if cfg.Primary.ValueBool() && primaryAddr == "" {
			primaryAddr = privateIP
		}
	}
	if primaryAddr == "" && len(configs) > 0 {
		primaryAddr = configs[0].PrivateIPAddress.ValueString()
	}

	listType := types.ObjectType{AttrTypes: ipConfigurationAttrTypes}
	plan.IPConfiguration = types.ListValueMust(listType, resolvedElems)
	plan.PrivateIPAddresses = types.ListValueMust(types.StringType, addresses)
	plan.PrivateIPAddress = types.StringValue(primaryAddr)

	return nil
}

func (r *NetworkInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.applyComputed(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Failed to materialize NIC", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworkInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.applyComputed(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Failed to materialize NIC", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworkInterfaceResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: simulated resource
}
