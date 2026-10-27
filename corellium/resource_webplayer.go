package corellium

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"terraform-provider-corellium/corellium/pkg/api"

	corellium "github.com/corellium/corellium-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &CorelliumWebPlayerResource{}
	_ resource.ResourceWithConfigure = &CorelliumWebPlayerResource{}
)

// helper function to simplify the provider implementation.
func NewCorelliumWebPlayerResource() resource.Resource {
	return &CorelliumWebPlayerResource{}
}

// the resource implementation.
type CorelliumWebPlayerResource struct {
	client *corellium.CorelliumClient
	host   string
}

type WebPlayerDataModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instanceid"`
	Identifier types.String `tfsdk:"identifier"`
	Project    types.String `tfsdk:"project"`
	Token      types.String `tfsdk:"token"`
	// Expiration       types.String  `tfsdk:"expiration"`
	Expiresinseconds types.Float64 `tfsdk:"expiresinseconds"`
	ClientId         types.String  `tfsdk:"clientid"`
	// LastActivity types.String  `tfsdk:"lastactivity"`
	// CreatedAt    types.String  `tfsdk:"createdat"`
	// UpdatedAt    types.String  `tfsdk:"updatedat"`
	Features Features `tfsdk:"features"`
}

type Features struct {
	Apps           types.Bool `tfsdk:"apps"`
	Console        types.Bool `tfsdk:"console"`
	Coretrace      types.Bool `tfsdk:"coretrace"`
	DeviceControl  types.Bool `tfsdk:"devicecontrol"`
	DeviceDelete   types.Bool `tfsdk:"devicedelete"`
	Files          types.Bool `tfsdk:"files"`
	Frida          types.Bool `tfsdk:"frida"`
	Images         types.Bool `tfsdk:"images"`
	Messaging      types.Bool `tfsdk:"messaging"`
	Netmon         types.Bool `tfsdk:"netmon"`
	Network        types.Bool `tfsdk:"network"`
	PortForwarding types.Bool `tfsdk:"portforwarding"`
	Profile        types.Bool `tfsdk:"profile"`
	Sensors        types.Bool `tfsdk:"sensors"`
	Settings       types.Bool `tfsdk:"settings"`
	Snapshots      types.Bool `tfsdk:"snapshots"`
	Strace         types.Bool `tfsdk:"strace"`
	System         types.Bool `tfsdk:"system"`
	Connect        types.Bool `tfsdk:"connect"`
}

// Metadata returns the resource type name.
func (d *CorelliumWebPlayerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webplayer"
}

