package corellium

import (
	"context"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-corellium/corellium/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumV1ProjectResource{}
	_ resource.ResourceWithConfigure = &CorelliumV1ProjectResource{}
)

// NewCorelliumV1ProjectResource is a helper function to simplify the provider implementation.
func NewCorelliumV1ProjectResource() resource.Resource {
	return &CorelliumV1ProjectResource{}
}

// CorelliumV1ProjectResource is the resource implementation.
type CorelliumV1ProjectResource struct {
	client *corellium.APIClient
}

type V1ProjectUserModel struct {
	// Id is the user ID.
	Id types.String `tfsdk:"id"`
	// Name is the user name.
	Name types.String `tfsdk:"name"`
	// Label is the user label.
	Label types.String `tfsdk:"label"`
	// Email is the user email.
	Email types.String `tfsdk:"email"`
	// Role is the user role.
	Role types.String `tfsdk:"role"`
}

type V1ProjectTeamModel struct {
	// Id is the team ID.
	Id types.String `tfsdk:"id"`
	// Label is the team label.
	Label types.String `tfsdk:"label"`
	// Role is the team role.
	Role types.String `tfsdk:"role"`
}

type V1ProjectSettingsModel struct {
	// Version is the project version.
	Version types.Number `tfsdk:"version"`
	// InternetAccess is a boolean that defines if the project has Internet access.
	InternetAccess types.Bool `tfsdk:"internet_access"`
	// Dhcp is a boolean that defines if the project has DHCP enabled.
	Dhcp types.Bool `tfsdk:"dhcp"`
}

type V1ProjectQuotasModel struct {
	// Name is the project name.
	Name types.String `tfsdk:"name"`
	// Core is the project cores quota.
	Cores types.Number `tfsdk:"cores"`
	// Instances is the project instances quota.
	// Instances is computed from cores. Instances are equal to cores * 2.5
	Instances types.Number `tfsdk:"instances"`
	// Ram is the project RAM quota.
	// Ram is computed from cores. Ram is equal to cores * 6144
	Ram types.Number `tfsdk:"ram"`
}

