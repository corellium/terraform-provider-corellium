package corellium

import (
	"context"
	"net/http"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corellium/corellium/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumV1UserResource{}
	_ resource.ResourceWithConfigure = &CorelliumV1UserResource{}
)

// NewCorelliumV1UserResource is a helper function to simplify the provider implementation.
func NewCorelliumV1UserResource() resource.Resource {
	return &CorelliumV1UserResource{}
}

// CorelliumV1UserResource is the resource implementation.
type CorelliumV1UserResource struct {
	client *corellium.APIClient
}

type V1UserDataModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Label         types.String `tfsdk:"label"`
	Email         types.String `tfsdk:"email"`
	Password      types.String `tfsdk:"password"`
	Administrator types.Bool   `tfsdk:"administrator"`
}

// Metadata returns the resource type name.
func (d *CorelliumV1UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1user"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_v1user "user" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumV1UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the created user",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the user",
				Required:    true,
			},
			"label": schema.StringAttribute{
				Description: "The label of the user",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email of the user",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password of the user",
				Optional:    true,
				Sensitive:   true,
			},
			"administrator": schema.BoolAttribute{
				Description: "The administrator flag of the user",
				Required:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumV1UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state V1UserDataModel

	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API Body expectgs map[string]interface{} as a parameter for the request Body (contains User data)
	userMap := map[string]interface{}{
		"name":          state.Name.ValueString(),
		"label":         state.Label.ValueString(),
		"email":         state.Email.ValueString(),
		"password":      state.Password.ValueString(),
		"administrator": state.Administrator.ValueBool(),
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	// Create the user
	// Just returns a map[string]interface{} with the user ID
	createdUser, r, err := d.client.UsersApi.V1CreateUser(auth).Body(userMap).Execute()
	if err != nil {
		if r.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unable to create the corellium user",
				"You do not have permission to create a user",
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to create the corellium user",
			err.Error(),
		)
		return
	}

	userID, ok := createdUser["id"].(string)
	if !ok {
		// Handle the case where the type assertion fails
		resp.Diagnostics.AddError(
			"Invalid or missing 'id' field in created user object",
			userID,
		)
		return
	}

	state.ID = types.StringValue(userID)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumV1UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// UserAPI does NOT have a GetUser endpoint so we will need to use the TeamsAPI to fetch the user data by ID
	var state V1UserDataModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	teams, _, err := d.client.TeamsApi.V1Teams(auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to fetch the corellium teams",
			err.Error(),
		)
		return
	}

	// Iterate over the teams and find the user
	for _, team := range teams {
		for _, user := range team.Users {
			if user.Id == state.ID.ValueString() {
				state.Name = types.StringValue(user.GetName())
				state.Label = types.StringValue(user.GetLabel())
				state.Email = types.StringValue(user.GetEmail())
				if !user.Administrator.IsSet() {
					state.Administrator = types.BoolValue(user.GetAdministrator())
				}
				break
			}
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from state
	var state V1UserDataModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// update is the proposed new state of the resource.
	var update V1UserDataModel
	diags = req.Plan.Get(ctx, &update)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Populate
	updatedUserMap := map[string]interface{}{
		"name":          update.Name.ValueString(),
		"label":         update.Label.ValueString(),
		"email":         update.Email.ValueString(),
		"administrator": update.Administrator.ValueBool(),
		"password":      update.Password.ValueString(),
	}

	// Update the user
	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	// Takes the user uuID as a parameter and a map[string]interface{} as a body containing the user data to update
	_, _, err := d.client.UsersApi.V1UpdateUser(auth, state.ID.ValueString()).Body(updatedUserMap).Execute()
	// Returns an empty body and a 200 status code on success
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update the corellium user: "+state.ID.ValueString(),
			err.Error(),
		)
		return
	}

	// Set updated state by using plan (update) values if successful. Keep the ID from the current state.
	state.Name = update.Name
	state.Label = update.Label
	state.Email = update.Email
	state.Administrator = update.Administrator
	state.Password = update.Password

	// Set refreshed state with the updated User data.
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumV1UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state V1UserDataModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the user
	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	// Takes the user uuID as a parameter
	_, _, err := d.client.UsersApi.V1DeleteUser(auth, state.ID.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete the corellium user: "+state.ID.ValueString(),
			err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (d *CorelliumV1UserResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
