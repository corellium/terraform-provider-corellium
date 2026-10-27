package corellium

import (
	"context"
	"io"
	"net/http"

	"terraform-provider-corellium/corellium/pkg/api"

	corellium "github.com/corellium/corellium-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumSnapshotResource{}
	_ resource.ResourceWithConfigure = &CorelliumSnapshotResource{}
)

// NewCorelliumSnapshotResource is a helper function to simplify the provider implementation.
func NewCorelliumSnapshotResource() resource.Resource {
	return &CorelliumSnapshotResource{}
}

// CorelliumSnapshotResource is the resource implementation.
type CorelliumSnapshotResource struct {
	client *corellium.CorelliumClient
}

type SnapshotStatusModel struct {
	// Task is the task name.
	// It can be either "creating", "deleting" or "none.
	Task types.String `tfsdk:"task"`
	// Create is a boolean that indicates if the snapshot is already created.
	Created types.Bool `tfsdk:"created"`
}

// SnapshotModel maps the resource schema data.
// https://github.com/corellium/corellium-go-api-client/blob/main/docs/Snapshot.md
type SnapshotModel struct {
	// Id is the snapshot id.
	Id types.String `tfsdk:"id"`
	// Name is the snapshot name.
	Name types.String `tfsdk:"name"`
	// Instance is the instance id.
	Instance types.String `tfsdk:"instance"`
	// Status is the snapshot status.
	Status *SnapshotStatusModel `tfsdk:"status"`
	// Date is the time when the snapshot was created.
	// Data is a float64 because the API returns a UNIX timestamp.
	Date  types.Float64 `tfsdk:"date"`
	Fresh types.Bool    `tfsdk:"fresh"`
	// Live snapshot (included state and memory).
	Live  types.Bool `tfsdk:"live"`
	Local types.Bool `tfsdk:"local"`
}

// Metadata returns the resource type name.
func (d *CorelliumSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_snapshot "snapshot" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumSnapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Snapshot id",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Snapshot name",
				Required:    true,
			},
			"instance": schema.StringAttribute{
				Description: "Instance id",
				Required:    true,
			},
			"status": schema.SingleNestedAttribute{
				Description: "Snapshot status",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"task": schema.StringAttribute{
						Description: "Snapshot task",
						Computed:    true,
					},
					"created": schema.BoolAttribute{
						Description: "Snapshot created",
						Optional:    true,
						Computed:    true,
					},
				},
			},
			"date": schema.Float64Attribute{
				Description: "Snapshot date",
				Computed:    true,
			},
			"fresh": schema.BoolAttribute{
				Description: "Snapshot fresh",
				Computed:    true,
			},
			"live": schema.BoolAttribute{
				Description: "Snapshot live",
				Computed:    true,
			},
			"local": schema.BoolAttribute{
				Description: "Snapshot local",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SnapshotModel

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createSnapshotBody := corellium.NewCreateSnapshotBody(plan.Name.ValueString())
	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	snapshot, r, err := d.client.Snapshots.Create(auth, plan.Instance.ValueString()).CreateSnapshotBody(*createSnapshotBody).Execute()
	if err != nil {
		if r.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Error creating snapshot",
				"You don't have permissions to create snapshots.",
			)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating snapshot",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error creating snapshot",
			"An unexpected error was encountered trying to create the snapshot:\n\n"+string(b),
		)
		return
	}

	plan.Id = types.StringValue(snapshot.GetId())
	plan.Name = types.StringValue(snapshot.GetName())
	plan.Instance = types.StringValue(snapshot.GetInstance())
	plan.Status = &SnapshotStatusModel{
		Task:    types.StringValue(snapshot.Status.GetTask()),
		Created: types.BoolValue(snapshot.Status.GetCreated()),
	}
	plan.Date = types.Float64Value(float64(snapshot.GetDate()))
	plan.Fresh = types.BoolValue(snapshot.GetFresh())
	plan.Live = types.BoolValue(snapshot.GetLive())
	plan.Local = types.BoolValue(snapshot.GetLocal())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SnapshotModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	snapshot, r, err := d.client.Snapshots.GetById(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get the snapshot",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error to get the snapshot",
			"An unexpected error was encountered trying to create the snapshot:\n\n"+string(b),
		)
		return
	}

	state.Id = types.StringValue(snapshot.GetId())
	state.Name = types.StringValue(snapshot.GetName())
	state.Instance = types.StringValue(snapshot.GetInstance())
	state.Status = &SnapshotStatusModel{
		Task:    types.StringValue(snapshot.Status.GetTask()),
		Created: types.BoolValue(snapshot.Status.GetCreated()),
	}
	state.Date = types.Float64Value(float64(snapshot.GetDate()))
	state.Fresh = types.BoolValue(snapshot.GetFresh())
	state.Live = types.BoolValue(snapshot.GetLive())
	state.Local = types.BoolValue(snapshot.GetLocal())

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state SnapshotModel
	// state is the current state of the resource.

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan SnapshotModel
	// plan is the proposed new state of the resource.

	diags = req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Instance.Equal(plan.Instance) {
		resp.Diagnostics.AddError(
			"Error updating snapshot",
			"It is not possible to update the snapshot's instance",
		)
	}

	createSnapshotBody := corellium.NewCreateSnapshotBody(plan.Name.ValueString())
	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	if !state.Name.Equal(plan.Name) {
		snapshot, r, err := d.client.Snapshots.Patch(auth, state.Instance.ValueString(), state.Id.ValueString()).CreateSnapshotBody(*createSnapshotBody).Execute()
		if err != nil {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating snapshot",
					"Coudn't read the response body: "+err.Error(),
				)
				return
			}

			resp.Diagnostics.AddError(
				"Error updating snapshot",
				"An unexpected error was encountered trying to update the snapshot:\n\n"+string(b),
			)
			return
		}

		state.Name = types.StringValue(snapshot.GetName())
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SnapshotModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.Snapshots.DeleteById(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting snapshot",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete snapshot",
			"An unexpected error was encountered trying to delete the snapshot:\n\n"+string(b))
		return
	}
}

// Configure adds the provider configured client to the resource.
func (d *CorelliumSnapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*ProviderData).Client
}
