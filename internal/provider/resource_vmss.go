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

// VMSS — Linux and Windows variants share the same shape.

type VMSSModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ResourceGroupName   types.String `tfsdk:"resource_group_name"`
	Location            types.String `tfsdk:"location"`
	SKU                 types.String `tfsdk:"sku"`
	Instances           types.Int64  `tfsdk:"instances"`
	AdminUsername       types.String `tfsdk:"admin_username"`
	AdminPassword       types.String `tfsdk:"admin_password"`
	NetworkInterface    types.List   `tfsdk:"network_interface"`
	OSDisk              types.List   `tfsdk:"os_disk"`
	SourceImageReference types.List  `tfsdk:"source_image_reference"`
	Zones               types.List   `tfsdk:"zones"`
	UpgradeMode         types.String `tfsdk:"upgrade_mode"`
	Overprovision       types.Bool   `tfsdk:"overprovision"`
	Tags                types.Map    `tfsdk:"tags"`
}

func vmssAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":   schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "VMSS ID."},
		"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "VMSS name."},
		"resource_group_name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Resource Group."},
		"location":      schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Azure region."},
		"sku":           schema.StringAttribute{Required: true, Description: "VM SKU (e.g. `Standard_D2s_v5`)."},
		"instances":     schema.Int64Attribute{Required: true, Description: "Initial instance count."},
		"admin_username": schema.StringAttribute{Required: true, Description: "Admin username."},
		"admin_password": schema.StringAttribute{Optional: true, Sensitive: true, Description: "Admin password (Linux requires this OR `admin_ssh_key`; this simulator accepts password only)."},
		"zones":         schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "Availability zones."},
		"upgrade_mode":  schema.StringAttribute{Optional: true, Description: "`Manual`, `Automatic`, or `Rolling`."},
		"overprovision": schema.BoolAttribute{Optional: true, Description: "Whether to overprovision VMs."},
		"tags":          schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
	}
}

func vmssBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"network_interface": schema.ListNestedBlock{
			Description: "NIC profile.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"name":    schema.StringAttribute{Required: true, Description: "NIC profile name."},
					"primary": schema.BoolAttribute{Optional: true, Description: "Primary NIC."},
				},
				Blocks: map[string]schema.Block{
					"ip_configuration": schema.ListNestedBlock{
						Description: "IP configuration.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name":      schema.StringAttribute{Required: true, Description: "Configuration name."},
								"subnet_id": schema.StringAttribute{Required: true, Description: "Subnet ID."},
								"primary":   schema.BoolAttribute{Optional: true, Description: "Primary IP config."},
								"load_balancer_backend_address_pool_ids": schema.ListAttribute{Optional: true, ElementType: types.StringType, Description: "LB backend pool IDs."},
							},
						},
					},
				},
			},
		},
		"os_disk": schema.ListNestedBlock{
			Description: "OS Disk.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"caching":              schema.StringAttribute{Required: true, Description: "Caching mode."},
					"storage_account_type": schema.StringAttribute{Required: true, Description: "Storage type."},
					"disk_size_gb":         schema.Int64Attribute{Optional: true, Description: "Disk size in GB."},
				},
			},
		},
		"source_image_reference": schema.ListNestedBlock{
			Description: "Image reference.",
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"publisher": schema.StringAttribute{Required: true, Description: "Publisher."},
					"offer":     schema.StringAttribute{Required: true, Description: "Offer."},
					"sku":       schema.StringAttribute{Required: true, Description: "SKU."},
					"version":   schema.StringAttribute{Required: true, Description: "Version."},
				},
			},
		},
	}
}

func vmssApplyComputed(plan *VMSSModel, subscriptionID string) {
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachineScaleSets/%s",
		subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
}

// --- Linux VMSS ---

var _ resource.Resource = &LinuxVMSSResource{}
var _ resource.ResourceWithConfigure = &LinuxVMSSResource{}

type LinuxVMSSResource struct{ subscriptionID string }

func NewLinuxVMSSResource() resource.Resource { return &LinuxVMSSResource{} }

func (r *LinuxVMSSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_linux_virtual_machine_scale_set"
}

func (r *LinuxVMSSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Linux Virtual Machine Scale Set.",
		Attributes:  vmssAttributes(),
		Blocks:      vmssBlocks(),
	}
}

func (r *LinuxVMSSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *LinuxVMSSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VMSSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	vmssApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LinuxVMSSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VMSSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *LinuxVMSSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VMSSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	vmssApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *LinuxVMSSResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Windows VMSS ---

var _ resource.Resource = &WindowsVMSSResource{}
var _ resource.ResourceWithConfigure = &WindowsVMSSResource{}

type WindowsVMSSResource struct{ subscriptionID string }

func NewWindowsVMSSResource() resource.Resource { return &WindowsVMSSResource{} }

func (r *WindowsVMSSResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_windows_virtual_machine_scale_set"
}

func (r *WindowsVMSSResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Windows Virtual Machine Scale Set.",
		Attributes:  vmssAttributes(),
		Blocks:      vmssBlocks(),
	}
}

func (r *WindowsVMSSResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *WindowsVMSSResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VMSSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	vmssApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsVMSSResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VMSSModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WindowsVMSSResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VMSSModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	vmssApplyComputed(&plan, r.subscriptionID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsVMSSResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
