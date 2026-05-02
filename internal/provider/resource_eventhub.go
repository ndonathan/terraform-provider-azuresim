package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &EventHubResource{}
var _ resource.ResourceWithConfigure = &EventHubResource{}

type EventHubResource struct{ subscriptionID string }

type EventHubModel struct {
	ID                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	NamespaceID         types.String `tfsdk:"namespace_id"`
	PartitionCount      types.Int64  `tfsdk:"partition_count"`
	MessageRetention    types.Int64  `tfsdk:"message_retention"`
	Status              types.String `tfsdk:"status"`
	CaptureDescription  types.List   `tfsdk:"capture_description"`
}

func NewEventHubResource() resource.Resource { return &EventHubResource{} }

func (r *EventHubResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eventhub"
}

func (r *EventHubResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Event Hub.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Event Hub ID."},
			"name": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Event Hub name."},
			"namespace_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Parent namespace ID."},
			"partition_count":   schema.Int64Attribute{Required: true, Description: "Number of partitions."},
			"message_retention": schema.Int64Attribute{Required: true, Description: "Retention in days."},
			"status":            schema.StringAttribute{Optional: true, Description: "`Active`, `Disabled`, or `SendDisabled`."},
		},
		Blocks: map[string]schema.Block{
			"capture_description": schema.ListNestedBlock{
				Description: "Optional capture configuration to Storage / Data Lake.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enabled":             schema.BoolAttribute{Required: true, Description: "Enable capture."},
						"encoding":            schema.StringAttribute{Required: true, Description: "`Avro` or `AvroDeflate`."},
						"interval_in_seconds": schema.Int64Attribute{Optional: true, Description: "Capture interval (60-900)."},
						"size_limit_in_bytes": schema.Int64Attribute{Optional: true, Description: "Capture size limit."},
						"skip_empty_archives": schema.BoolAttribute{Optional: true, Description: "Skip empty archives."},
					},
					Blocks: map[string]schema.Block{
						"destination": schema.ListNestedBlock{
							Description: "Capture destination.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"name":                  schema.StringAttribute{Required: true, Description: "Destination type (`EventHubArchive.AzureBlockBlob` or `EventHubArchive.AzureDataLake`)."},
									"archive_name_format":   schema.StringAttribute{Required: true, Description: "Archive path format."},
									"blob_container_name":   schema.StringAttribute{Required: true, Description: "Target blob container."},
									"storage_account_id":    schema.StringAttribute{Required: true, Description: "Target storage account ID."},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *EventHubResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *EventHubResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EventHubModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/eventhubs/%s",
		strings.TrimRight(plan.NamespaceID.ValueString(), "/"), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EventHubResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EventHubModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EventHubResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EventHubModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/eventhubs/%s",
		strings.TrimRight(plan.NamespaceID.ValueString(), "/"), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *EventHubResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
