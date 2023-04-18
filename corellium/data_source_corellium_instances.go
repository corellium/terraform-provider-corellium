package corellium

import (
	"context"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	ID           types.String   `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Key          types.String   `tfsdk:"key"`
	Flavor       types.String   `tfsdk:"flavor"`
	Project      types.String   `tfsdk:"project"`
	State        types.String   `tfsdk:"state"`
	StateChanged types.String   `tfsdk:"state_changed"`
	StartedAt    types.String   `tfsdk:"started_at"`
	UserTask     types.String   `tfsdk:"user_task"`
	TaskState    types.String   `tfsdk:"task_state"`
	Error        types.String   `tfsdk:"error"`
	BootOptions  BootOptions    `tfsdk:"bootOptions"`
	ServiceIP    types.String   `tfsdk:"serviceIp"`
	WifiIP       types.String   `tfsdk:"wifiIp"`
	SecondaryIP  types.String   `tfsdk:"secondaryIp"`
	Services     Services       `tfsdk:"services"`
	Panicked     types.Bool     `tfsdk:"panicked"`
	Created      types.String   `tfsdk:"created"`
	Model        types.String   `tfsdk:"model"`
	FWPackage    types.String   `tfsdk:"fwpackage"`
	OS           types.String   `tfsdk:"os"`
	Agent        Agent          `tfsdk:"agent"`
	Netmon       Netmon         `tfsdk:"netmon"`
	ExposePort   types.String   `tfsdk:"expose_port"`
	Fault        types.Bool     `tfsdk:"fault"`
	Patches      []types.String `tfsdk:"patches"`
	CreatedBy    CreatedBy      `tfsdk:"created_by"`
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

type Services struct {
	VPN VPN `tfsdk:"vpn"`
}

type VPN struct {
	Proxy     []Proxy    `tfsdk:"proxy"`
	Listeners []Listener `tfsdk:"listeners"`
}

type Proxy struct {
	// Proxy configuration attributes
}

type Listener struct {
	// Listener configuration attributes
}

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
	resp.TypeName = req.ProviderTypeName + "_v1instances"
}

// Schema defines the schema for the data source.
func (d *v1GetInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"key": schema.StringAttribute{
							Computed: true,
						},
						"flavor": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"project": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"state_changed": schema.StringAttribute{
							Computed: true,
						},
						"started_at": schema.StringAttribute{
							Computed: true,
						},
						"user_task": schema.StringAttribute{
							Computed: true,
						},
						"task_state": schema.StringAttribute{
							Computed: true,
						},
						"error": schema.StringAttribute{
							Computed: true,
						},
						"service_ip": schema.StringAttribute{
							Computed: true,
						},
						"wifi_ip": schema.StringAttribute{
							Computed: true,
						},
						"secondary_ip": schema.StringAttribute{
							Computed: true,
						},
						"services": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"vpn": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"proxy": schema.StringAttribute{
													Computed: true,
												},
												"listeners": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"panicked": schema.BoolAttribute{
							Computed: true,
						},
						"created": schema.StringAttribute{
							Computed: true,
						},
						"model": schema.StringAttribute{
							Computed: true,
						},
						"fwpackage": schema.StringAttribute{
							Computed: true,
						},
						"os": schema.StringAttribute{
							Computed: true,
						},
						"netmon": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"hash": schema.StringAttribute{
										Computed: true,
									},
									"info": schema.StringAttribute{
										Computed: true,
									},
									"enabled": schema.BoolAttribute{
										Computed: true,
									},
								},
							},
						},
						"expose_port": schema.StringAttribute{
							Computed: true,
						},
						"fault": schema.BoolAttribute{
							Computed: true,
						},
						"patches": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									// TODO: Patches string here
								},
							},
						},
						"created_by": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"username": schema.StringAttribute{
										Computed: true,
									},
									"label": schema.StringAttribute{
										Computed: true,
									},
									"deleted": schema.BoolAttribute{
										Computed: true,
									},
								},
							},
						},
						"boot_options": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"boot_args": schema.StringAttribute{
										Computed: true,
									},
									"restore_boot_args": schema.StringAttribute{
										Computed: true,
									},
									"udid": schema.StringAttribute{
										Computed: true,
									},
									"ecid": schema.StringAttribute{
										Computed: true,
									},
									"random_seed": schema.StringAttribute{
										Computed: true,
									},
									"pac": schema.BoolAttribute{
										Computed: true,
									},
									"aprr": schema.BoolAttribute{
										Computed: true,
									},
									"additional_tags": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												// TODO: additional_tags string
											},
										},
									},
								},
							},
						},
						"agent": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"hash": schema.StringAttribute{
										Computed: true,
									},
									"info": schema.StringAttribute{
										Computed: true,
									},
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
	var state v1GetInstanceModel

	instances, _, err := d.client.InstancesApi.V1GetInstances(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Instances",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for index, instance := range instances {
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

		if instance.BootOptions != nil {
			instanceState.BootOptions = BootOptions{
				BootArgs:        types.StringValue(instance.BootOptions.GetBootArgs()),
				RestoreBootArgs: types.StringValue(instance.BootOptions.GetRestoreBootArgs()),
				UDID:            types.StringValue(instance.BootOptions.GetUdid()),
				ECID:            types.StringValue(instance.BootOptions.GetEcid()),
				RandomSeed:      types.StringValue(instance.BootOptions.GetRandomSeed()),
				PAC:             types.BoolValue(instance.BootOptions.GetPac()),
				APRR:            types.BoolValue(instance.BootOptions.GetAprr()),
			}
			//for _, tag := range instance.BootOptions.AdditionalTags {
			//	instanceState.BootOptions.AdditionalTags = append(instanceState.BootOptions.AdditionalTags, types.StringValue(tag))
			//}
		}

		if instance.Services != nil {
			// instanceState.Services = Services{}

			if instance.Services.Vpn != nil {
				// instanceState.Services.VPN = VPN{}

				//for i, proxy := range instance.Services.Vpn.GetProxy() {
				//	instance.Services.Vpn.Listeners[i]
				//}

				/*for _, proxy := range instance.Services.VPN.Proxy {
					instanceState.Services.VPN.Proxy = append(instanceState.Services.VPN.Proxy, Proxy{})
				}
				for _, listener := range instance.Services.VPN.Listeners {
					instanceState.Services.VPN.Listeners = append(instanceState.Services.VPN.Listeners, Listener{})
				}*/
			}
		}

		if instance.Agent.IsSet() {
			instanceState.Agent = Agent{
				Hash: types.StringValue(instance.Agent.Get().GetHash()),
				Info: types.StringValue(instance.Agent.Get().GetInfo()),
			}
		}

		if instance.Netmon.IsSet() {
			instanceState.Netmon = Netmon{
				Hash:    types.StringValue(instance.Netmon.Get().GetHash()),
				Info:    types.StringValue(instance.Netmon.Get().GetInfo()),
				Enabled: types.BoolValue(instance.Netmon.Get().HasEnabled()),
			}
		}

		if instance.CreatedBy != nil {
			instanceState.CreatedBy = CreatedBy{
				ID:       types.StringValue(instance.CreatedBy.GetId()),
				Username: types.StringValue(instance.CreatedBy.GetUsername()),
				Label:    types.StringValue(instance.CreatedBy.GetLabel()),
				Deleted:  types.BoolValue(instance.CreatedBy.GetDeleted()),
			}
		}

		instances[index] = instance
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
