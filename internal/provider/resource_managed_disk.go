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

var _ resource.Resource = &ManagedDiskResource{}
var _ resource.ResourceWithConfigure = &ManagedDiskResource{}

type ManagedDiskResource struct {
	subscriptionID string
}

type ManagedDiskModel struct {
	ID                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	ResourceGroupName          types.String `tfsdk:"resource_group_name"`
	Location                   types.String `tfsdk:"location"`
	StorageAccountType         types.String `tfsdk:"storage_account_type"`
	CreateOption               types.String `tfsdk:"create_option"`
	DiskSizeGB                 types.Int64  `tfsdk:"disk_size_gb"`
	SourceURI                  types.String `tfsdk:"source_uri"`
	SourceResourceID           types.String `tfsdk:"source_resource_id"`
	ImageReferenceID           types.String `tfsdk:"image_reference_id"`
	OSType                     types.String `tfsdk:"os_type"`
	Tier                       types.String `tfsdk:"tier"`
	MaxShares                  types.Int64  `tfsdk:"max_shares"`
	Zone                       types.String `tfsdk:"zone"`
	NetworkAccessPolicy        types.String `tfsdk:"network_access_policy"`
	PublicNetworkAccessEnabled types.Bool   `tfsdk:"public_network_access_enabled"`
	DiskIOPSReadWrite          types.Int64  `tfsdk:"disk_iops_read_write"`
	DiskMBpsReadWrite          types.Int64  `tfsdk:"disk_mbps_read_write"`
	Tags                       types.Map    `tfsdk:"tags"`
}

func NewManagedDiskResource() resource.Resource {
	return &ManagedDiskResource{}
}

func (r *ManagedDiskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_managed_disk"
}

func (r *ManagedDiskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Managed Disk.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The simulated Managed Disk ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Name of the Managed Disk.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"resource_group_name": schema.StringAttribute{
				Required: true, Description: "Name of the Resource Group.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"location": schema.StringAttribute{
				Required: true, Description: "Azure region.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"storage_account_type": schema.StringAttribute{
				Required:    true,
				Description: "`Standard_LRS`, `Premium_LRS`, `StandardSSD_LRS`, `UltraSSD_LRS`, `Premium_ZRS`, `StandardSSD_ZRS`.",
			},
			"create_option": schema.StringAttribute{
				Required: true, Description: "`Empty`, `Copy`, `FromImage`, `Import`, `Restore`, `Upload`.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"disk_size_gb": schema.Int64Attribute{
				Optional: true, Computed: true, Description: "Disk size in GB.",
			},
			"source_uri":         schema.StringAttribute{Optional: true, Description: "Source URI when `create_option` is `Import`."},
			"source_resource_id": schema.StringAttribute{Optional: true, Description: "Source resource ID when `create_option` is `Copy` or `Restore`."},
			"image_reference_id": schema.StringAttribute{Optional: true, Description: "Image reference ID when `create_option` is `FromImage`."},
			"os_type":            schema.StringAttribute{Optional: true, Description: "`Linux` or `Windows`."},
			"tier":               schema.StringAttribute{Optional: true, Description: "Performance tier (e.g. `P30`)."},
			"max_shares":         schema.Int64Attribute{Optional: true, Description: "Maximum number of VMs that can share the disk."},
			"zone": schema.StringAttribute{
				Optional: true, Description: "Availability zone (`1`, `2`, or `3`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"network_access_policy":         schema.StringAttribute{Optional: true, Description: "`AllowAll`, `AllowPrivate`, or `DenyAll`."},
			"public_network_access_enabled": schema.BoolAttribute{Optional: true, Description: "Whether public network access is allowed."},
			"disk_iops_read_write":          schema.Int64Attribute{Optional: true, Description: "Provisioned IOPS (UltraSSD/PremiumV2 only)."},
			"disk_mbps_read_write":          schema.Int64Attribute{Optional: true, Description: "Provisioned bandwidth in MB/s (UltraSSD/PremiumV2 only)."},
			"tags": schema.MapAttribute{
				Optional: true, ElementType: types.StringType,
				Description: "Tags.",
			},
		},
	}
}

func (r *ManagedDiskResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ManagedDiskResource) applyComputed(plan *ManagedDiskModel) {
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/disks/%s",
		r.subscriptionID, plan.ResourceGroupName.ValueString(), plan.Name.ValueString(),
	))
	if plan.DiskSizeGB.IsNull() || plan.DiskSizeGB.IsUnknown() {
		plan.DiskSizeGB = types.Int64Value(30)
	}
}

func (r *ManagedDiskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ManagedDiskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ManagedDiskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ManagedDiskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ManagedDiskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ManagedDiskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state ManagedDiskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ManagedDiskResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
