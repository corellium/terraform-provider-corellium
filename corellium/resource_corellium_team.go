package corellium

import (
	"context"
	"io"

	"terraform-provider-corellium/corellium/pkg/api"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumV1TeamResource{}
	_ resource.ResourceWithConfigure = &CorelliumV1TeamResource{}
)

// NewCorelliumV1TeamResource is a helper function to simplify the provider implementation.
func NewCorelliumV1TeamResource() resource.Resource {
	return &CorelliumV1TeamResource{}
}

// CorelliumV1TeamResource is the resource implementation.
type CorelliumV1TeamResource struct {
	client *corellium.APIClient
}

// https://github.com/aimoda/go-corellium-api-client/blob/main/docs/User.md
type V1TeamUserModel struct {
	// Id is the user ID.
	Id types.String `tfsdk:"id"`
}

// V1TeamModel maps the resource schema data.
// https://github.com/aimoda/go-corellium-api-client/blob/main/docs/Team.md
type V1TeamModel struct {
	// Id is the team ID.
	Id types.String `tfsdk:"id"`
	// Label is the team label.
	Label types.String `tfsdk:"label"`
	// Users is the list of users.
	Users []V1TeamUserModel `tfsdk:"users"`
}

// Metadata returns the resource type name.
func (d *CorelliumV1TeamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1team"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_v1team "team" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumV1TeamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Team id",
				Computed:    true,
			},
			"label": schema.StringAttribute{
				Description: "Team label",
				Required:    true,
			},
			"users": schema.ListNestedAttribute{
				Description: "Team users",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "User id",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumV1TeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan V1TeamModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	t := corellium.NewCreateTeam(plan.Label.ValueString())

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	team, r, err := d.client.TeamsApi.V1TeamCreate(auth).CreateTeam(*t).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating team",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error creating team",
			"An unexpected error was encountered trying to create the team:\n\n"+string(b),
		)
		return
	}

	if len(plan.Users) > 0 {
		for _, user := range plan.Users {
			r, err := d.client.TeamsApi.V1AddUserToTeam(auth, team.GetId(), user.Id.ValueString()).Execute()
			if err != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error adding user to team",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Error adding user to team",
					"An unexpected error was encountered trying to add user to team:\n\n"+string(b),
				)
				return
			}
		}
	}

	plan.Id = types.StringValue(team.GetId())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumV1TeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state V1TeamModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	teams, r, err := d.client.TeamsApi.V1Teams(auth).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get the teams",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error to get the teams",
			"An unexpected error was encountered trying to create the team:\n\n"+string(b),
		)
		return
	}

	// NOTICE: This is a workaround for the fact that the API doesn't support getting a single team by ID.
	for _, team := range teams {
		if team.GetId() == state.Id.ValueString() {
			state.Id = types.StringValue(team.Id)
			state.Label = types.StringValue(team.Label)
			state.Users = make([]V1TeamUserModel, len(team.Users))
			// TODO: add the user model instead of the only ID.
			for i, user := range team.Users {
				state.Users[i].Id = types.StringValue(user.Id)
			}

			break
		}
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1TeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state V1TeamModel
	// state is the current state of the resource.

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan V1TeamModel
	// plan is the proposed new state of the resource.

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	if !state.Label.Equal(plan.Label) {
		t := corellium.NewCreateTeam(plan.Label.ValueString())
		r, err := d.client.TeamsApi.V1TeamChange(auth, state.Id.ValueString()).CreateTeam(*t).Execute()
		if err != nil {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error updating team",
					"Coudn't read the response body: "+err.Error(),
				)
				return
			}

			resp.Diagnostics.AddError(
				"Error updating team",
				"An unexpected error was encountered trying to update the team:\n\n"+string(b),
			)
			return
		}
	}

	if plan.Users == nil {
		for i, user := range state.Users {
			r, err := d.client.TeamsApi.V1RemoveUserFromTeam(auth, state.Id.ValueString(), user.Id.ValueString()).Execute()
			if err != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error removing user from team",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Error removing user from team",
					"An unexpected error was encountered trying to remove user from team:\n\n"+string(b),
				)
				return
			}

			state.Users = append(state.Users[:i], state.Users[i+1:]...)
		}

		state.Users = nil
	}

	// NOTICE: The API doesn't support updating the users of a team, so we need to remove the users that are not in the plan
	// and add the users that are in the plan but not in the state.

	// NOTICE: It is probably exist a better way to do this without two loops, but It works for now.

	// When the exist in state, but not in plan, remove the user from the team.
	if state.Users != nil {
		for i, user := range state.Users {
			var found bool
			for _, u := range plan.Users {
				if u.Id.Equal(user.Id) {
					found = true
					break
				}
			}

			if !found {
				r, err := d.client.TeamsApi.V1RemoveUserFromTeam(auth, state.Id.ValueString(), user.Id.ValueString()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error removing user from team",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error removing user from team",
						"An unexpected error was encountered trying to remove user from team:\n\n"+string(b),
					)
					return
				}

				state.Users = append(state.Users[:i], state.Users[i+1:]...)
			}
		}
	}

	// When the exist in plan, but not in state, add the user to the team.
	if plan.Users != nil {
		for _, user := range plan.Users {
			var found bool
			for _, u := range state.Users {
				if u.Id.Equal(user.Id) {
					found = true
					break
				}
			}

			if !found {
				r, err := d.client.TeamsApi.V1AddUserToTeam(auth, state.Id.ValueString(), user.Id.ValueString()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error adding user to team",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error adding user to team",
						"An unexpected error was encountered trying to add user to team:\n\n"+string(b),
					)
					return
				}

				state.Users = append(state.Users, user)
			}
		}
	}

	state.Label = plan.Label

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumV1TeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state V1TeamModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.TeamsApi.V1TeamDelete(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting team",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete team",
			"An unexpected error was encountered trying to delete the team:\n\n"+string(b))
		return
	}
}

// Configure adds the provider configured client to the resource.
func (d *CorelliumV1TeamResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
