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

func sbChildID(namespaceID, kind, name string) string {
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(namespaceID, "/"), kind, name)
}

// --- Service Bus Queue ---

var _ resource.Resource = &ServiceBusQueueResource{}
var _ resource.ResourceWithConfigure = &ServiceBusQueueResource{}

type ServiceBusQueueResource struct{ subscriptionID string }

type ServiceBusQueueModel struct {
	ID                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	NamespaceID                        types.String `tfsdk:"namespace_id"`
	MaxSizeInMegabytes                 types.Int64  `tfsdk:"max_size_in_megabytes"`
	LockDuration                       types.String `tfsdk:"lock_duration"`
	RequiresDuplicateDetection         types.Bool   `tfsdk:"requires_duplicate_detection"`
	RequiresSession                    types.Bool   `tfsdk:"requires_session"`
	DeadLetteringOnMessageExpiration   types.Bool   `tfsdk:"dead_lettering_on_message_expiration"`
	MaxDeliveryCount                   types.Int64  `tfsdk:"max_delivery_count"`
	DefaultMessageTTL                  types.String `tfsdk:"default_message_ttl"`
	DuplicateDetectionHistoryTimeWindow types.String `tfsdk:"duplicate_detection_history_time_window"`
	EnablePartitioning                 types.Bool   `tfsdk:"enable_partitioning"`
	EnableExpress                      types.Bool   `tfsdk:"enable_express"`
	Status                             types.String `tfsdk:"status"`
}

func NewServiceBusQueueResource() resource.Resource { return &ServiceBusQueueResource{} }

func (r *ServiceBusQueueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicebus_queue"
}

func (r *ServiceBusQueueResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates a Service Bus Queue.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Queue ID."},
			"name":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Queue name."},
			"namespace_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Parent namespace ID."},
			"max_size_in_megabytes":                   schema.Int64Attribute{Optional: true, Description: "Max queue size in MB."},
			"lock_duration":                           schema.StringAttribute{Optional: true, Description: "ISO 8601 duration (e.g. `PT30S`)."},
			"requires_duplicate_detection":            schema.BoolAttribute{Optional: true, Description: "Enable duplicate detection."},
			"requires_session":                        schema.BoolAttribute{Optional: true, Description: "Enable sessions."},
			"dead_lettering_on_message_expiration":    schema.BoolAttribute{Optional: true, Description: "Dead-letter on TTL expiry."},
			"max_delivery_count":                      schema.Int64Attribute{Optional: true, Description: "Max delivery count before dead-letter."},
			"default_message_ttl":                     schema.StringAttribute{Optional: true, Description: "Default message TTL (ISO 8601)."},
			"duplicate_detection_history_time_window": schema.StringAttribute{Optional: true, Description: "Duplicate detection window (ISO 8601)."},
			"enable_partitioning":                     schema.BoolAttribute{Optional: true, Description: "Enable partitioning."},
			"enable_express":                          schema.BoolAttribute{Optional: true, Description: "Enable express queue."},
			"status":                                  schema.StringAttribute{Optional: true, Description: "`Active`, `Disabled`, `SendDisabled`, or `ReceiveDisabled`."},
		},
	}
}

