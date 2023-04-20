package corellium

import (
	"context"

	"terraform-provider-corellium/corellium/pkg/api"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &v1GetInstancesDataSource{}
	_ datasource.DataSourceWithConfigure = &v1GetInstancesDataSource{}
)

// NewCorelliumDataSource is a helper function to simplify the provider implementation.
func NewCorelliumV1GetInstancesSource() datasource.DataSource {
	return &v1GetInstancesDataSource{}
}

// corelliumDataSource is the data source implementation.
type v1GetInstancesDataSource struct {
	client *corellium.APIClient
}

// corelliumDataSourceModel maps the data source schema data.

type v1GetInstancesDataSourceModel struct {
	Instances []v1GetInstanceModel `tfsdk:"instances"`
}

// coffeesModel maps coffees schema data.
type v1GetInstanceModel struct {
	ID           types.String  `tfsdk:"id"`
	Name         types.String  `tfsdk:"name"`
	Key          types.String  `tfsdk:"key"`
	Flavor       types.String  `tfsdk:"flavor"`
	Project      types.String  `tfsdk:"project"`
	State        types.String  `tfsdk:"state"`
	StateChanged types.String  `tfsdk:"state_changed"`
	StartedAt    types.String  `tfsdk:"started_at"`
	UserTask     types.String  `tfsdk:"user_task"`
	TaskState    types.String  `tfsdk:"task_state"`
	Error        types.String  `tfsdk:"error"`
	BootOptions  []BootOptions `tfsdk:"boot_options"`
	ServiceIP    types.String  `tfsdk:"service_ip"`
	WifiIP       types.String  `tfsdk:"wifi_ip"`
	SecondaryIP  types.String  `tfsdk:"secondary_ip"`
	// Services     []Services     `tfsdk:"services"`
	Panicked   types.Bool     `tfsdk:"panicked"`
	Created    types.String   `tfsdk:"created"`
	Model      types.String   `tfsdk:"model"`
	FWPackage  types.String   `tfsdk:"fwpackage"`
	OS         types.String   `tfsdk:"os"`
	Agent      []Agent        `tfsdk:"agent"`
	Netmon     []Netmon       `tfsdk:"netmon"`
	ExposePort types.String   `tfsdk:"expose_port"`
	Fault      types.Bool     `tfsdk:"fault"`
	Patches    []types.String `tfsdk:"patches"`
	CreatedBy  []CreatedBy    `tfsdk:"created_by"`
}

type BootOptions struct {
	BootArgs        types.String   `tfsdk:"boot_args"`
	RestoreBootArgs types.String   `tfsdk:"restore_boot_args"`
	UDID            types.String   `tfsdk:"udid"`
	ECID            types.String   `tfsdk:"ecid"`
	RandomSeed      types.String   `tfsdk:"random_seed"`
	PAC             types.Bool     `tfsdk:"pac"`
	APRR            types.Bool     `tfsdk:"aprr"`
	AdditionalTags  []types.String `tfsdk:"additional_tags"`
}

// type Services struct {
// 	VPN VPN `tfsdk:"vpn"`
// }

//	type VPN struct {
//		Proxy     types.String `tfsdk:"proxy"`
//		Listeners types.String `tfsdk:"listeners"`
//	}
type Agent struct {
	Hash types.String `tfsdk:"hash"`
	Info types.String `tfsdk:"info"`
}

