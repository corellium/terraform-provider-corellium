package corellium

import (
	"context"
	"io"
	"net/http"
	"time"

	"terraform-provider-corellium/corellium/pkg/api"

	corellium "github.com/corellium/corellium-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumInstanceResource{}
	_ resource.ResourceWithConfigure = &CorelliumInstanceResource{}
)

// NewCorelliumInstanceResource is a helper function to simplify the provider implementation.
func NewCorelliumInstanceResource() resource.Resource {
	return &CorelliumInstanceResource{}
}

// CorelliumInstanceResource is the resource implementation.
type CorelliumInstanceResource struct {
	client *corellium.CorelliumClient
}

type InstanceVPNModel struct {
	Proxy     types.Map `tfsdk:"proxy"`
	Listeners types.Map `tfsdk:"listeners"`
}

type InstanceServicesModel struct {
	VPN *InstanceVPNModel `tfsdk:"vpn"`
}

type InstanceAgentModel struct {
	Hash types.String `tfsdk:"hash"`
	Info types.String `tfsdk:"info"`
}

type InstanceNetmonModel struct {
	Hash    types.String `tfsdk:"hash"`
	Info    types.String `tfsdk:"info"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

type InstanceCreatedByModel struct { // TODO: rename to user or use the user model.
	// Id is the user who created the instance.
	Id types.String `tfsdk:"id"`
	// Username is the username of the user who created the instance.
	Username types.String `tfsdk:"username"`
	// Label is the label of the user who created the instance.
	Label types.String `tfsdk:"label"`
	// Deleted is the deleted status of the user who created the instance.
	Deleted types.Bool `tfsdk:"deleted"`
}

type InstanceBootOptionsModel struct {
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
	// InstanceStateOn is the state of the instance when it is running.
	InstanceStateOn = "on"
	// InstanceStateOff is the state of the instance when it is powered off.
	InstanceStateOff = "off"
	// InstanceStatePaused is the state of the instance when it is paused.
	InstanceStatePaused = "paused"
	// InstanceStateCreating is the state of the instance when it is being created.
	InstanceStateCreating = "creating"
	// InstanceStateDeleting is the state of the instance when it is being deleted.
	InstanceStateDeleting = "deleting"
)

const (
	// InstancesTaskStateNone is the state of the instance when there is no task.
	InstancesTaskStateNone = "none"
	// InstancesTaskStateBuilding is the state of the instance when it is being built.
	InstancesTaskStateBuilding = "building"
)

// coffeesModel maps coffees schema data.
type InstanceModel struct {
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
	BootOptions *InstanceBootOptionsModel `tfsdk:"boot_options"`
	ServiceIP   types.String              `tfsdk:"service_ip"`
	WifiIP      types.String              `tfsdk:"wifi_ip"`
	SecondaryIP types.String              `tfsdk:"secondary_ip"`
	// Services    *InstanceServicesModel    `tfsdk:"services"` // TODO: find the right type for this.
	Panicked types.Bool `tfsdk:"panicked"`
	// Created is the time the instance was created.
	Created   types.String `tfsdk:"created"`
	Model     types.String `tfsdk:"model"`
	FWPackage types.String `tfsdk:"fwpackage"`
	// OS is the version of the operating system running on the instance, e.g. 14.3 for iOS, or 11.0.0 for Android.
	OSVersion  types.String         `tfsdk:"os_version"`
	OSBuild    types.String         `tfsdk:"os_build"`
	Agent      *InstanceAgentModel  `tfsdk:"agent"`
	Netmon     *InstanceNetmonModel `tfsdk:"netmon"`
	ExposePort types.String         `tfsdk:"expose_port"`
	Fault      types.Bool           `tfsdk:"fault"`
	// Patches is the list of patches applied to the instance.
	//    - jailbroken The instance should be jailbroken (default).
	//    - nonjailbroken The instance should not be jailbroken.
	//    - corelliumd The instance should not be jailbroken but should profile API agent.
	Patches types.List `tfsdk:"patches"`
	// CreatedBy is the user who created the instance.
	CreatedBy *InstanceCreatedByModel `tfsdk:"created_by"`
	// WaitForReady is a boolean that indicates if the resource should wait for the instance to be ready.
	WaitForReady types.Bool `tfsdk:"wait_for_ready"`
	// WaitForReadyTimeout is the timeout in seconds to wait for the instance to be ready.
	// WaitForReadyTimeout is a amount in seconds.
	WaitForReadyTimeout types.Int64 `tfsdk:"wait_for_ready_timeout"`
}

// Metadata returns the resource type name.
func (d *CorelliumInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
	// TypeName is the name of the resource type, which must be unique within the provider.
	// This is used to identify the resource type in state and plan files.
	// i.e: resource corellium_instance "instance" { ... }
}

// Schema defines the schema for the resource.
func (d *CorelliumInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Instance type",
				Computed:    true,
			},
			"project": schema.StringAttribute{
				Description: "Instance project",
				Optional:    true,
				Computed:    true,
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
			"os_version": schema.StringAttribute{
				Description: "Instance os version",
				Required:    true,
			},
			"os_build": schema.StringAttribute{
				Description: "Instance os build",
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

// WaitForReadyTimeout is the default timeout for waiting for an instance to be ready.
const WaitForReadyTimeout = 900

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan InstanceModel

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())

	if plan.Project.IsNull() || plan.Project.IsUnknown() {
		projects, r, err := d.client.Projects.List(auth).Execute()
		if err != nil {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error get default project",
					"Coudn't read the response body from project name check: "+err.Error(),
				)
				return
			}

			resp.Diagnostics.AddError(
				"Error get default project",
				"Couldn't get projects to create instance: "+string(b),
			)
			return
		}

		// Assuming that the default project is the first one in the list.
		plan.Project = types.StringValue(projects.Results[0].Id)

		diags = resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createInstanceBody := corellium.NewCreateInstanceBody(plan.Project.ValueString())
	createInstanceBody.SetFlavor(plan.Flavor.ValueString())
	createInstanceBody.SetOsVersion(plan.OSVersion.ValueString())
	createInstanceBody.SetName(plan.Name.ValueString())
	createInstanceBody.SetOsBuild(plan.OSBuild.ValueString())
	created, r, err := d.client.Instances.Create(auth).CreateInstanceBody(*createInstanceBody).Execute()
	if err != nil {
		if r.StatusCode == http.StatusForbidden {
			resp.Diagnostics.AddError(
				"Error creating instance",
				"You don't have permission to create instances in this project.",
			)
			return
		}

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
				instance, r, err := d.client.Instances.Get(auth, created.GetId()).Execute()
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
						"Error creating instance",
						"An unexpected error was encountered trying to check if instance is ready:\n\n"+string(b),
					)
					return nil, "", err
				}

				return instance, string(instance.GetState()), nil
			},
			Pending: []string{
				InstanceStateCreating,
				InstanceStateDeleting,
			},
			Target: []string{
				InstanceStateOn,
				InstanceStateOff,
				InstanceStatePaused,
			},
			Delay:      5 * time.Second,
			MinTimeout: 5 * time.Second,
			Timeout: func() time.Duration {
				if plan.WaitForReadyTimeout.IsUnknown() {
					return WaitForReadyTimeout
				}

				return time.Duration(plan.WaitForReadyTimeout.ValueInt64())
			}() * time.Second,
		}

		if _, err := createStateConf.WaitForStateContext(ctx); err != nil {
			resp.Diagnostics.AddError(
				"Error creating instance",
				"Coudn't create the instance: "+err.Error(),
			)

			return
		}
	}

	instance, r, err := d.client.Instances.Get(auth, created.GetId()).Execute()
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
	plan.StartedAt = types.StringValue(instance.GetStartedAt().Format(time.RFC3339))
	plan.UserTask = types.StringValue(instance.GetUserTask())
	plan.TaskState = types.StringValue(instance.GetTaskState())
	plan.Error = types.StringValue(instance.GetError())

	additionalTags, diags := types.ListValueFrom(ctx, types.StringType, instance.BootOptions.GetAdditionalTags())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.BootOptions = &InstanceBootOptionsModel{
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

	plan.Services = &InstanceServicesModel{
		VPN: &InstanceVPNModel{
			Proxy:     proxy,
			Listeners: listeners,
		},
	}*/

	plan.Panicked = types.BoolValue(instance.GetPanicked())
	plan.Created = types.StringValue(instance.GetCreated().UTC().String())
	plan.Model = types.StringValue(instance.GetModel())
	plan.FWPackage = types.StringValue(instance.GetFwPackage())
	plan.OSBuild = types.StringValue(instance.GetOsBuild())
	plan.OSVersion = types.StringValue(instance.GetOsVersion())
	plan.Agent = &InstanceAgentModel{
		Hash: types.StringValue(instance.Agent.Get().GetHash()),
		Info: types.StringValue(instance.Agent.Get().GetInfo()),
	}
	netmon := instance.GetNetMon()
	plan.Netmon = &InstanceNetmonModel{
		Hash:    types.StringValue(netmon.GetHash()),
		Info:    types.StringValue(netmon.GetInfo()),
		Enabled: types.BoolValue(netmon.GetEnabled()),
	}

	plan.ExposePort = types.StringValue(instance.GetExposePort())
	plan.Fault = types.BoolValue(instance.GetFault())

	patches, diags := types.ListValueFrom(ctx, types.StringType, instance.GetPatches())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Patches = patches

	plan.CreatedBy = &InstanceCreatedByModel{
		Id:       types.StringValue(instance.CreatedBy.Get().GetId()),
		Username: types.StringValue(instance.CreatedBy.Get().GetUsername()),
		Label:    types.StringValue(instance.CreatedBy.Get().GetLabel()),
		Deleted:  types.BoolValue(instance.CreatedBy.Get().GetDeleted()),
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state InstanceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	instance, r, err := d.client.Instances.Get(auth, state.Id.ValueString()).Execute()
	if err != nil {
		if r.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}

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
	state.StartedAt = types.StringValue(instance.GetStartedAt().Format(time.RFC3339))
	state.UserTask = types.StringValue(instance.GetUserTask())
	state.TaskState = types.StringValue(instance.GetTaskState())
	state.Error = types.StringValue(instance.GetError())

	additionalTags, diags := types.ListValueFrom(ctx, types.StringType, instance.BootOptions.GetAdditionalTags())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.BootOptions = &InstanceBootOptionsModel{
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

	state.Services = &InstanceServicesModel{
		VPN: &InstanceVPNModel{
			Proxy:     proxy,
			Listeners: listeners,
		},
	}*/

	state.Panicked = types.BoolValue(instance.GetPanicked())
	state.Created = types.StringValue(instance.GetCreated().UTC().String())
	state.Model = types.StringValue(instance.GetModel())
	state.FWPackage = types.StringValue(instance.GetFwPackage())
	state.OSBuild = types.StringValue(instance.GetOsBuild())
	state.OSVersion = types.StringValue(instance.GetOsVersion())
	state.Agent = &InstanceAgentModel{
		Hash: types.StringValue(instance.Agent.Get().GetHash()),
		Info: types.StringValue(instance.Agent.Get().GetInfo()),
	}
	state.Netmon = &InstanceNetmonModel{
		Hash:    types.StringValue(instance.NetMon.GetHash()),
		Info:    types.StringValue(instance.NetMon.GetInfo()),
		Enabled: types.BoolValue(instance.NetMon.GetEnabled()),
	}
	state.ExposePort = types.StringValue(instance.GetExposePort())
	state.Fault = types.BoolValue(instance.GetFault())

	patches, diags := types.ListValueFrom(ctx, types.StringType, instance.GetPatches())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Patches = patches

	state.CreatedBy = &InstanceCreatedByModel{
		Id:       types.StringValue(instance.CreatedBy.Get().GetId()),
		Username: types.StringValue(instance.CreatedBy.Get().GetUsername()),
		Label:    types.StringValue(instance.CreatedBy.Get().GetLabel()),
		Deleted:  types.BoolValue(instance.CreatedBy.Get().GetDeleted()),
	}

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state InstanceModel
	// state is the current state of the resource.

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan InstanceModel
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

	if !state.OSBuild.Equal(plan.OSBuild) {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"OS build cannot be changed",
		)
		return
	}

	if !state.OSVersion.Equal(plan.OSVersion) {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"OS version cannot be changed",
		)
		return
	}

	if state.State.ValueString() == InstanceStateCreating {
		resp.Diagnostics.AddWarning(
			"Warning updating instance",
			"Instance is currently being created. Changes to its state cannot be applied.",
		)
	}

	if state.State.ValueString() == InstanceStateDeleting {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"Instance is currently being deleted. Changes cannot be applied.",
		)
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	patchInstanceBody := corellium.NewPatchInstanceBody()
	patchInstanceBody.SetName(plan.Name.ValueString())
	instance, r, err := d.client.Instances.Patch(auth, state.Id.ValueString()).PatchInstanceBody(*patchInstanceBody).Execute()
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
	state.StartedAt = types.StringValue(instance.GetStartedAt().Format(time.RFC3339))
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

	state.Services = &InstanceServicesModel{
		VPN: &InstanceVPNModel{
			Proxy:     proxy,
			Listeners: listeners,
		},
	}*/

	state.Panicked = types.BoolValue(instance.GetPanicked())
	state.Created = types.StringValue(instance.GetCreated().UTC().String())
	state.Model = types.StringValue(instance.GetModel())
	state.FWPackage = types.StringValue(instance.GetFwPackage())
	state.OSBuild = types.StringValue(instance.GetOsBuild())
	state.OSVersion = types.StringValue(instance.GetOsVersion())

	state.Agent = &InstanceAgentModel{
		Hash: types.StringValue(instance.Agent.Get().GetHash()),
		Info: types.StringValue(instance.Agent.Get().GetInfo()),
	}

	netmon := instance.GetNetMon()
	state.Netmon = &InstanceNetmonModel{
		Hash:    types.StringValue(netmon.GetHash()),
		Info:    types.StringValue(netmon.GetInfo()),
		Enabled: types.BoolValue(netmon.GetEnabled()),
	}

	state.ExposePort = types.StringValue(instance.GetExposePort())
	state.Fault = types.BoolValue(instance.GetFault())

	patches, diags := types.ListValueFrom(ctx, types.StringType, instance.GetPatches())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Patches = patches

	if instance.CreatedBy.IsSet() {
		createdBy := instance.CreatedBy.Get()

		state.CreatedBy = &InstanceCreatedByModel{
			Id:       types.StringValue(createdBy.GetId()),
			Username: types.StringValue(createdBy.GetUsername()),
			Label:    types.StringValue(createdBy.GetLabel()),
			Deleted:  types.BoolValue(createdBy.GetDeleted()),
		}
	} else {
		state.CreatedBy = nil
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state InstanceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.Instances.Delete(auth, state.Id.ValueString()).Execute()
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
			instance, _, _ := d.client.Instances.Get(auth, state.Id.ValueString()).Execute()
			if instance != nil {
				return instance, string(instance.GetState()), nil
			}

			return deleteStateStructure{Id: state.Id.ValueString()}, deleteState, nil
		},
		Pending: []string{
			InstanceStateDeleting,
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
func (d *CorelliumInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*ProviderData).Client
}
