package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &CosmosDBAccountResource{}
var _ resource.ResourceWithConfigure = &CosmosDBAccountResource{}

type CosmosDBAccountResource struct {
	subscriptionID string
}

type CosmosDBAccountModel struct {
	ID                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	ResourceGroupName               types.String `tfsdk:"resource_group_name"`
	Location                        types.String `tfsdk:"location"`
	OfferType                       types.String `tfsdk:"offer_type"`
	Kind                            types.String `tfsdk:"kind"`
	AutomaticFailoverEnabled        types.Bool   `tfsdk:"automatic_failover_enabled"`
	MultipleWriteLocationsEnabled   types.Bool   `tfsdk:"multiple_write_locations_enabled"`
	IPRangeFilter                   types.String `tfsdk:"ip_range_filter"`
	PublicNetworkAccessEnabled      types.Bool   `tfsdk:"public_network_access_enabled"`
	LocalAuthenticationDisabled     types.Bool   `tfsdk:"local_authentication_disabled"`
	ConsistencyPolicy               types.List   `tfsdk:"consistency_policy"`
	GeoLocation                     types.List   `tfsdk:"geo_location"`
	Capabilities                    types.List   `tfsdk:"capabilities"`
	Endpoint                        types.String `tfsdk:"endpoint"`
	PrimaryKey                      types.String `tfsdk:"primary_key"`
	SecondaryKey                    types.String `tfsdk:"secondary_key"`
	PrimaryReadonlyKey              types.String `tfsdk:"primary_readonly_key"`
	SecondaryReadonlyKey            types.String `tfsdk:"secondary_readonly_key"`
	Tags                            types.Map    `tfsdk:"tags"`
}

func NewCosmosDBAccountResource() resource.Resource { return &CosmosDBAccountResource{} }

func (r *CosmosDBAccountResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cosmosdb_account"
}

func (r *CosmosDBAccountResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Simulates an Azure Cosmos DB Account.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true, Description: "The Cosmos DB Account ID.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required: true, Description: "Account name (3-44 lowercase alphanumeric/dashes, globally unique).",
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
			"offer_type":                       schema.StringAttribute{Required: true, Description: "Always `Standard`."},
			"kind":                             schema.StringAttribute{Optional: true, Description: "`GlobalDocumentDB` (default), `MongoDB`, or `Parse`."},
			"automatic_failover_enabled":       schema.BoolAttribute{Optional: true, Description: "Enable automatic failover."},
			"multiple_write_locations_enabled": schema.BoolAttribute{Optional: true, Description: "Enable multi-master."},
			"ip_range_filter":                  schema.StringAttribute{Optional: true, Description: "Comma-separated IP ranges."},
			"public_network_access_enabled":    schema.BoolAttribute{Optional: true, Description: "Allow public network access."},
			"local_authentication_disabled":    schema.BoolAttribute{Optional: true, Description: "Disable key-based auth."},
			"endpoint": schema.StringAttribute{
				Computed: true, Description: "Simulated endpoint (`https://<name>.documents.azure.com:443/`).",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"primary_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"secondary_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated secondary key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"primary_readonly_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated primary readonly key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"secondary_readonly_key": schema.StringAttribute{
				Computed: true, Sensitive: true, Description: "Simulated secondary readonly key.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"tags": schema.MapAttribute{Optional: true, ElementType: types.StringType, Description: "Tags."},
		},
		Blocks: map[string]schema.Block{
			"consistency_policy": schema.ListNestedBlock{
				Description: "Consistency policy.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"consistency_level":         schema.StringAttribute{Required: true, Description: "`Strong`, `BoundedStaleness`, `Session`, `ConsistentPrefix`, or `Eventual`."},
						"max_interval_in_seconds":   schema.Int64Attribute{Optional: true, Description: "Bounded-staleness window in seconds."},
						"max_staleness_prefix":      schema.Int64Attribute{Optional: true, Description: "Bounded-staleness lag (operations)."},
					},
				},
			},
			"geo_location": schema.ListNestedBlock{
				Description: "Geographic replica.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"location":          schema.StringAttribute{Required: true, Description: "Azure region."},
						"failover_priority": schema.Int64Attribute{Required: true, Description: "Failover priority (0 = primary)."},
						"zone_redundant":    schema.BoolAttribute{Optional: true, Description: "Use zone redundancy."},
					},
				},
			},
			"capabilities": schema.ListNestedBlock{
				Description: "Account capability (e.g. `EnableServerless`, `EnableMongo`).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{Required: true, Description: "Capability name."},
					},
				},
			},
		},
	}
}

func (r *CosmosDBAccountResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.subscriptionID = req.ProviderData.(string)
}

func (r *CosmosDBAccountResource) applyComputed(plan *CosmosDBAccountModel) {
	name := plan.Name.ValueString()
	rg := plan.ResourceGroupName.ValueString()
	plan.ID = types.StringValue(fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s",
		r.subscriptionID, rg, name,
	))
	plan.Endpoint = types.StringValue(fmt.Sprintf("https://%s.documents.azure.com:443/", name))
	mkKey := func(prefix string) string {
		sum := sha1.Sum([]byte(prefix + "/" + rg + "/" + name))
		padded := make([]byte, 64)
		for i := range padded {
			padded[i] = sum[i%len(sum)]
		}
		return base64.StdEncoding.EncodeToString(padded)
	}
	plan.PrimaryKey = types.StringValue(mkKey("cosmos-primary"))
	plan.SecondaryKey = types.StringValue(mkKey("cosmos-secondary"))
	plan.PrimaryReadonlyKey = types.StringValue(mkKey("cosmos-primary-ro"))
	plan.SecondaryReadonlyKey = types.StringValue(mkKey("cosmos-secondary-ro"))
}

func (r *CosmosDBAccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CosmosDBAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CosmosDBAccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CosmosDBAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CosmosDBAccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CosmosDBAccountModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state CosmosDBAccountModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = state.ID
	plan.Endpoint = state.Endpoint
	plan.PrimaryKey = state.PrimaryKey
	plan.SecondaryKey = state.SecondaryKey
	plan.PrimaryReadonlyKey = state.PrimaryReadonlyKey
	plan.SecondaryReadonlyKey = state.SecondaryReadonlyKey
	r.applyComputed(&plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CosmosDBAccountResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
}
