package corellium

import (
	"context"
	"io"
	"time"

	"github.com/aimoda/go-corellium-api-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"terraform-provider-corellium/corellium/pkg/api"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumV1InstanceResource{}
	_ resource.ResourceWithConfigure = &CorelliumV1InstanceResource{}
)

// NewCorelliumV1InstanceResource is a helper function to simplify the provider implementation.
func NewCorelliumV1InstanceResource() resource.Resource {
	return &CorelliumV1InstanceResource{}
}

// CorelliumV1InstanceResource is the resource implementation.
type CorelliumV1InstanceResource struct {
	client *corellium.APIClient
}

type V1InstanceVPNModel struct {
	Proxy     types.Map `tfsdk:"proxy"`
	Listeners types.Map `tfsdk:"listeners"`
}

type V1InstanceServicesModel struct {
	VPN *V1InstanceVPNModel `tfsdk:"vpn"`
}

type V1InstanceAgentModel struct {
	Hash types.String `tfsdk:"hash"`
	Info types.String `tfsdk:"info"`
}

type V1InstanceNetmonModel struct {
	Hash    types.String `tfsdk:"hash"`
	Info    types.String `tfsdk:"info"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

type V1InstanceCreatedByModel struct { // TODO: rename to user or use the user model.
	// Id is the user who created the instance.
	Id types.String `tfsdk:"id"`
	// Username is the username of the user who created the instance.
	Username types.String `tfsdk:"username"`
	// Label is the label of the user who created the instance.
	Label types.String `tfsdk:"label"`
	// Deleted is the deleted status of the user who created the instance.
	Deleted types.Bool `tfsdk:"deleted"`
}

type V1InstanceBootOptionsModel struct {
	BootArgs        types.String `tfsdk:"boot_args"`
	RestoreBootArgs types.String `tfsdk:"restore_boot_args"`
	UDID            types.String `tfsdk:"udid"`
	ECID            types.String `tfsdk:"ecid"`
	RandomSeed      types.String `tfsdk:"random_seed"`
	PAC             types.Bool   `tfsdk:"pac"`
	APRR            types.Bool   `tfsdk:"aprr"`
	// AdditionalTags is the additional tags of the instance.
	// AdditionalTags can assume the following values:
	// kalloc: Enable kalloc/kfree trace access via GDB (Enterprise only)
	// gpu: Enable cloud GPU acceleration (Extra costs incurred, cloud only)
	// no-keyboard: Enable keyboard passthrough from web interface
	// nodevmode: Disable developer mode on iOS16 and greater
	// sep-cons-ext: Patch SEPOS to print debug messages to console
	// iboot-jailbreak: Patch iBoot to disable signature checks
	// llb-jailbreak: Patch LLB to disable signature checks
	// rom-jailbreak: Patch BootROM to disable signature checks
	AdditionalTags types.List `tfsdk:"additional_tags"`
}

const (
	// V1InstanceStateOn is the state of the instance when it is running.
	V1InstanceStateOn = "on"
	// V1InstanceStateOff is the state of the instance when it is powered off.
	V1InstanceStateOff = "off"
	// V1InstanceStatePaused is the state of the instance when it is paused.
	V1InstanceStatePaused = "paused"
	// V1InstanceStateCreating is the state of the instance when it is being created.
	V1InstanceStateCreating = "creating"
	// V1InstanceStateDeleting is the state of the instance when it is being deleted.
	V1InstanceStateDeleting = "deleting"
)

const (
	// V1InstancesTaskStateNone is the state of the instance when there is no task.
	V1InstancesTaskStateNone = "none"
	// V1InstancesTaskStateBuilding is the state of the instance when it is being built.
	V1InstancesTaskStateBuilding = "building"
)

// coffeesModel maps coffees schema data.
type V1InstanceModel struct {
	// Id is the unique identifier of the instance.
	Id types.String `tfsdk:"id"`
	// Name is the name of the instance.
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
	// Flavor is the flavor of the instance.
	// A flavor is a device model, what can be a Android or iOS device.
	// The following flavors are examples of supported flavors for Android:
	//    - ranchu (for Generic Android devices)
	//    - google-nexus-4
	//    - google-nexus-5
	//    - google-nexus-5x
	//    - google-nexus-6
	//    - google-nexus-6p
	//    - google-nexus-9
	//    - google-pixel
	//    - google-pixel-2
	//    - google-pixel-3
	//    - htc-one-m8
	//    - huawei-p8
	//    - samsung-galaxy-s-duos 
	//
	// The following flavors are examples for iOS:
	//    - iphone6
	//    - iphone6plus
	//    - ipodtouch6
	//    - ipadmini4wifi
	//    - iphone6s
	//    - iphone6splus
	//    - iphonese
	//    - iphone7
	//    - iphone7plus
	//    - iphone8
	//    - iphone8plus
	//    - iphonex
	//    - iphonexs
	//    - iphonexsmax
	//    - iphonexsmaxww
	//    - iphonexr
	//    - iphone11
	//    - iphone11pro
	//    - iphone11promax
	//    - iphonese2
	//    - iphone12m
	//    - iphone12
	//    - iphone12p
	//    - iphone12pm
	//    - iphone13
	//    - iphone13m
	//    - iphone13p
	//    - iphone13pm
	Flavor types.String `tfsdk:"flavor"`
	Type   types.String `tfsdk:"type"`
	// Project is the id of the project the instance belongs to.
	Project types.String `tfsdk:"project"`
	// State is the current state of the instance.
	// State can assume the following values:
	// on - The instance is running.
	// off - The instance is powered off.
	// paused - The instance is paused.
	// creating - The instance is being created.
	// deleting - The instance is being deleted.
	State        types.String `tfsdk:"state"`
	StateChanged types.String `tfsdk:"state_changed"`
	// StartedAt is the time the instance was started.
	StartedAt types.String `tfsdk:"started_at"`
	UserTask  types.String `tfsdk:"user_task"`
	TaskState types.String `tfsdk:"task_state"`
	// Error is the error message of the instance.
	Error types.String `tfsdk:"error"`
	// BootOptions is the boot options of the instance.
	BootOptions *V1InstanceBootOptionsModel `tfsdk:"boot_options"`
	ServiceIP   types.String                `tfsdk:"service_ip"`
	WifiIP      types.String                `tfsdk:"wifi_ip"`
	SecondaryIP types.String                `tfsdk:"secondary_ip"`
	// Services    *V1InstanceServicesModel    `tfsdk:"services"` // TODO: find the right type for this.
	Panicked types.Bool `tfsdk:"panicked"`
	// Created is the time the instance was created.
	Created   types.String `tfsdk:"created"`
	Model     types.String `tfsdk:"model"`
	FWPackage types.String `tfsdk:"fwpackage"`
	// OS is the version of the operating system running on the instance, e.g. 14.3 for iOS, or 11.0.0 for Android.
	OS         types.String           `tfsdk:"os"`
	Agent      *V1InstanceAgentModel  `tfsdk:"agent"`
	Netmon     *V1InstanceNetmonModel `tfsdk:"netmon"`
	ExposePort types.String           `tfsdk:"expose_port"`
	Fault      types.Bool             `tfsdk:"fault"`
	// Patches is the list of patches applied to the instance.
	//    - jailbroken The instance should be jailbroken (default).
	//    - nonjailbroken The instance should not be jailbroken.
	//    - corelliumd The instance should not be jailbroken but should profile API agent.
	Patches types.List `tfsdk:"patches"`
	// CreatedBy is the user who created the instance.
	CreatedBy *V1InstanceCreatedByModel `tfsdk:"created_by"`
	// WaitForReady is a boolean that indicates if the resource should wait for the instance to be ready.
	WaitForReady types.Bool `tfsdk:"wait_for_ready"`
	// WaitForReadyTimeout is the timeout in seconds to wait for the instance to be ready.
	// WaitForReadyTimeout is a amount in seconds.
	WaitForReadyTimeout types.Int64 `tfsdk:"wait_for_ready_timeout"`
}

// Metadata returns the resource type name.
func (d *CorelliumV1InstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_v1instance"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_v1instance "instance" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumV1InstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
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
				Required: true,
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
				Default:     nil,
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
	}
}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumV1InstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan V1InstanceModel

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	i := corellium.NewInstanceCreateOptions(plan.Flavor.ValueString(), plan.Project.ValueString(), plan.OS.ValueString())
	i.SetName(plan.Name.ValueString())
	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	created, r, err := d.client.InstancesApi.V1CreateInstance(auth).InstanceCreateOptions(*i).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating instance",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error creating instance",
			"An unexpected error was encountered trying to create the instance:\n\n"+string(b),
		)
		return
	}

	plan.Id = types.StringValue(created.GetId())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.WaitForReady.IsUnknown() && plan.WaitForReady.ValueBool() {
		createStateConf := &retry.StateChangeConf{
			Refresh: func() (interface{}, string, error) {
				instance, r, err := d.client.InstancesApi.V1GetInstance(auth, created.GetId()).Execute()
				if err != nil {
					b, err := io.ReadAll(r.Body)
					if err != nil {
						resp.Diagnostics.AddError(
							"Error creating instance",
							"Coudn't read the response body when checking if instnace is ready: "+err.Error(),
						)
						return nil, "", err
					}

					resp.Diagnostics.AddError(
						"Error creating snapshot",
						"An unexpected error was encountered trying to check if instance is ready:\n\n"+string(b),
					)
					return nil, "", err
				}

				return instance, string(instance.GetState()), nil
			},
			Pending: []string{
				V1InstanceStateCreating,
				V1InstanceStateDeleting,
			},
			Target: []string{
				V1InstanceStateOn,
				V1InstanceStateOff,
				V1InstanceStatePaused,
			},
			Delay:      5 * time.Second,
			MinTimeout: 5 * time.Second,
			Timeout: func() time.Duration {
				if plan.WaitForReadyTimeout.IsUnknown() {
					return 300
				}

				return time.Duration(plan.WaitForReadyTimeout.ValueInt64())
			}() * time.Second,
		}

		if _, err := createStateConf.WaitForStateContext(ctx); err != nil {
			resp.Diagnostics.AddError(
				"Error creating snapshot",
				"Coudn't create the snapshot: "+err.Error(),
			)

			return
		}
	}

	instance, r, err := d.client.InstancesApi.V1GetInstance(auth, created.GetId()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get instance",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error get instance",
			"An unexpected error was encountered trying to get the instance:\n\n"+string(b),
		)
		return
	}

	plan.Name = types.StringValue(instance.GetName())
	plan.Key = types.StringValue(instance.GetKey())
	plan.Flavor = types.StringValue(instance.GetFlavor())
	plan.Type = types.StringValue(instance.GetType())
	plan.Project = types.StringValue(instance.GetProject())
	plan.State = types.StringValue(string(instance.GetState()))
	plan.StateChanged = types.StringValue(instance.GetStateChanged().UTC().String())
	plan.StartedAt = types.StringValue(instance.GetStartedAt())
	plan.UserTask = types.StringValue(instance.GetUserTask())
	plan.TaskState = types.StringValue(instance.GetTaskState())
	plan.Error = types.StringValue(instance.GetError())

	additionalTags, diags := types.ListValueFrom(ctx, types.StringType, instance.BootOptions.GetAdditionalTags())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.BootOptions = &V1InstanceBootOptionsModel{
		BootArgs:        types.StringValue(instance.BootOptions.GetBootArgs()),
		RestoreBootArgs: types.StringValue(instance.BootOptions.GetRestoreBootArgs()),
		UDID:            types.StringValue(instance.BootOptions.GetUdid()),
		ECID:            types.StringValue(instance.BootOptions.GetEcid()),
		RandomSeed:      types.StringValue(instance.BootOptions.GetRandomSeed()),
		PAC:             types.BoolValue(instance.BootOptions.GetPac()),
		APRR:            types.BoolValue(instance.BootOptions.GetAprr()),
		AdditionalTags:  additionalTags,
	}

	/*plan.ServiceIP = types.StringValue(instance.GetServiceIp()) // TODO: find the right type for this.
	plan.WifiIP = types.StringValue(instance.GetWifiIp())
	plan.SecondaryIP = types.StringValue(instance.GetSecondaryIp())


	proxy, diags := types.MapValueFrom(ctx, types.StringType, instance.Services.GetVpn().Proxy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listeners, diags := types.MapValueFrom(ctx, types.StringType, instance.Services.GetVpn().Listeners)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Services = &V1InstanceServicesModel{
		VPN: &V1InstanceVPNModel{
			Proxy:     proxy,
			Listeners: listeners,
		},
	}*/

	plan.Panicked = types.BoolValue(instance.GetPanicked())
	plan.Created = types.StringValue(instance.GetCreated().UTC().String())
	plan.Model = types.StringValue(instance.GetModel())
	plan.FWPackage = types.StringValue(instance.GetFwpackage())
	plan.OS = types.StringValue(instance.GetOs())
	plan.Agent = &V1InstanceAgentModel{
		Hash: types.StringValue(instance.Agent.Get().GetHash()),
		Info: types.StringValue(instance.Agent.Get().GetInfo()),
	}
	plan.Netmon = &V1InstanceNetmonModel{
		Hash:    types.StringValue(instance.Netmon.Get().GetHash()),
		Info:    types.StringValue(instance.Netmon.Get().GetInfo()),
		Enabled: types.BoolValue(instance.Netmon.Get().GetEnabled()),
	}
	plan.ExposePort = types.StringValue(instance.GetExposePort())
	plan.Fault = types.BoolValue(instance.GetFault())

	patches, diags := types.ListValueFrom(ctx, types.StringType, instance.GetPatches())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Patches = patches

	plan.CreatedBy = &V1InstanceCreatedByModel{
		Id:       types.StringValue(instance.CreatedBy.GetId()),
		Username: types.StringValue(instance.CreatedBy.GetUsername()),
		Label:    types.StringValue(instance.CreatedBy.GetLabel()),
		Deleted:  types.BoolValue(instance.CreatedBy.GetDeleted()),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumV1InstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state V1InstanceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	instance, r, err := d.client.InstancesApi.V1GetInstance(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error get the instance",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error to get the instance",
			"An unexpected error was encountered trying to get the instance:\n\n"+string(b),
		)
		return
	}

	state.Id = types.StringValue(instance.GetId())
	state.Name = types.StringValue(instance.GetName())
	state.Key = types.StringValue(instance.GetKey())
	state.Flavor = types.StringValue(instance.GetFlavor())
	state.Type = types.StringValue(instance.GetType())
	state.Project = types.StringValue(instance.GetProject())
	state.State = types.StringValue(string(instance.GetState()))
	state.StateChanged = types.StringValue(instance.GetStateChanged().UTC().String())
	state.StartedAt = types.StringValue(instance.GetStartedAt())
	state.UserTask = types.StringValue(instance.GetUserTask())
	state.TaskState = types.StringValue(instance.GetTaskState())
	state.Error = types.StringValue(instance.GetError())

	additionalTags, diags := types.ListValueFrom(ctx, types.StringType, instance.BootOptions.GetAdditionalTags())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.BootOptions = &V1InstanceBootOptionsModel{
		BootArgs:        types.StringValue(instance.BootOptions.GetBootArgs()),
		RestoreBootArgs: types.StringValue(instance.BootOptions.GetRestoreBootArgs()),
		UDID:            types.StringValue(instance.BootOptions.GetUdid()),
		ECID:            types.StringValue(instance.BootOptions.GetEcid()),
		RandomSeed:      types.StringValue(instance.BootOptions.GetRandomSeed()),
		PAC:             types.BoolValue(instance.BootOptions.GetPac()),
		APRR:            types.BoolValue(instance.BootOptions.GetAprr()),
		AdditionalTags:  additionalTags,
	}

	state.ServiceIP = types.StringValue(instance.GetServiceIp())
	state.WifiIP = types.StringValue(instance.GetWifiIp())
	state.SecondaryIP = types.StringValue(instance.GetSecondaryIp())

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

	state.Panicked = types.BoolValue(instance.GetPanicked())
	state.Created = types.StringValue(instance.GetCreated().UTC().String())
	state.Model = types.StringValue(instance.GetModel())
	state.FWPackage = types.StringValue(instance.GetFwpackage())
	state.OS = types.StringValue(instance.GetOs())
	state.Agent = &V1InstanceAgentModel{
		Hash: types.StringValue(instance.Agent.Get().GetHash()),
		Info: types.StringValue(instance.Agent.Get().GetInfo()),
	}
	state.Netmon = &V1InstanceNetmonModel{
		Hash:    types.StringValue(instance.Netmon.Get().GetHash()),
		Info:    types.StringValue(instance.Netmon.Get().GetInfo()),
		Enabled: types.BoolValue(instance.Netmon.Get().GetEnabled()),
	}
	state.ExposePort = types.StringValue(instance.GetExposePort())
	state.Fault = types.BoolValue(instance.GetFault())

	patches, diags := types.ListValueFrom(ctx, types.StringType, instance.GetPatches())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Patches = patches

	state.CreatedBy = &V1InstanceCreatedByModel{
		Id:       types.StringValue(instance.CreatedBy.GetId()),
		Username: types.StringValue(instance.CreatedBy.GetUsername()),
		Label:    types.StringValue(instance.CreatedBy.GetLabel()),
		Deleted:  types.BoolValue(instance.CreatedBy.GetDeleted()),
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumV1InstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state V1InstanceModel
	// state is the current state of the resource.

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan V1InstanceModel
	// plan is the proposed new state of the resource.

	diags = req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.Flavor.Equal(plan.Flavor) {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"Flavor cannot be changed",
		)
		return
	}

	if !state.Project.Equal(plan.Project) {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"Project cannot be changed",
		)
		return
	}

	if !state.OS.Equal(plan.OS) {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"OS cannot be changed",
		)
		return
	}

	p := corellium.NewPatchInstanceOptions()
	p.SetName(plan.Name.ValueString())

	if state.State.ValueString() == V1InstanceStateCreating {
		resp.Diagnostics.AddWarning(
			"Warning updating instance",
			"Instance is currently being created. Changes to its state cannot be applied.",
		)
	}

	if state.State.ValueString() == V1InstanceStateDeleting {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"Instance is currently being deleted. Changes cannot be applied.",
		)
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	instance, r, err := d.client.InstancesApi.V1PatchInstance(auth, state.Id.ValueString()).PatchInstanceOptions(*p).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating instance",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error updating instance",
			"An unexpected error was encountered trying to update the instance:\n\n"+string(b),
		)
		return
	}

	state.Id = types.StringValue(instance.GetId())
	state.Name = types.StringValue(instance.GetName())
	state.Key = types.StringValue(instance.GetKey())
	state.Flavor = types.StringValue(instance.GetFlavor())
	state.Type = types.StringValue(instance.GetType())
	state.Project = types.StringValue(instance.GetProject())
	state.State = types.StringValue(string(instance.GetState()))
	state.StateChanged = types.StringValue(instance.GetStateChanged().UTC().String())
	state.StartedAt = types.StringValue(instance.GetStartedAt())
	state.UserTask = types.StringValue(instance.GetUserTask())
	state.TaskState = types.StringValue(instance.GetTaskState())
	state.Error = types.StringValue(instance.GetError())

	state.BootOptions.BootArgs = types.StringValue(instance.BootOptions.GetBootArgs())
	state.BootOptions.RestoreBootArgs = types.StringValue(instance.BootOptions.GetRestoreBootArgs())
	state.BootOptions.UDID = types.StringValue(instance.BootOptions.GetUdid())
	state.BootOptions.ECID = types.StringValue(instance.BootOptions.GetEcid())
	state.BootOptions.RandomSeed = types.StringValue(instance.BootOptions.GetRandomSeed())
	state.BootOptions.PAC = types.BoolValue(instance.BootOptions.GetPac())
	state.BootOptions.APRR = types.BoolValue(instance.BootOptions.GetAprr())

	additionalTags, diags := types.ListValueFrom(ctx, types.StringType, instance.BootOptions.GetAdditionalTags())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.BootOptions.AdditionalTags = additionalTags

	state.ServiceIP = types.StringValue(instance.GetServiceIp())
	state.WifiIP = types.StringValue(instance.GetWifiIp())
	state.SecondaryIP = types.StringValue(instance.GetSecondaryIp())

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

	state.Panicked = types.BoolValue(instance.GetPanicked())
	state.Created = types.StringValue(instance.GetCreated().UTC().String())
	state.Model = types.StringValue(instance.GetModel())
	state.FWPackage = types.StringValue(instance.GetFwpackage())
	state.OS = types.StringValue(instance.GetOs())

	state.Agent.Hash = types.StringValue(instance.Agent.Get().GetHash())
	state.Agent.Info = types.StringValue(instance.Agent.Get().GetInfo())

	state.Netmon.Hash = types.StringValue(instance.Netmon.Get().GetHash())
	state.Netmon.Info = types.StringValue(instance.Netmon.Get().GetInfo())
	state.Netmon.Enabled = types.BoolValue(instance.Netmon.Get().GetEnabled())

	state.ExposePort = types.StringValue(instance.GetExposePort())
	state.Fault = types.BoolValue(instance.GetFault())

	patches, diags := types.ListValueFrom(ctx, types.StringType, instance.GetPatches())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Patches = patches

	state.CreatedBy = &V1InstanceCreatedByModel{
		Id:       types.StringValue(instance.CreatedBy.GetId()),
		Username: types.StringValue(instance.CreatedBy.GetUsername()),
		Label:    types.StringValue(instance.CreatedBy.GetLabel()),
		Deleted:  types.BoolValue(instance.CreatedBy.GetDeleted()),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumV1InstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state V1InstanceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.InstancesApi.V1DeleteInstance(auth, state.Id.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting instance",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to delete instance",
			"An unexpected error was encountered trying to delete the instance:\n\n"+string(b))
		return
	}

	type deleteStateStructure struct {
		Id string
	}

	const deleteState = "deleted"

	deleteStateConf := &retry.StateChangeConf{
		Refresh: func() (interface{}, string, error) {
			instance, _, _ := d.client.InstancesApi.V1GetInstance(auth, state.Id.ValueString()).Execute()
			if instance != nil {
				return instance, string(instance.GetState()), nil
			}

			return deleteStateStructure{Id: state.Id.ValueString()}, deleteState, nil
		},
		Pending: []string{
			V1InstanceStateDeleting,
		},
		Target: []string{
			deleteState,
		},
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
		Timeout:    5 * time.Minute,
	}

	if _, err = deleteStateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error deleting instance",
			"An unexpected error was encountered trying to delete the instance:\n\n"+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (d *CorelliumV1InstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*corellium.APIClient)
}