func (r *ServiceBusQueueResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ServiceBusQueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceBusQueueModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(sbChildID(plan.NamespaceID.ValueString(), "queues", plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusQueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceBusQueueModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceBusQueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceBusQueueModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(sbChildID(plan.NamespaceID.ValueString(), "queues", plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusQueueResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Service Bus Topic ---

var _ resource.Resource = &ServiceBusTopicResource{}
var _ resource.ResourceWithConfigure = &ServiceBusTopicResource{}

type ServiceBusTopicResource struct{ subscriptionID string }

type ServiceBusTopicModel struct {
	ID                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	NamespaceID                        types.String `tfsdk:"namespace_id"`
	MaxSizeInMegabytes                 types.Int64  `tfsdk:"max_size_in_megabytes"`
	RequiresDuplicateDetection         types.Bool   `tfsdk:"requires_duplicate_detection"`
	DefaultMessageTTL                  types.String `tfsdk:"default_message_ttl"`
	DuplicateDetectionHistoryTimeWindow types.String `tfsdk:"duplicate_detection_history_time_window"`
	EnablePartitioning                 types.Bool   `tfsdk:"enable_partitioning"`
	EnableExpress                      types.Bool   `tfsdk:"enable_express"`
	SupportOrdering                    types.Bool   `tfsdk:"support_ordering"`
	Status                             types.String `tfsdk:"status"`
}

func NewServiceBusTopicResource() resource.Resource { return &ServiceBusTopicResource{} }

func (r *ServiceBusTopicResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicebus_topic"
}

func (r *ServiceBusTopicResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates a Service Bus Topic.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Topic ID."},
			"name":         schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Topic name."},
			"namespace_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Parent namespace ID."},
			"max_size_in_megabytes":                   schema.Int64Attribute{Optional: true, Description: "Max topic size in MB."},
			"requires_duplicate_detection":            schema.BoolAttribute{Optional: true, Description: "Enable duplicate detection."},
			"default_message_ttl":                     schema.StringAttribute{Optional: true, Description: "Default message TTL (ISO 8601)."},
			"duplicate_detection_history_time_window": schema.StringAttribute{Optional: true, Description: "Duplicate detection window (ISO 8601)."},
			"enable_partitioning":                     schema.BoolAttribute{Optional: true, Description: "Enable partitioning."},
			"enable_express":                          schema.BoolAttribute{Optional: true, Description: "Enable express topic."},
			"support_ordering":                        schema.BoolAttribute{Optional: true, Description: "Preserve ordering."},
			"status":                                  schema.StringAttribute{Optional: true, Description: "`Active` or `Disabled`."},
		},
	}
}

func (r *ServiceBusTopicResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ServiceBusTopicResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceBusTopicModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(sbChildID(plan.NamespaceID.ValueString(), "topics", plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusTopicResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceBusTopicModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceBusTopicResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceBusTopicModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(sbChildID(plan.NamespaceID.ValueString(), "topics", plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusTopicResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}

// --- Service Bus Subscription ---

var _ resource.Resource = &ServiceBusSubscriptionResource{}
var _ resource.ResourceWithConfigure = &ServiceBusSubscriptionResource{}

type ServiceBusSubscriptionResource struct{ subscriptionID string }

type ServiceBusSubscriptionModel struct {
	ID                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	TopicID                           types.String `tfsdk:"topic_id"`
	MaxDeliveryCount                  types.Int64  `tfsdk:"max_delivery_count"`
	LockDuration                      types.String `tfsdk:"lock_duration"`
	DefaultMessageTTL                 types.String `tfsdk:"default_message_ttl"`
	DeadLetteringOnFilterEvaluationError types.Bool `tfsdk:"dead_lettering_on_filter_evaluation_error"`
	DeadLetteringOnMessageExpiration  types.Bool   `tfsdk:"dead_lettering_on_message_expiration"`
	EnableBatchedOperations           types.Bool   `tfsdk:"enable_batched_operations"`
	RequiresSession                   types.Bool   `tfsdk:"requires_session"`
	ForwardTo                         types.String `tfsdk:"forward_to"`
	ForwardDeadLetteredMessagesTo     types.String `tfsdk:"forward_dead_lettered_messages_to"`
	Status                            types.String `tfsdk:"status"`
}

func NewServiceBusSubscriptionResource() resource.Resource { return &ServiceBusSubscriptionResource{} }

func (r *ServiceBusSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_servicebus_subscription"
}

func (r *ServiceBusSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates a Service Bus Subscription.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}, Description: "Subscription ID."},
			"name":     schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Subscription name."},
			"topic_id": schema.StringAttribute{Required: true, PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}, Description: "Parent topic ID."},
			"max_delivery_count":                       schema.Int64Attribute{Required: true, Description: "Max delivery count before dead-letter."},
			"lock_duration":                            schema.StringAttribute{Optional: true, Description: "ISO 8601 duration."},
			"default_message_ttl":                      schema.StringAttribute{Optional: true, Description: "Default TTL (ISO 8601)."},
			"dead_lettering_on_filter_evaluation_error": schema.BoolAttribute{Optional: true, Description: "Dead-letter when a filter errors."},
			"dead_lettering_on_message_expiration":     schema.BoolAttribute{Optional: true, Description: "Dead-letter on TTL expiry."},
			"enable_batched_operations":                schema.BoolAttribute{Optional: true, Description: "Enable batched ops."},
			"requires_session":                         schema.BoolAttribute{Optional: true, Description: "Require sessions."},
			"forward_to":                               schema.StringAttribute{Optional: true, Description: "Auto-forward target queue/topic."},
			"forward_dead_lettered_messages_to":        schema.StringAttribute{Optional: true, Description: "Auto-forward dead-letter target."},
			"status":                                   schema.StringAttribute{Optional: true, Description: "`Active` or `Disabled`."},
		},
	}
}

func (r *ServiceBusSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *ServiceBusSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ServiceBusSubscriptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/subscriptions/%s",
		strings.TrimRight(plan.TopicID.ValueString(), "/"), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ServiceBusSubscriptionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ServiceBusSubscriptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServiceBusSubscriptionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(fmt.Sprintf("%s/subscriptions/%s",
		strings.TrimRight(plan.TopicID.ValueString(), "/"), plan.Name.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceBusSubscriptionResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