type V1ProjectKeyModel struct {
	// Id is the Identifier project key ID.
	// Id is called Identifier in the API.
	Id types.String `tfsdk:"id"`
	// Label is the project key label.
	Label types.String `tfsdk:"label"`
	// Kind is the project key kind.
	// It can be "ssh" or "adb".
	Kind types.String `tfsdk:"kind"`
	// Key is the project key.
	Key types.String `tfsdk:"key"`
	// Fingerprint is the project key fingerprint.
	Fingerprint types.String `tfsdk:"fingerprint"`
	// CreateAt is the project key creation date.
	CreatedAt types.String `tfsdk:"created_at"`
	// UpdateAt is the project key last update date.
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// V1ProjectModel maps the resource schema data.
// https://github.com/aimoda/go-corellium-api-client/blob/main/docs/Project.md
type V1ProjectModel struct { // TODO: add quotas_used model to the schema.
	// Id is the project ID.
	// The project ID is a uuid, universally unique identifier.
	Id types.String `tfsdk:"id"`
	// Name is the project name.
	Name types.String `tfsdk:"name"`
	// Settings is the project settings.
	Settings *V1ProjectSettingsModel `tfsdk:"settings"`
	// Quotas is the project quotas.
	Quotas *V1ProjectQuotasModel `tfsdk:"quotas"`
	// Users is the project users.
	Users []V1ProjectUserModel `tfsdk:"users"`
	// Teams is the project teams.
	Teams []V1ProjectTeamModel `tfsdk:"teams"`
	// Keys is a list of the project authroized keys.
	Keys []V1ProjectKeyModel `tfsdk:"keys"`
	// CreatedAt is the project creation date.
	CreatedAt types.String `tfsdk:"created_at"`
	// UpdatedAt is the project last update date.
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (d *CorelliumV1ProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1project"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_v1project "project" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumV1ProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Project id",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Project name",
				Required:    true,
			},
			"settings": schema.SingleNestedAttribute{
				Description: "Project settings",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"version": schema.NumberAttribute{
						Description: "Project version",
						Required:    true,
					},
					"internet_access": schema.BoolAttribute{
						Description: "Project internet access",
						Required:    true,
					},
					"dhcp": schema.BoolAttribute{
						Description: "Project dhcp",
						Required:    true,
					},
				},
			},
			"quotas": schema.SingleNestedAttribute{
				Description: "Project quotas",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "Project quota name",
						Computed:    true,
					},
					"cores": schema.NumberAttribute{
						Description: "Project quota cores",
						Required:    true,
					},
					"instances": schema.NumberAttribute{
						Description: "Project quota instances",
						Computed:    true,
					},
					"ram": schema.NumberAttribute{
						Description: "Project quota ram",
						Computed:    true,
					},
				},
			},
			"users": schema.ListNestedAttribute{
				Description: "Project users",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Project user id",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "Project user name",
							Computed:    true,
						},
						"label": schema.StringAttribute{
							Description: "Project user label",
							Computed:    true,
						},
						"email": schema.StringAttribute{
							Description: "Project user email",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "Project user role",
							Required:    true,
						},
					},
				},
			},
			"teams": schema.ListNestedAttribute{
				Description: "Project teams",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Project team id",
							Required:    true,
						},
						"label": schema.StringAttribute{
							Description: "Project team label",
							Computed:    true,
						},
						"role": schema.StringAttribute{
							Description: "Project team role",
							Required:    true,
						},
					},
				},
			},
			"keys": schema.ListNestedAttribute{
				Description: "Project keys",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ProjectKey ID",
							Computed:    true,
						},
						"label": schema.StringAttribute{
							Description: "ProjectKey label",
							Required:    true,
						},
						"kind": schema.StringAttribute{
							Description: "ProjectKey kind",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("ssh", "adb"),
							},
						},
						"key": schema.StringAttribute{
							Description: "ProjectKey key",
							Required:    true,
						},
						"fingerprint": schema.StringAttribute{
							Description: "ProjectKey fingerprint",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "ProjectKey creation date",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "ProjectKey last update date",
							Computed:    true,
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Project created at",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Project updated at",
				Optional:    true, // TODO: Check if the `Optional` flag is required.
				Computed:    true,
			},
		},
	}
}

// locker is a helper struct to force the creation of each project sequentially.
// This is required becuase Terraform Framework executes each create operation concurrently, what causes problems when
// creating multiple projects at the same time.
type locker struct {
	mu *sync.Mutex
}

