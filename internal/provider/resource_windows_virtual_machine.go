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

var _ resource.Resource = &WindowsVirtualMachineResource{}
var _ resource.ResourceWithConfigure = &WindowsVirtualMachineResource{}

type WindowsVirtualMachineResource struct {
	subscriptionID string
}

type WindowsVirtualMachineModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	ResourceGroupName    types.String `tfsdk:"resource_group_name"`
	Location             types.String `tfsdk:"location"`
	Size                 types.String `tfsdk:"size"`
	AdminUsername        types.String `tfsdk:"admin_username"`
	AdminPassword        types.String `tfsdk:"admin_password"`
	ComputerName         types.String `tfsdk:"computer_name"`
	NetworkInterfaceIDs  types.List   `tfsdk:"network_interface_ids"`
	OSDisk               types.List   `tfsdk:"os_disk"`
	SourceImageReference types.List   `tfsdk:"source_image_reference"`
	LicenseType          types.String `tfsdk:"license_type"`
	HotpatchingEnabled   types.Bool   `tfsdk:"hotpatching_enabled"`
	PatchMode            types.String `tfsdk:"patch_mode"`
	Zone                 types.String `tfsdk:"zone"`
	Tags                 types.Map    `tfsdk:"tags"`
}

func NewWindowsVirtualMachineResource() resource.Resource {
	return &WindowsVirtualMachineResource{}
}

func (r *WindowsVirtualMachineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_windows_virtual_machine"
}

func (r *WindowsVirtualMachineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Windows Virtual Machine.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Virtual Machine ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "VM name.",
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
			"size":           schema.StringAttribute{Required: true, Description: "VM size (e.g. `Standard_D2s_v5`)."},
			"admin_username": schema.StringAttribute{Required: true, Description: "Local administrator username."},
			"admin_password": schema.StringAttribute{Required: true, Sensitive: true, Description: "Local administrator password."},
			"computer_name":  schema.StringAttribute{Optional: true, Description: "Computer name (defaults to VM name when omitted)."},
			"network_interface_ids": schema.ListAttribute{
				Required: true, ElementType: types.StringType,
				Description: "NIC IDs to attach. The first is treated as primary.",
			},
			"license_type":        schema.StringAttribute{Optional: true, Description: "`None`, `Windows_Client`, or `Windows_Server`."},
			"hotpatching_enabled": schema.BoolAttribute{Optional: true, Description: "Enable hotpatching (Server 2022 Datacenter Azure Edition only)."},
			"patch_mode":          schema.StringAttribute{Optional: true, Description: "`Manual`, `AutomaticByOS`, or `AutomaticByPlatform`."},
			"zone": schema.StringAttribute{
				Optional: true, Description: "Availability zone.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType, Description: "Tags.",
			},
		},
		Blocks: map[string]schema.Block{
			"os_disk": schema.ListNestedBlock{
				Description: "OS Disk configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"caching":              schema.StringAttribute{Required: true, Description: "`None`, `ReadOnly`, or `ReadWrite`."},
						"storage_account_type": schema.StringAttribute{Required: true, Description: "`Standard_LRS`, `StandardSSD_LRS`, `Premium_LRS`, etc."},
						"disk_size_gb":         schema.Int64Attribute{Optional: true, Description: "Disk size in GB."},
						"name":                 schema.StringAttribute{Optional: true, Description: "Optional disk name."},
					},
				},
			},
			"source_image_reference": schema.ListNestedBlock{
				Description: "Source image reference.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"publisher": schema.StringAttribute{Required: true, Description: "e.g. `MicrosoftWindowsServer`."},
						"offer":     schema.StringAttribute{Required: true, Description: "e.g. `WindowsServer`."},
						"sku":       schema.StringAttribute{Required: true, Description: "e.g. `2022-datacenter-azure-edition`."},
						"version":   schema.StringAttribute{Required: true, Description: "e.g. `latest`."},
					},
				},
			},
		},
	}
}

func (r *WindowsVirtualMachineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *WindowsVirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WindowsVirtualMachineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/virtualMachines/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsVirtualMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WindowsVirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WindowsVirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsVirtualMachineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state WindowsVirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WindowsVirtualMachineResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