// Schema defines the schema for the resource.
func (d *CorelliumWebPlayerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The sessionID of the created web player",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"instanceid": schema.StringAttribute{
				Required: true,
			},
			"identifier": schema.StringAttribute{
				Required: false,
				Optional: false,
				Computed: true,
			},
			"project": schema.StringAttribute{
				Required: true,
			},
			"expiresinseconds": schema.Float64Attribute{
				Required: true,
			},
			// "expiration": schema.StringAttribute{
			// 	Computed: true,
			// },
			"clientid": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"token": schema.StringAttribute{
				// Sensitive: true,
				Required: false,
				Optional: false,
				Computed: true,
			},
			// "lastactivity": schema.StringAttribute{
			// 	Computed: true,
			// },
			// "createdate": schema.StringAttribute{
			// 	Computed: true,
			// },
			// "updatedate": schema.StringAttribute{
			// 	Computed: true,
			// },
			"features": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"apps": schema.BoolAttribute{
						Optional: true,
					},
					"console": schema.BoolAttribute{
						Optional: true,
					},
					"coretrace": schema.BoolAttribute{
						Optional: true,
					},
					"devicecontrol": schema.BoolAttribute{
						Optional: true,
					},
					"devicedelete": schema.BoolAttribute{
						Optional: true,
					},
					"files": schema.BoolAttribute{
						Optional: true,
					},
					"frida": schema.BoolAttribute{
						Optional: true,
					},
					"images": schema.BoolAttribute{
						Optional: true,
					},
					"messaging": schema.BoolAttribute{
						Optional: true,
					},
					"netmon": schema.BoolAttribute{
						Optional: true,
					},
					"network": schema.BoolAttribute{
						Optional: true,
					},
					"portforwarding": schema.BoolAttribute{
						Optional: true,
					},
					"profile": schema.BoolAttribute{
						Optional: true,
					},
					"sensors": schema.BoolAttribute{
						Optional: true,
					},
					"settings": schema.BoolAttribute{
						Optional: true,
					},
					"snapshots": schema.BoolAttribute{
						Optional: true,
					},
					"strace": schema.BoolAttribute{
						Optional: true,
					},
					"system": schema.BoolAttribute{
						Optional: true,
					},
					"connect": schema.BoolAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (d *CorelliumWebPlayerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state WebPlayerDataModel

	diags := req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	webPlayerFeatures := corellium.NewWebPlayerFeatures()
	webPlayerFeatures.SetApps(state.Features.Apps.ValueBool())
	webPlayerFeatures.SetConsole((state.Features.Console.ValueBool()))
	webPlayerFeatures.SetCoretrace(state.Features.Coretrace.ValueBool())
	webPlayerFeatures.SetDeviceControl(state.Features.DeviceControl.ValueBool())
	webPlayerFeatures.SetDeviceDelete(state.Features.DeviceDelete.ValueBool())
	webPlayerFeatures.SetFiles(state.Features.Files.ValueBool())
	webPlayerFeatures.SetFrida(state.Features.Frida.ValueBool())
	webPlayerFeatures.SetImages(state.Features.Images.ValueBool())
	webPlayerFeatures.SetMessaging(state.Features.Messaging.ValueBool())
	webPlayerFeatures.SetNetmon(state.Features.Netmon.ValueBool())
	webPlayerFeatures.SetNetwork(state.Features.Network.ValueBool())
	webPlayerFeatures.SetPortForwarding(state.Features.PortForwarding.ValueBool())
	webPlayerFeatures.SetProfile(state.Features.Profile.ValueBool())
	webPlayerFeatures.SetSensors(state.Features.Sensors.ValueBool())
	webPlayerFeatures.SetSettings(state.Features.Settings.ValueBool())
	webPlayerFeatures.SetSnapshots(state.Features.Snapshots.ValueBool())
	webPlayerFeatures.SetStrace(state.Features.Strace.ValueBool())
	webPlayerFeatures.SetSystem(state.Features.System.ValueBool())
	webPlayerFeatures.SetConnect(state.Features.Connect.ValueBool())

	// Check to see if Instance exists. Corellium does a check for project and not for instance.
	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	instance, r, err := d.client.Instances.Get(auth, state.InstanceId.ValueString()).Execute()
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

	if instance == nil || len(instance.GetId()) == 0 {
		resp.Diagnostics.AddError(
			"Error! Instance with instanceId= "+state.InstanceId.ValueString()+" not found",
			"The instance with id "+state.InstanceId.ValueString()+" does not exist",
		)
		return
	}

	projectId := state.Project.ValueString()
	instanceId := state.InstanceId.ValueString()
	expiresIn := float32(state.Expiresinseconds.ValueFloat64())

	createWebPlayerSessionBody := corellium.NewCreateWebPlayerSessionBody(projectId, instanceId, expiresIn, *webPlayerFeatures)
	session, r, err := d.client.WebPlayer.Create(auth).CreateWebPlayerSessionBody(*createWebPlayerSessionBody).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			if r.StatusCode == http.StatusForbidden {
				resp.Diagnostics.AddError(
					"Error creating a web player session",
					"You don't have permission to create a web player session",
				)
				return
			}
			resp.Diagnostics.AddError(
				"Error creating a web player session",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error creating a web player session",
			"An unexpected error was encountered trying to create the web player session:\n\n"+string(b),
		)
		return
	}

	// fmt.Println("*********************************")
	// SESSION.EXPIRATION IS 1683227776611   (sent 3600 got 1683227776611)
	// expirationTimestamp := float64(session.GetExpiration()) // 1683227776611
	// seconds := int64(expirationTimestamp / 1000)
	// nanoseconds := int64((expirationTimestamp - float64(seconds)*1000) * 1e6)
	// expirationTime := time.Unix(seconds, nanoseconds)
	// expirationTime := time.Unix(seconds, 0)

	// Note: depending on how
	// expiration := expirationTime.UTC().Round(time.Second).Format("2006-01-02T15:04:05.000Z07:00")
	// expiration := expirationTime.UTC().Format("2006-01-02T15:04:05.000Z07:00")

	// SESSION.EXPIRATION IS 1683227776611   (sent 3600 got 1683227776611)
	// expirationTimestamp := float64(session.GetExpiration()) // 1683227776611
	// seconds := int64(expirationTimestamp / 1000)
	// nanoseconds := int64((expirationTimestamp - float64(seconds)*1000) * 1e6)
	// expirationTime := time.Unix(seconds, 0)
	// expiration := expirationTime.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	// fmt.Println(expirationTimestamp)
	// fmt.Println(seconds)
	// fmt.Println(nanoseconds)
	// fmt.Println(expirationTime)
	// fmt.Println(expiration)
	// fmt.Println(r)
	state.Identifier = types.StringValue(session.GetIdentifier())
	state.ID = state.Identifier
	state.Token = types.StringValue(session.GetToken())
	state.ClientId = types.StringValue(state.ClientId.ValueString())

	//returns a float value for create endpoint but a string for in ISO-8601 format e.g. 2022-05-06T02:39:23.000Z everywhere else....
	// expirationTimestamp := float64(session.GetExpiration())
	// expirationTime := time.Unix(int64(expirationTimestamp), 0)
	// expiration := expirationTime.UTC().Format(time.RFC3339)

	// state.Expiration = types.StringValue(expiration)

	// expiration := time.Unix().UTC().Format("2022-05-06T02:39:23.000Z")
	// expiration := strconv.FormatFloat(float64(session.GetExpiration()), 'f', -1, 64)
	// state.Expiration = types.StringValue(expiration)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *CorelliumWebPlayerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebPlayerDataModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())

	var sessions []WebPlayerDataModelManual
	var err error
	sessions, err = GetWebPlayerManual(auth, "https://"+d.host+"/api", state.Identifier.ValueString())

	// session, r, err := d.client.WebPlayerApi.WebPlayerSessionInfo(auth, state.Identifier.ValueString()).Execute()
	if err != nil {
		// b, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	resp.Diagnostics.AddError(
		// 		"Error retrieving the session",
		// 		"Coudn't read the response body: "+err.Error(),
		// 	)
		// 	return
		// }

		resp.Diagnostics.AddError(
			"Unable to fetch the web player session",
			// "An unexpected error was encountered trying to read the session:\n\n"+err.Error()+"\n\n"+string(b))
			"An unexpected error was encountered trying to read the session:\n\n"+err.Error())

		return
	}
	if len(sessions) == 0 {
		resp.Diagnostics.AddError(
			"Fetching the web player session returned an empty response",
			"The API returned an empty array. The web player session may have expired or doesn't exist.")
		return
	}

	// Should only have one element
	// state.InstanceId = sessions.InstanceId
	state.Token = types.StringValue(sessions[0].Token)
	// state.Expiration = types.StringValue(sessions[0].Expiration)
	state.ClientId = types.StringValue(sessions[0].ClientId)
	state.Identifier = types.StringValue(sessions[0].Identifier)
	state.Project = types.StringValue(sessions[0].Project)
	state.Features.Apps = types.BoolValue(sessions[0].Features.Apps)
	state.Features.Console = types.BoolValue(sessions[0].Features.Console)
	state.Features.Coretrace = types.BoolValue(sessions[0].Features.Coretrace)
	state.Features.DeviceControl = types.BoolValue(sessions[0].Features.DeviceControl)
	state.Features.DeviceDelete = types.BoolValue(sessions[0].Features.DeviceDelete)
	state.Features.Files = types.BoolValue(sessions[0].Features.Files)
	state.Features.Frida = types.BoolValue(sessions[0].Features.Frida)
	state.Features.Images = types.BoolValue(sessions[0].Features.Images)
	state.Features.Messaging = types.BoolValue(sessions[0].Features.Messaging)
	state.Features.Netmon = types.BoolValue(sessions[0].Features.Netmon)
	state.Features.Network = types.BoolValue(sessions[0].Features.Network)
	state.Features.PortForwarding = types.BoolValue(sessions[0].Features.PortForwarding)
	state.Features.Profile = types.BoolValue(sessions[0].Features.Profile)
	state.Features.Sensors = types.BoolValue(sessions[0].Features.Sensors)
	state.Features.Settings = types.BoolValue(sessions[0].Features.Settings)
	state.Features.Snapshots = types.BoolValue(sessions[0].Features.Snapshots)
	state.Features.Strace = types.BoolValue(sessions[0].Features.Strace)
	state.Features.System = types.BoolValue(sessions[0].Features.System)
	state.Features.Connect = types.BoolValue(sessions[0].Features.Connect)

	// state.InstanceId = types.StringValue(session.GetIdentifier())
	// state.Token = types.StringValue(session.GetToken())
	// state.Expiration = types.Float64Value(float64(session.GetExpiration()))

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (d *CorelliumWebPlayerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	//No update endpoint for the WebPlayerAPI
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *CorelliumWebPlayerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebPlayerDataModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	auth := context.WithValue(ctx, corellium.ContextAccessToken, api.GetAccessToken())
	r, err := d.client.WebPlayer.Delete(auth, state.Identifier.ValueString()).Execute()
	if err != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving the session",
				"Coudn't read the response body: "+err.Error(),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Unable to fetch the web player session",
			"An unexpected error was encountered trying to read the session:\n\n"+err.Error()+"\n\n"+string(b))

		return
	}

}