var l locker = locker{mu: &sync.Mutex{}}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumV1ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan V1ProjectModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())

	l.mu.Lock()
	defer l.mu.Unlock()
	projects, r, err := d.client.ProjectsApi.V1GetProjects(auth).Execute()
	if err != nil {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating project",
				"Coudn't read the response body from project name check: "+err.Error(),
			)
			return
		}
	}

	for _, project := range projects {
		if project.GetName() == plan.Name.ValueString() {
			resp.Diagnostics.AddError(
				"Error creating project",
				"A project with the name "+plan.Name.ValueString()+" already exists",
			)
			return
		}
	}

	BigFloatToFloat32 := func(bf *big.Float) float32 {
		f, _ := bf.Float32()
		return f
	}

	p := corellium.NewProjectWithDefaults()
	p.SetName(plan.Name.ValueString())

	s := corellium.NewProjectSettingsWithDefaults()
	s.SetVersion(BigFloatToFloat32(plan.Settings.Version.ValueBigFloat()))
	s.SetInternetAccess(plan.Settings.InternetAccess.ValueBool())
	s.SetDhcp(plan.Settings.Dhcp.ValueBool())

	p.SetSettings(*s)

	q := corellium.NewProjectQuotaWithDefaults()
	q.SetCores(BigFloatToFloat32(plan.Quotas.Cores.ValueBigFloat()))

	p.SetQuotas(*q)

	created, r, err := d.client.ProjectsApi.V1CreateProject(auth).Project(*p).Execute()
	if err != nil {
		if r.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Error creating project",
				"You don't have permission to create a project.",
			)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error creating project",
			"An unexpected error was encountered trying to create the project:\n\n"+string(b),
		)
		return
	}

	plan.Id = types.StringValue(created.GetId())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, r, err := d.client.ProjectsApi.V1GetProject(auth, created.GetId()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get the creatd project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error to get created project",
			"An unexpected error was encountered trying to create the project:\n\n"+string(b),
		)
		return
	}

	if len(plan.Users) > 0 {
		for i, user := range plan.Users {
			teams, r, err := d.client.TeamsApi.V1Teams(auth).Execute()
			if err != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error get team with all users",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Error to get the team with all users",
					"An unexpected error was encountered trying to get the team with all users:\n\n"+string(b),
				)
				return
			}

			// NOTICE: This is a workaround for the fact that the API doesn't support getting a single team by ID.
			var found bool
			for _, t := range teams {
				if t.Id == "all-users" {
					for _, u := range t.Users {
						if u.Id == user.Id.ValueString() {
							user.Id = types.StringValue(u.Id)
							user.Name = types.StringValue(u.Name)
							user.Label = types.StringValue(u.Label)
							user.Email = types.StringValue(u.Email)

							found = true
							break
						}
					}

					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Error get the users",
					"User with ID "+user.Id.ValueString()+" not found",
				)
				return
			}

			r, err = d.client.RolesApi.V1AddUserRoleToProject(auth, project.GetId(), user.Id.ValueString(), user.Role.ValueString()).Execute()
			if err != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error adding user to project",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Error adding user to project",
					"An unexpected error was encountered trying to add user to project:\n\n"+string(b),
				)
				return
			}

			plan.Users[i] = user
		}
	}

	if len(plan.Teams) > 0 {
		for i, team := range plan.Teams {
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
			var found bool
			for _, t := range teams {
				if t.Id == team.Id.ValueString() {
					team.Id = types.StringValue(t.Id)
					team.Label = types.StringValue(t.Label)
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Error get the teams",
					"Team with ID "+team.Id.ValueString()+" not found",
				)
				return
			}

			r, err = d.client.RolesApi.V1AddTeamRoleToProject(auth, project.GetId(), team.Id.ValueString(), team.Role.ValueString()).Execute()
			if err != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error adding team to project",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Error adding team to project",
					"An unexpected error was encountered trying to add team to project:\n\n"+string(b),
				)
				return
			}

			plan.Teams[i] = team
		}
	}

	if plan.Keys != nil && len(plan.Keys) > 0 {
		for i, key := range plan.Keys {
			p := corellium.NewProjectKey(key.Kind.ValueString(), key.Key.ValueString())
			projectKey, r, err := d.client.ProjectsApi.V1AddProjectKey(auth, created.Id).ProjectKey(*p).Execute()
			if err != nil {
				if r.StatusCode == http.StatusForbidden {
					resp.Diagnostics.AddError(
						"Error creating project key",
						"You don't have permission to create an project key in this project.",
					)
					return
				}

				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error creating project key",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Error creating project key",
					"An unexpected error was encountered trying to create the project key:\n\n"+string(b),
				)
				return
			}

			plan.Keys[i].Id = types.StringValue(projectKey.GetIdentifier())
			plan.Keys[i].Label = types.StringValue(plan.Keys[i].Label.ValueString())
			plan.Keys[i].Kind = types.StringValue(projectKey.GetKind())
			plan.Keys[i].Key = types.StringValue(projectKey.GetKey())
			plan.Keys[i].Fingerprint = types.StringValue(projectKey.GetFingerprint())
			plan.Keys[i].CreatedAt = types.StringValue(projectKey.GetCreatedAt().String())
			plan.Keys[i].UpdatedAt = types.StringValue(projectKey.GetUpdatedAt().String())

			resp.State.Set(ctx, plan)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	plan.Id = types.StringValue(project.GetId())
	plan.Name = types.StringValue(project.GetName())

	plan.Settings.Version = types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion())))
	plan.Settings.InternetAccess = types.BoolValue(project.Settings.GetInternetAccess())
	plan.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())

	plan.Quotas.Name = types.StringValue(project.GetName())
	plan.Quotas.Cores = types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores())))
	plan.Quotas.Instances = types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances())))
	plan.Quotas.Ram = types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam())))

	plan.CreatedAt = types.StringValue(time.Now().Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumV1ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state V1ProjectModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	project, r, err := d.client.ProjectsApi.V1GetProject(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to read project",
			"An unexpected error was encountered trying to read the image:\n\n"+string(b))
		return
	}

	if state.Users != nil {
		for i, user := range state.Users {
			teams, r, err := d.client.TeamsApi.V1Teams(auth).Execute()
			if err != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error to get team with all users",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Error to get the teams",
					"An unexpected error was encountered trying to get the team with all users:\n\n"+string(b),
				)
				return
			}

			// NOTICE: This is a workaround for the fact that the API doesn't support getting a single team by ID.
			var found bool
			for _, t := range teams {
				if t.Id == "all-users" {
					for _, u := range t.Users {
						if u.Id == user.Id.ValueString() {
							user.Id = types.StringValue(u.Id)
							user.Name = types.StringValue(u.Name)
							user.Email = types.StringValue(u.Email)

							found = true
							break
						}
					}

					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Error get the users",
					"User with ID "+user.Id.ValueString()+" not found",
				)
				return
			}

			state.Users[i] = user
		}
	}

	if state.Teams != nil {
		for i, team := range state.Teams {
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

			var found bool
			// NOTICE: This is a workaround for the fact that the API doesn't support getting a single team by ID.
			for _, t := range teams {
				if t.Id == team.Id.ValueString() {
					team.Id = types.StringValue(t.Id)
					team.Label = types.StringValue(t.Label)
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Error get the teams",
					"Team with ID "+team.Id.ValueString()+" not found",
				)
				return
			}

			state.Teams[i] = team
		}
	}

	if state.Keys != nil {
		for _, key := range state.Keys {
			projectKeys, r, err := d.client.ProjectsApi.V1GetProjectKeys(auth, project.Id).Execute()
			if err != nil {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error reading the project keys",
						"Coudn't read the response body: "+err.Error(),
					)
					return
				}

				resp.Diagnostics.AddError(
					"Unable to read project keys",
					"An unexpected error was encountered trying to read the project keys from the project:\n\n"+string(b))
				return
			}

			var found bool
			for _, k := range projectKeys {
				if k.GetIdentifier() == key.Id.ValueString() {
					key.Id = types.StringValue(k.GetIdentifier())
					key.Label = types.StringValue(key.Label.String())
					key.Kind = types.StringValue(k.GetKind())
					key.Key = types.StringValue(k.GetKey())
					key.Fingerprint = types.StringValue(k.GetFingerprint())
					key.CreatedAt = types.StringValue(k.GetCreatedAt().String())
					key.UpdatedAt = types.StringValue(k.GetUpdatedAt().String())
					found = true
					break
				}
			}

			if !found {
				resp.Diagnostics.AddError(
					"Error reading the project keys",
					"Project key with ID "+key.Id.ValueString()+" not found",
				)
				return
			}
		}
	}

	state.Id = types.StringValue(project.GetId())
	state.Name = types.StringValue(project.GetName())

	state.Settings.Version = types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion())))
	state.Settings.InternetAccess = types.BoolValue(project.Settings.GetInternetAccess())
	state.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())

	state.Quotas.Name = types.StringValue(project.GetName())
	state.Quotas.Cores = types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores())))
	state.Quotas.Instances = types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances())))
	state.Quotas.Ram = types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam())))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state V1ProjectModel
	// state is the current state of the resource.

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan V1ProjectModel
	// plan is the proposed new state of the resource.

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	BigFloatToFloat32 := func(bf *big.Float) float32 {
		f, _ := bf.Float32()
		return f
	}

	p := corellium.NewProject(state.Id.ValueString())
	p.SetName(plan.Name.ValueString())

	s := corellium.NewProjectSettings()
	s.SetVersion(BigFloatToFloat32(plan.Settings.Version.ValueBigFloat()))
	s.SetInternetAccess(plan.Settings.InternetAccess.ValueBool())
	s.SetDhcp(plan.Settings.Dhcp.ValueBool())

	p.SetSettings(*s)

	q := corellium.NewProjectQuota()
	q.SetCores(BigFloatToFloat32(plan.Quotas.Cores.ValueBigFloat()))

	p.SetQuotas(*q)

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	project, r, err := d.client.ProjectsApi.V1UpdateProject(auth, state.Id.ValueString()).Project(*p).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error updating project",
			"An unexpected error was encountered trying to update the project:\n\n"+string(b),
		)
		return
	}

	if len(state.Users) > 0 {
		for i, user := range state.Users {
			var found bool
			for _, u := range plan.Users {
				if u.Id.Equal(user.Id) {
					found = true
					break
				}
			}

			if !found {
				r, err := d.client.RolesApi.V1RemoveUserRoleFromProject(auth, state.Id.ValueString(), user.Id.ValueString(), user.Role.ValueString()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error removing user from project",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error removing user from project",
						"An unexpected error was encountered trying to remove user from project:\n\n"+string(b),
					)
					return
				}

				// This snippet removes a user from the state on each iteration.
				if len(state.Users) > 1 {
					if i < len(state.Users)-1 {
						// Removes the user from the state by copying the slice without the user.
						state.Users = append(state.Users[:i], state.Users[i+1:]...)
					} else {
						// However, if the user is the last one, we can just remove it from the slice.
						state.Users = state.Users[:i]
					}
				} else {
					// If the users attribute is empty, we need to set the state to empty too.
					state.Users = []V1ProjectUserModel{}

					// However, when the plan has users atributes set to nil, we need to set the state to nil too.
					if plan.Users == nil {
						state.Users = nil
					}
				}
			}
		}
	}

	if len(plan.Users) > 0 {
		for _, user := range plan.Users {
			var found bool
			for _, u := range state.Users {
				if u.Id.Equal(user.Id) {
					found = true
					break
				}
			}

			if !found {
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
						"An unexpected error was encountered trying to get the team with all users:\n\n"+string(b),
					)
					return
				}

				// NOTICE: This is a workaround for the fact that the API doesn't support getting a single team by ID.
				for _, t := range teams {
					if t.Id == "all-users" {
						for _, u := range t.Users {
							if u.Id == user.Id.ValueString() {
								user.Id = types.StringValue(u.Id)
								user.Name = types.StringValue(u.Name)
								user.Label = types.StringValue(u.Label)
								user.Email = types.StringValue(u.Email)

								break
							}
						}

						break
					}
				}

				r, err = d.client.RolesApi.V1AddUserRoleToProject(auth, state.Id.ValueString(), user.Id.ValueString(), user.Role.ValueString()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error adding user to project",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error adding user to project",
						"An unexpected error was encountered trying to add user to project:\n\n"+string(b),
					)
					return
				}

				state.Users = append(state.Users, user)
			}
		}
	}

	if len(state.Teams) > 0 {
		for i, team := range state.Teams {
			var found bool
			for _, t := range plan.Teams {
				if t.Id.Equal(team.Id) {
					found = true
					break
				}
			}

			if !found {
				r, err := d.client.RolesApi.V1RemoveTeamRoleFromProject(auth, state.Id.ValueString(), team.Id.ValueString(), team.Role.ValueString()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error removing team from project",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error removing team from project",
						"An unexpected error was encountered trying to remove team from project:\n\n"+string(b),
					)
					return
				}

				// This snippet removes a team from the state on each iteration.
				if len(state.Teams) > 1 {
					if i < len(state.Teams)-1 {
						// Removes the team from the state by copying the slice without the team.
						state.Teams = append(state.Teams[:i], state.Teams[i+1:]...)
					} else {
						// However, if the team is the last one, we can just remove it from the slice.
						state.Teams = state.Teams[:i]
					}
				} else {
					// If the teams attribute is empty, we need to set the state to empty too.
					state.Teams = []V1ProjectTeamModel{}

					// However, when the plan has teams atributes set to nil, we need to set the state to nil too.
					if plan.Teams == nil {
						state.Teams = nil
					}
				}
			}
		}
	}

	if len(plan.Teams) > 0 {
		for _, team := range plan.Teams {
			var found bool
			for _, t := range state.Teams {
				if t.Id.Equal(team.Id) {
					found = true
					break
				}
			}

			if !found {
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
				for _, t := range teams {
					if t.Id == team.Id.ValueString() {
						team.Id = types.StringValue(t.Id)
						team.Label = types.StringValue(t.Label)

						break
					}
				}

				r, err = d.client.RolesApi.V1AddTeamRoleToProject(auth, state.Id.ValueString(), team.Id.ValueString(), team.Role.ValueString()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error adding team to project",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error adding team to project",
						"An unexpected error was encountered trying to add team to project:\n\n"+string(b),
					)
					return
				}

				state.Teams = append(state.Teams, team)
			}
		}
	}

	if len(state.Keys) > 0 {
		for i, key := range state.Keys {
			var found bool
			for _, k := range plan.Keys {
				if k.Id.Equal(key.Id) {
					found = true
					break
				}
			}

			if !found {
				r, err := d.client.ProjectsApi.V1RemoveProjectKey(auth, state.Id.ValueString(), key.Id.ValueString()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error removing key from project",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error removing key from project",
						"An unexpected error was encountered trying to remove key from project:\n\n"+string(b),
					)
					return
				}

				// This snippet removes a key from the state on each iteration.
				if len(state.Keys) > 1 {
					if i < len(state.Keys)-1 {
						// Removes the key from the state by copying the slice without the key.
						state.Keys = append(state.Keys[:i], state.Keys[i+1:]...)
					} else {
						// However, if the key is the last one, we can just remove it from the slice.
						state.Keys = state.Keys[:i]
					}
				} else {
					// If the keys attribute is empty, we need to set the state to empty too.
					state.Keys = []V1ProjectKeyModel{}

					// However, when the plan has keys atributes set to nil, we need to set the state to nil too.
					if plan.Keys == nil {
						state.Keys = nil
					}
				}
			}
		}
	}

	if len(plan.Keys) > 0 {
		for i, key := range plan.Keys {
			var found bool
			for _, k := range state.Keys {
				if k.Id.Equal(key.Id) {
					found = true
					break
				}
			}

			if !found {
				p := corellium.NewProjectKey(key.Kind.ValueString(), key.Key.ValueString())
				key, r, err := d.client.ProjectsApi.V1AddProjectKey(auth, state.Id.ValueString()).ProjectKey(*p).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error adding key to project",
							"Coudn't read the response body: "+err.Error(),
						)
						return
					}

					resp.Diagnostics.AddError(
						"Error adding key to project",
						"An unexpected error was encountered trying to add key to project:\n\n"+string(b),
					)
					return
				}

				state.Keys = append(state.Keys, V1ProjectKeyModel{
					Id:          types.StringValue(key.GetIdentifier()),
					Label:       plan.Keys[i].Label,
					Kind:        types.StringValue(key.Kind),
					Key:         types.StringValue(key.Key),
					Fingerprint: types.StringValue(key.GetFingerprint()),
					// CreatedAt:   types.StringValue(time.Time(key.CreatedAt).String()),
					// UpdatedAt:   types.StringValue(time.Time(key.UpdatedAt).String()),
				})
			}
		}
	}

	state.Id = types.StringValue(project.GetId())
	state.Name = types.StringValue(project.GetName())

	state.Settings.Version = types.NumberValue(big.NewFloat(float64(project.Settings.GetVersion())))
	state.Settings.InternetAccess = types.BoolValue(project.Settings.GetInternetAccess())
	state.Settings.Dhcp = types.BoolValue(project.Settings.GetDhcp())

	state.Quotas.Name = types.StringValue(project.GetName())
	state.Quotas.Cores = types.NumberValue(big.NewFloat(float64(project.Quotas.GetCores())))
	state.Quotas.Instances = types.NumberValue(big.NewFloat(float64(project.Quotas.GetInstances())))
	state.Quotas.Ram = types.NumberValue(big.NewFloat(float64(project.Quotas.GetRam())))

	state.UpdatedAt = types.StringValue(time.Now().Format(time.RFC3339))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumV1ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state V1ProjectModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.ProjectsApi.V1DeleteProject(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting project",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete project",
			"An unexpected error was encountered trying to delete the project:\n\n"+string(b))
		return
	}
}

// Configure adds the provider configured client to the resource.
func (d *CorelliumV1ProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
