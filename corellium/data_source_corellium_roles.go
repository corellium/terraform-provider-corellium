package corellium

import (
	"context"
	"io"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corellium/corellium/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &V1RolesDataSource{}
	_ datasource.DataSourceWithConfigure = &V1RolesDataSource{}
)

// NewCorelliumV1RolesDataSource is a helper function to simplify the provider implementation.
func NewCorelliumV1RolesDataSource() datasource.DataSource {
	return &V1RolesDataSource{}
}

// V1RolesDataSource is the data source implementation.
type V1RolesDataSource struct {
	client *corellium.APIClient
}

type V1RoleModel struct {
	// Role is the role name.
	Role types.String `tfsdk:"role"`
	// Project is the project ID.
	Project types.String `tfsdk:"project"`
	// User is the user ID.
	User types.String `tfsdk:"user"`
}

// V1RolesModel maps the data source schema data.
// https://github.com/aimoda/go-corellium-api-client/blob/main/docs/RolesApi.md#v1roles
type V1RolesModel struct {
	// Id is the data source required Id.
	// Each data source should has a Id.
	Id types.String `tfsdk:"id"`
	// Project is the project ID.
	// Project is optional, if not set, all the roles from all projects and users will be returned.
	// If set, only the roles from the given project will be returned.
	Project types.String `tfsdk:"project"`
	// TODO: add filter for user and team.
	// Roles is the list of roles.
	Roles []V1RoleModel `tfsdk:"roles"`
}

// Metadata returns the data source type name.
func (d *V1RolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1roles"
	// TypeName is the name of the data resource type, which must be unique within the provider.
	// This is used to identify the data resource type in state and plan files.
	// i.e: data corellium_v1roles "roles" { ... }
}

// Schema defines the schema for the data source.
func (d *V1RolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Data source ID",
				Computed:    true,
			},
			"project": schema.StringAttribute{
				Description: "Project ID",
				Optional:    true,
			},
			"roles": schema.ListNestedAttribute{
				Description: "List of roles",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role": schema.StringAttribute{
							Description: "Role name",
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("admin", "_member_"), // NOTICE: Is _member_ right?
							},
						},
						"project": schema.StringAttribute{
							Description: "Project ID",
							Optional:    true,
							Computed:    true,
						},
						"user": schema.StringAttribute{
							Description: "User ID",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *V1RolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state V1RolesModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	// auth is the context with the access token, what is required by the API client.
	roles, r, err := d.client.RolesApi.V1Roles(auth).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error gettings roles",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error gettings roles",
			"An unexpected error was encountered trying to get the roles:\n\n"+string(b),
		)
		return
	}

	if len(roles) == 0 {
		resp.Diagnostics.AddWarning(
			"No roles found",
			"No roles were found for the given project and user.",
		)
		return
	}

	appendRoleToState := func(role, project, user string) {
		state.Roles = append(state.Roles, V1RoleModel{
			Role:    types.StringValue(role),
			Project: types.StringValue(project),
			User:    types.StringValue(user),
		})
	}

	// Check if the list per role as project filter.
	// If the project is set, but no roles are found, we return all the roles from all projects and users.
	var found bool
	if !state.Project.IsNull() {
		for _, role := range roles {
			if role.GetProject() == state.Project.ValueString() {
				found = true

				appendRoleToState(role.Role, role.Project, role.User)
			}
		}
	}

	// If the project is not set, we return all the roles from all projects and users.
	if !found {
		for _, role := range roles {
			appendRoleToState(role.Role, role.Project, role.User)
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating UUID",
			"An unexpected error was encountered trying to generate the ID:\n\n"+err.Error(),
		)
		return
	}

	state.Id = types.StringValue(id)

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *V1RolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