type Netmon struct {
	Hash    types.String `tfsdk:"hash"`
	Info    types.String `tfsdk:"info"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

type CreatedBy struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Label    types.String `tfsdk:"label"`
	Deleted  types.Bool   `tfsdk:"deleted"`
}

// Metadata returns the data source type name.
func (d *v1GetInstancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1getinstances"
}

// Schema defines the schema for the data source.
func (d *v1GetInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Optional: true,
						},
						"name": schema.StringAttribute{
							Optional: true,
						},
						"key": schema.StringAttribute{
							Optional: true,
						},
						"flavor": schema.StringAttribute{
							Optional: true,
						},
						"type": schema.StringAttribute{
							Optional: true,
						},
						"project": schema.StringAttribute{
							Optional: true,
						},
						"state": schema.StringAttribute{
							Optional: true,
						},
						"state_changed": schema.StringAttribute{
							Optional: true,
						},
						"started_at": schema.StringAttribute{
							Optional: true,
						},
						"user_task": schema.StringAttribute{
							Optional: true,
						},
						"task_state": schema.StringAttribute{
							Optional: true,
						},
						"error": schema.StringAttribute{
							Optional: true,
						},
						"service_ip": schema.StringAttribute{
							Optional: true,
						},
						"wifi_ip": schema.StringAttribute{
							Optional: true,
						},
						"secondary_ip": schema.StringAttribute{
							Optional: true,
						},
						// "services": schema.SingleNestedAttribute{
						// 	Optional: true,
						// 	Attributes: map[string]schema.Attribute{
						// 		"vpn": schema.SingleNestedAttribute{
						// 			Optional: true,
						// 			Attributes: map[string]schema.Attribute{
						// 				"proxy": schema.StringAttribute{
						// 					Optional: true,
						// 				},
						// 				"listeners": schema.StringAttribute{
						// 					Optional: true,
						// 				},
						// 			},
						// 		},
						// 	},
						// },
						"panicked": schema.BoolAttribute{
							Optional: true,
						},
						"created": schema.StringAttribute{
							Optional: true,
						},
						"model": schema.StringAttribute{
							Optional: true,
						},
						"fwpackage": schema.StringAttribute{
							Optional: true,
						},
						"os": schema.StringAttribute{
							Optional: true,
						},
						"netmon": schema.SingleNestedAttribute{
							Description: "Netmon",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"hash": schema.StringAttribute{
									Optional: true,
								},
								"info": schema.StringAttribute{
									Optional: true,
								},
								"enabled": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
						"expose_port": schema.StringAttribute{
							Optional: true,
						},
						"fault": schema.BoolAttribute{
							Optional: true,
						},
						"patches": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
						},
						"created_by": schema.SingleNestedAttribute{
							Description: "User who created the instance",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Optional: true,
								},
								"username": schema.StringAttribute{
									Optional: true,
								},
								"label": schema.StringAttribute{
									Optional: true,
								},
								"deleted": schema.BoolAttribute{
									Optional: true,
								},
							},
						},
						"boot_options": schema.SingleNestedAttribute{
							Description: "Boot options for the instance",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"boot_args": schema.StringAttribute{
									Optional: true,
								},
								"restore_boot_args": schema.StringAttribute{
									Optional: true,
								},
								"udid": schema.StringAttribute{
									Optional: true,
								},
								"ecid": schema.StringAttribute{
									Optional: true,
								},
								"random_seed": schema.StringAttribute{
									Optional: true,
								},
								"pac": schema.BoolAttribute{
									Optional: true,
								},
								"aprr": schema.BoolAttribute{
									Optional: true,
								},
								"additional_tags": schema.SetAttribute{
									ElementType: types.StringType,
									Optional:    true,
								},
							},
						},
						"agent": schema.SingleNestedAttribute{
							Description: "Agent information for the instance",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"hash": schema.StringAttribute{
									Optional: true,
								},
								"info": schema.StringAttribute{
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *v1GetInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state v1GetInstancesDataSourceModel

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	instances, _, err := d.client.InstancesApi.V1GetInstances(auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Instances",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, instance := range instances {
		instanceState := v1GetInstanceModel{
			ID:           types.StringValue(instance.GetId()),
			Name:         types.StringValue(instance.GetName()),
			Key:          types.StringValue(instance.GetKey()),
			Flavor:       types.StringValue(instance.GetFlavor()),
			Project:      types.StringValue(instance.GetProject()),
			State:        types.StringValue(string(instance.GetState())),
			StateChanged: types.StringValue(instance.GetStateChanged().GoString()),
			StartedAt:    types.StringValue(instance.GetStartedAt()),
			UserTask:     types.StringValue(instance.GetUserTask()),
			TaskState:    types.StringValue(instance.GetTaskState()),
			Error:        types.StringValue(instance.GetError()),
			ServiceIP:    types.StringValue(instance.GetServiceIp()),
			WifiIP:       types.StringValue(instance.GetWifiIp()),
			SecondaryIP:  types.StringValue(instance.GetSecondaryIp()),
			Panicked:     types.BoolValue(instance.GetPanicked()),
			Created:      types.StringValue(instance.GetCreated().String()),
			Model:        types.StringValue(instance.GetModel()),
			FWPackage:    types.StringValue(instance.GetFwpackage()),
			OS:           types.StringValue(instance.GetOs()),
			ExposePort:   types.StringValue(instance.GetExposePort()),
			Fault:        types.BoolValue(instance.GetFault()),
		}

		// if instance.Services != nil {
		// 	if instance.Services.Vpn != nil {
		// 		instanceState.Services = []Services{
		// 			{
		// 				VPN: []VPN{
		// 					{
		// 						Proxy: types.StringValue(instance.Services.Vpn.GetProxy()),
		// 				},
		// 			},
		// 		}
		// 	}
		// }

		// Convert AdditionalTags from []InstanceBootOptionsAdditionalTag to []basetypes.StringValue
		additionalTags := make([]basetypes.StringValue, len(instance.BootOptions.GetAdditionalTags()))
		for i, tag := range instance.BootOptions.GetAdditionalTags() {
			additionalTags[i] = basetypes.NewStringValue(string(tag)) // Use constructor function to create a new basetypes.StringValue
		}
		if instance.BootOptions != nil {
			instanceState.BootOptions = []BootOptions{
				{
					BootArgs:        types.StringValue(instance.BootOptions.GetBootArgs()),
					RestoreBootArgs: types.StringValue(instance.BootOptions.GetRestoreBootArgs()),
					UDID:            types.StringValue(instance.BootOptions.GetUdid()),
					ECID:            types.StringValue(instance.BootOptions.GetEcid()),
					RandomSeed:      types.StringValue(instance.BootOptions.GetRandomSeed()),
					PAC:             types.BoolValue(instance.BootOptions.GetPac()),
					APRR:            types.BoolValue(instance.BootOptions.GetAprr()),
					AdditionalTags:  additionalTags,

					// AdditionalTags:  []basetypes.StringValue(instance.BootOptions.GetAdditionalTags()),
					// AdditionalTags:  types.StringValue(instance.BootOptions.GetAdditionalTags()),
				},
			}
		}

		if instance.Agent.IsSet() {
			instanceState.Agent = []Agent{
				{
					Hash: types.StringValue(instance.Agent.Get().GetHash()),
					Info: types.StringValue(instance.Agent.Get().GetInfo()),
				},
			}
		}

		if instance.Netmon.IsSet() {
			instanceState.Netmon = []Netmon{
				{
					Hash:    types.StringValue(instance.Netmon.Get().GetHash()),
					Info:    types.StringValue(instance.Netmon.Get().GetInfo()),
					Enabled: types.BoolValue(instance.Netmon.Get().HasEnabled()),
				},
			}
		}

		if instance.CreatedBy != nil {
			instanceState.CreatedBy = []CreatedBy{
				{
					ID:       types.StringValue(instance.CreatedBy.GetId()),
					Username: types.StringValue(instance.CreatedBy.GetUsername()),
					Label:    types.StringValue(instance.CreatedBy.GetLabel()),
					Deleted:  types.BoolValue(instance.CreatedBy.GetDeleted()),
				},
			}
		}

		// instances[index] = instance
		state.Instances = append(state.Instances, instanceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *v1GetInstancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
