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
	ID           types.Int64    `tfsdk:"id"`
	Name         types.String   `tfsdk:"name"`
	Key          types.String   `tfsdk:"key"`
	Flavor       types.String   `tfsdk:"flavor"`
	Project      types.Int64    `tfsdk:"project"`
	State        types.String   `tfsdk:"state"`
	StateChanged types.String   `tfsdk:"stateChanged"`
	StartedAt    types.String   `tfsdk:"startedAt"`
	UserTask     types.String   `tfsdk:"userTask"`
	TaskState    types.String   `tfsdk:"taskState"`
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
	ExposePort   types.String   `tfsdk:"exposePort"`
	Fault        types.Bool     `tfsdk:"fault"`
	Patches      []types.String `tfsdk:"patches"`
	CreatedBy    CreatedBy      `tfsdk:"createdBy"`
}

type BootOptions struct {
	BootArgs        types.String   `tfsdk:"bootArgs"`
	RestoreBootArgs types.String   `tfsdk:"restoreBootArgs"`
	UDID            types.Int64    `tfsdk:"udid"`
	ECID            types.String   `tfsdk:"ecid"`
	RandomSeed      types.String   `tfsdk:"randomSeed"`
	PAC             types.Bool     `tfsdk:"pac"`
	APRR            types.Bool     `tfsdk:"aprr"`
	AdditionalTags  []types.String `tfsdk:"additionalTags"`
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
	ID       types.Int64  `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Label    types.String `tfsdk:"label"`
	Deleted  types.Bool   `tfsdk:"deleted"`
}

// Metadata returns the data source type name.
func (d *v1GetInstancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1ready"
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
						"stateChanged": schema.StringAttribute{
							Computed: true,
						},
						"startedAt": schema.StringAttribute{
							Computed: true,
						},
						"userTask": schema.StringAttribute{
							Computed: true,
						},
						"taskState": schema.StringAttribute{
							Computed: true,
						},
						"error": schema.StringAttribute{
							Computed: true,
						},
						"serviceIp": schema.StringAttribute{
							Computed: true,
						},
						"wifiIp": schema.StringAttribute{
							Computed: true,
						},
						"secondaryIp": schema.StringAttribute{
							Computed: true,
						},
						"services": schema.NestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"vpn": schema.NestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"proxy": schema.ListNestedAttribute{
													Computed: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															Type: schema.TypeString,
														},
													},
												},
												"listeners": schema.ListNestedAttribute{
													Computed: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															Type: schema.TypeString,
														},
													},
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
						"netmon": schema.NestedAttribute{
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
						"exposePort": schema.StringAttribute{
							Computed: true,
						},
						"fault": schema.BoolAttribute{
							Computed: true,
						},
						"patches": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									Type: schema.TypeString,
								},
							},
						},
						"createdBy": schema.NestedAttribute{
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
						"bootOptions": schema.NestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"bootArgs": schema.StringAttribute{
										Computed: true,
									},
									"restoreBootArgs": schema.StringAttribute{
										Computed: true,
									},
									"udid": schema.StringAttribute{
										Computed: true,
									},
									"ecid": schema.StringAttribute{
										Computed: true,
									},
									"randomSeed": schema.StringAttribute{
										Computed: true,
									},
									"pac": schema.BoolAttribute{
										Computed: true,
									},
									"aprr": schema.BoolAttribute{
										Computed: true,
									},
									"additionalTags": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												Type: schema.TypeString,
											},
										},
									},
								},
							},
						},
						"agent": schema.NestedAttribute{
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

	instances, r, err := d.client.InstancesApi.V1GetInstances(ctx).Execute()
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
			ID:           types.Int64(instance.ID),
			Name:         types.StringValue(instance.Name),
			Key:          types.StringValue(instance.Key),
			Flavor:       types.StringValue(instance.Flavor),
			Project:      types.Int64(instance.Project),
			State:        types.StringValue(instance.State),
			StateChanged: types.StringValue(instance.StateChanged),
			StartedAt:    types.StringValue(instance.StartedAt),
			UserTask:     types.StringValue(instance.UserTask),
			TaskState:    types.StringValue(instance.TaskState),
			Error:        types.StringValue(instance.Error),
			ServiceIP:    types.StringValue(instance.ServiceIP),
			WifiIP:       types.StringValue(instance.WifiIP),
			SecondaryIP:  types.StringValue(instance.SecondaryIP),
			Panicked:     types.BoolValue(instance.Panicked),
			Created:      types.StringValue(instance.Created),
			Model:        types.StringValue(instance.Model),
			FWPackage:    types.StringValue(instance.FWPackage),
			OS:           types.StringValue(instance.OS),
			ExposePort:   types.StringValue(instance.ExposePort),
			Fault:        types.BoolValue(instance.Fault),
		}

		if instance.BootOptions != nil {
			instanceState.BootOptions = BootOptions{
				BootArgs:        types.StringValue(instance.BootOptions.BootArgs),
				RestoreBootArgs: types.StringValue(instance.BootOptions.RestoreBootArgs),
				UDID:            types.Int64(instance.BootOptions.UDID),
				ECID:            types.StringValue(instance.BootOptions.ECID),
				RandomSeed:      types.StringValue(instance.BootOptions.RandomSeed),
				PAC:             types.BoolValue(instance.BootOptions.PAC),
				APRR:            types.BoolValue(instance.BootOptions.APRR),
			}
			for _, tag := range instance.BootOptions.AdditionalTags {
				instanceState.BootOptions.AdditionalTags = append(instanceState.BootOptions.AdditionalTags, types.StringValue(tag))
			}
		}

		if instance.Services != nil {
			instanceState.Services = Services{}
			if instance.Services.VPN != nil {
				instanceState.Services.VPN = VPN{}
				for _, proxy := range instance.Services.VPN.Proxy {
					instanceState.Services.VPN.Proxy = append(instanceState.Services.VPN.Proxy, Proxy{})
				}
				for _, listener := range instance.Services.VPN.Listeners {
					instanceState.Services.VPN.Listeners = append(instanceState.Services.VPN.Listeners, Listener{})
				}
			}
		}

		if instance.Agent != nil {
			instanceState.Agent = Agent{
				Hash: types.StringValue(instance.Agent.Hash),
				Info: types.StringValue(instance.Agent.Info),
			}
		}

		if instance.Netmon != nil {
			instanceState.Netmon = Netmon{
				Hash:    types.StringValue(instance.Netmon.Hash),
				Info:    types.StringValue(instance.Netmon.Info),
				Enabled: types.BoolValue(instance.Netmon.Enabled),
			}
		}

		if instance.CreatedBy != nil {
			instanceState.CreatedBy = CreatedBy{
				ID:       types.Int64(instance.CreatedBy.ID),
				Username: types.StringValue(instance.CreatedBy.Username),
				Label:    types.StringValue(instance.CreatedBy.Label),
				Deleted:  types.BoolValue(instance.CreatedBy.Deleted),
			}
		}

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
