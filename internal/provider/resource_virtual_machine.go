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

var _ resource.Resource = &VirtualMachineResource{}
var _ resource.ResourceWithConfigure = &VirtualMachineResource{}

type VirtualMachineResource struct {
	subscriptionID string
}

type VirtualMachineModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	ResourceGroupName    types.String `tfsdk:"resource_group_name"`
	Location             types.String `tfsdk:"location"`
	VMSize               types.String `tfsdk:"vm_size"`
	AdminUsername        types.String `tfsdk:"admin_username"`
	AdminPassword        types.String `tfsdk:"admin_password"`
	NetworkInterfaceIDs  types.List   `tfsdk:"network_interface_ids"`
	OSDisk               types.List   `tfsdk:"os_disk"`
	SourceImageReference types.List   `tfsdk:"source_image_reference"`
	Tags                 types.Map    `tfsdk:"tags"`
}

type OSDiskModel struct {
	Caching            types.String `tfsdk:"caching"`
	StorageAccountType types.String `tfsdk:"storage_account_type"`
	DiskSizeGB         types.Int64  `tfsdk:"disk_size_gb"`
}

type SourceImageReferenceModel struct {
	Publisher types.String `tfsdk:"publisher"`
	Offer     types.String `tfsdk:"offer"`
	SKU       types.String `tfsdk:"sku"`
	Version   types.String `tfsdk:"version"`
}

func NewVirtualMachineResource() resource.Resource {
	return &VirtualMachineResource{}
}

func (r *VirtualMachineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine"
}

func (r *VirtualMachineResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Virtual Machine.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The simulated Virtual Machine ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the Virtual Machine.",
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
			"vm_size": schema.StringAttribute{
				Required:    true,
				Description: "The size of the Virtual Machine (e.g. Standard_DS1_v2).",
			},
			"admin_username": schema.StringAttribute{
				Optional:    true,
				Description: "The admin username for the VM.",
			},
			"admin_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The admin password for the VM.",
			},
			"network_interface_ids": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A list of Network Interface IDs to attach to this VM.",
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "A mapping of tags to assign to the resource.",
			},
		},
		Blocks: map[string]schema.Block{
			"os_disk": schema.ListNestedBlock{
				Description: "The OS Disk configuration.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"caching": schema.StringAttribute{
							Required:    true,
							Description: "The caching type (None, ReadOnly, ReadWrite).",
						},
						"storage_account_type": schema.StringAttribute{
							Required:    true,
							Description: "The storage account type (Standard_LRS, Premium_LRS, etc.).",
						},
						"disk_size_gb": schema.Int64Attribute{
							Optional:    true,
							Description: "The size of the OS Disk in GB.",
						},
					},
				},
			},
			"source_image_reference": schema.ListNestedBlock{
				Description: "The source image reference.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"publisher": schema.StringAttribute{
							Required:    true,
							Description: "The publisher of the image (e.g. Canonical).",
						},
						"offer": schema.StringAttribute{
							Required:    true,
							Description: "The offer of the image (e.g. UbuntuServer).",
						},
						"sku": schema.StringAttribute{
							Required:    true,
							Description: "The SKU of the image (e.g. 18.04-LTS).",
						},
						"version": schema.StringAttribute{
							Required:    true,
							Description: "The version of the image (e.g. latest).",
						},
					},
				},
			},
		},
	}
}

func (r *VirtualMachineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *VirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VirtualMachineModel
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

func (r *VirtualMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *VirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VirtualMachineModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VirtualMachineModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *VirtualMachineResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// No-op: simulated resource
}