// Configure adds the provider configured client to the resource.
func (d *CorelliumWebPlayerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData := req.ProviderData.(*ProviderData)
	d.client = providerData.Client
	d.host = providerData.Host
}

// *******************************************************************************************************************************
// Custom struct and methods for the WebPlayerAPI that serve as a workaround for the current state of the API.
type WebPlayerDataModelManual struct {
	InstanceId       string  `tfsdk:"instanceid"`
	Identifier       string  `tfsdk:"identifier"`
	Project          string  `tfsdk:"project"`
	Token            string  `tfsdk:"token"`
	Expiration       string  `tfsdk:"expiration"`
	Expiresinseconds float64 `tfsdk:"expiresinseconds"`
	ClientId         string  `tfsdk:"clientid"`
	// LastActivity string  `tfsdk:"lastactivity"`
	// CreatedAt    string  `tfsdk:"createdat"`
	// UpdatedAt    string  `tfsdk:"updatedat"`
	Features FeaturesManual `tfsdk:"features"`
}

type FeaturesManual struct {
	Apps           bool `tfsdk:"apps"`
	Console        bool `tfsdk:"console"`
	Coretrace      bool `tfsdk:"coretrace"`
	DeviceControl  bool `tfsdk:"devicecontrol"`
	DeviceDelete   bool `tfsdk:"devicedelete"`
	Files          bool `tfsdk:"files"`
	Frida          bool `tfsdk:"frida"`
	Images         bool `tfsdk:"images"`
	Messaging      bool `tfsdk:"messaging"`
	Netmon         bool `tfsdk:"netmon"`
	Network        bool `tfsdk:"network"`
	PortForwarding bool `tfsdk:"portforwarding"`
	Profile        bool `tfsdk:"profile"`
	Sensors        bool `tfsdk:"sensors"`
	Settings       bool `tfsdk:"settings"`
	Snapshots      bool `tfsdk:"snapshots"`
	Strace         bool `tfsdk:"strace"`
	System         bool `tfsdk:"system"`
	Connect        bool `tfsdk:"connect"`
}

func GetWebPlayerManual(ctx context.Context, url string, sessionId string) ([]WebPlayerDataModelManual, error) {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}, // matches the provider's main API client.
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url+"/v1/webplayer/"+sessionId, nil)
	if err != nil {
		return nil, err
	}

	// Get access token from context and add it to the request header
	accessToken, ok := ctx.Value(corellium.ContextAccessToken).(string)
	if !ok {
		return nil, errors.New("access token not found in context")
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching firmware data: %s", resp.Status)
	}

	var sessions []WebPlayerDataModelManual
	err = json.NewDecoder(resp.Body).Decode(&sessions)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// *******************************************************************************************************************************
