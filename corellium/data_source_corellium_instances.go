package corellium

import (
	"context"
	"net/http"

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
	_ datasource.DataSource              = &V1InstancesDataSource{}
	_ datasource.DataSourceWithConfigure = &V1InstancesDataSource{}
)

// NewCorelliumV1InstancesSource is a helper function to simplify the provider implementation.
func NewCorelliumV1InstancesSource() datasource.DataSource {
	return &V1InstancesDataSource{}
}

// V1InstancesDataSource is the data source implementation.
type V1InstancesDataSource struct {
	client *corellium.APIClient
}

// V1InstancesDataSourceModel maps the data source schema data.
type V1InstancesDataSourceModel struct {
	Id        types.String      `tfsdk:"id"`
	Instances []V1InstanceModel `tfsdk:"instances"`
}

// Metadata returns the data source type name.
func (d *V1InstancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1instances"
}

// Schema defines the schema for the data source.
func (d *V1InstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"instances": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Instance id",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Instance name",
							Optional:    true,
							Computed:    true,
						},
						"key": schema.StringAttribute{
							Description: "Instance key",
							Computed:    true,
						},
						"flavor": schema.StringAttribute{
							Description: "Instance flavor",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "Instance type",
							Computed:    true,
						},
						"project": schema.StringAttribute{
							Description: "Instance project",
							Required:    true,
						},
						"state": schema.StringAttribute{
							Description: "Instance state",
							Optional:    true,
							Computed:    true,
							Validators: []validator.String{
								stringvalidator.OneOf("on", "off", "paused"),
							},
						},
						"state_changed": schema.StringAttribute{
							Description: "Instance state changed",
							Computed:    true,
						},
						"started_at": schema.StringAttribute{
							Description: "Instance started at",
							Computed:    true,
						},
						"user_task": schema.StringAttribute{
							Description: "Instance user task",
							Computed:    true,
						},
						"task_state": schema.StringAttribute{
							Description: "Instance task state",
							Computed:    true,
						},
						"error": schema.StringAttribute{
							Description: "Instance error",
							Computed:    true,
						},
						"boot_options": schema.SingleNestedAttribute{
							Description: "Instance boot options",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"boot_args": schema.StringAttribute{
									Description: "Instance boot args",
									Computed:    true,
								},
								"restore_boot_args": schema.StringAttribute{
									Description: "Instance restore boot args",
									Computed:    true,
								},
								"udid": schema.StringAttribute{
									Description: "Instance boot options udid",
									Computed:    true,
								},
								"ecid": schema.StringAttribute{
									Description: "Instance boot options ecid",
									Computed:    true,
								},
								"random_seed": schema.StringAttribute{
									Description: "Instance boot options random seed",
									Computed:    true,
								},
								"pac": schema.BoolAttribute{
									Description: "Instance boot options pac",
									Computed:    true,
								},
								"aprr": schema.BoolAttribute{
									Description: "Instance boot options aprr",
									Computed:    true,
								},
								"additional_tags": schema.ListAttribute{
									// TODO: add validation to this list.
									Description: "Instance boot options additional tags",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"service_ip": schema.StringAttribute{
							Description: "Instance service ip",
							Computed:    true,
						},
						"wifi_ip": schema.StringAttribute{
							Description: "Instance wifi ip",
							Computed:    true,
						},
						"secondary_ip": schema.StringAttribute{
							Description: "Instance secondary ip",
							Computed:    true,
						},
						/*"services": schema.SingleNestedAttribute{ // TODO: find the right type for this.
							Description: "Instance services",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"vpn": schema.SingleNestedAttribute{
									Description: "Instance services vpn",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"proxy": schema.MapAttribute{
											Computed: true,
										},
										"listeners": schema.MapAttribute{
											Computed: true,
										},
									},
								},
							},
						},*/
						"panicked": schema.BoolAttribute{
							Description: "Instance panicked",
							Computed:    true,
						},
						"created": schema.StringAttribute{
							Description: "Instance created",
							Computed:    true,
						},
						"model": schema.StringAttribute{
							Description: "Instance model",
							Computed:    true,
						},
						"fwpackage": schema.StringAttribute{
							Description: "Instance fwpackage",
							Computed:    true,
						},
						"os": schema.StringAttribute{
							Description: "Instance os",
							Required:    true,
						},
						"agent": schema.SingleNestedAttribute{
							Description: "Instance agent",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"hash": schema.StringAttribute{
									Description: "Instance agent hash",
									Computed:    true,
								},
								"info": schema.StringAttribute{
									Description: "Instance agent info",
									Computed:    true,
								},
							},
						},
						"netmon": schema.SingleNestedAttribute{
							Description: "Instance agent",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"hash": schema.StringAttribute{
									Description: "Instance netmon hash",
									Computed:    true,
								},
								"info": schema.StringAttribute{
									Description: "Instance netmon info",
									Computed:    true,
								},
								"enabled": schema.BoolAttribute{
									Description: "Instance netmon enabled",
									Computed:    true,
								},
							},
						},
						"expose_port": schema.StringAttribute{
							Description: "Instance expose port",
							Computed:    true,
						},
						"fault": schema.BoolAttribute{
							Description: "Instance fault",
							Computed:    true,
						},
						"patches": schema.ListAttribute{
							Description: "Instance patches",
							Computed:    true,
							ElementType: types.StringType,
						},
						"created_by": schema.SingleNestedAttribute{
							Description: "Instance created by",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description: "Instance user id",
									Computed:    true,
								},
								"username": schema.StringAttribute{
									Description: "Instance user username",
									Computed:    true,
								},
								"label": schema.StringAttribute{
									Description: "Instance user label",
									Computed:    true,
								},
								"deleted": schema.BoolAttribute{
									Description: "Instance user deleted status",
									Computed:    true,
								},
							},
						},
						"wait_for_ready": schema.BoolAttribute{
							Description: "Wait for ready",
							Optional:    true,
						},
						"wait_for_ready_timeout": schema.Int64Attribute{
							Description: "Wait for ready timeout",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *V1InstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state V1InstancesDataSourceModel

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	instances, r, err := d.client.InstancesApi.V1GetInstances(auth).Execute()
	if err != nil {
		if r.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Unable to Read Instances",
				"You do not have permission to access instances",
			)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to Read Instances",
			err.Error(),
		)
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error generating UUID",
			"An unexpected error was encountered trying to generate UUID\n\n"+err.Error(),
		)
		return
	}

	state.Id = types.StringValue(id)
	state.Instances = make([]V1InstanceModel, len(instances))
	for i, instance := range instances {
		state.Instances[i].Id = types.StringValue(instance.GetId())
		state.Instances[i].Name = types.StringValue(instance.GetName())
		state.Instances[i].Key = types.StringValue(instance.GetKey())
		state.Instances[i].Flavor = types.StringValue(instance.GetFlavor())
		state.Instances[i].Type = types.StringValue(instance.GetType())
		state.Instances[i].Project = types.StringValue(instance.GetProject())
		state.Instances[i].State = types.StringValue(string(instance.GetState()))
		state.Instances[i].StateChanged = types.StringValue(instance.GetStateChanged().UTC().String())
		state.Instances[i].StartedAt = types.StringValue(instance.GetStartedAt())
		state.Instances[i].UserTask = types.StringValue(instance.GetUserTask())
		state.Instances[i].TaskState = types.StringValue(instance.GetTaskState())
		state.Instances[i].Error = types.StringValue(instance.GetError())

		additionalTags, diags := types.ListValueFrom(ctx, types.StringType, instance.BootOptions.GetAdditionalTags())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		state.Instances[i].BootOptions = &V1InstanceBootOptionsModel{
			BootArgs:        types.StringValue(instance.BootOptions.GetBootArgs()),
			RestoreBootArgs: types.StringValue(instance.BootOptions.GetRestoreBootArgs()),
			UDID:            types.StringValue(instance.BootOptions.GetUdid()),
			ECID:            types.StringValue(instance.BootOptions.GetEcid()),
			RandomSeed:      types.StringValue(instance.BootOptions.GetRandomSeed()),
			PAC:             types.BoolValue(instance.BootOptions.GetPac()),
			APRR:            types.BoolValue(instance.BootOptions.GetAprr()),
			AdditionalTags:  additionalTags,
		}

		state.Instances[i].ServiceIP = types.StringValue(instance.GetServiceIp())
		state.Instances[i].WifiIP = types.StringValue(instance.GetWifiIp())
		state.Instances[i].SecondaryIP = types.StringValue(instance.GetSecondaryIp())

		/*proxy, diags := types.MapValueFrom(ctx, types.StringType, instance.Services.GetVpn().Proxy) // TODO: find the right type for this.
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		listeners, diags := types.MapValueFrom(ctx, types.StringType, instance.Services.GetVpn().Listeners)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		state.Services = &V1InstanceServicesModel{
			VPN: &V1InstanceVPNModel{
				Proxy:     proxy,
				Listeners: listeners,
			},
		}*/

		state.Instances[i].Panicked = types.BoolValue(instance.GetPanicked())
		state.Instances[i].Created = types.StringValue(instance.GetCreated().UTC().String())
		state.Instances[i].Model = types.StringValue(instance.GetModel())
		state.Instances[i].FWPackage = types.StringValue(instance.GetFwpackage())
		state.Instances[i].OS = types.StringValue(instance.GetOs())
		state.Instances[i].Agent = &V1InstanceAgentModel{
			Hash: types.StringValue(instance.Agent.Get().GetHash()),
			Info: types.StringValue(instance.Agent.Get().GetInfo()),
		}
		state.Instances[i].Netmon = &V1InstanceNetmonModel{
			Hash:    types.StringValue(instance.Netmon.Get().GetHash()),
			Info:    types.StringValue(instance.Netmon.Get().GetInfo()),
			Enabled: types.BoolValue(instance.Netmon.Get().GetEnabled()),
		}
		state.Instances[i].ExposePort = types.StringValue(instance.GetExposePort())
		state.Instances[i].Fault = types.BoolValue(instance.GetFault())

		patches, diags := types.ListValueFrom(ctx, types.StringType, instance.GetPatches())
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		state.Instances[i].Patches = patches

		state.Instances[i].CreatedBy = &V1InstanceCreatedByModel{
			Id:       types.StringValue(instance.CreatedBy.GetId()),
			Username: types.StringValue(instance.CreatedBy.GetUsername()),
			Label:    types.StringValue(instance.CreatedBy.GetLabel()),
			Deleted:  types.BoolValue(instance.CreatedBy.GetDeleted()),
		}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *V1InstancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
