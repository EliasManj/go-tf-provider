package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"terraform-provider-movies/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewMovieResource() resource.Resource {
	return &MovieResource{}
}

type MovieResource struct {
	client *client.Client
}

type MovieResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Title types.String `tfsdk:"title"`
}

func (r *MovieResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create implements resource.Resource.
func (r *MovieResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Read Terraform plan data into the model
	var plan *MovieResourceModel
	var raw string
	diags := req.Plan.Get(ctx, &plan)
	req.Plan.Raw.As(&raw)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create APi request
	movie := client.Movie{
		ID:    plan.ID.ValueString(),
		Title: plan.Title.ValueString(),
	}

	movieJson, err := json.Marshal(movie)
	if err != nil {
		resp.Diagnostics.AddError("Error encoding JSON:", err.Error())
		return
	}

	movieRequest, err := http.NewRequest("POST", r.client.HostUrl+"/movies", bytes.NewBuffer(movieJson))
	movieRequest.Header.Set("Content-Type", "application/json")

	body, err := r.client.DoRequest(movieRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating movie",
			"Could not create movie, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	var createdMovie client.Movie
	err = json.Unmarshal(body, &createdMovie)
	if err != nil {
		return
	}
	plan.ID = types.StringValue(createdMovie.ID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	return
}

// Delete implements resource.Resource.
func (*MovieResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Metadata implements resource.Resource.
func (*MovieResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_movie"
}

// Read implements resource.Resource.
func (*MovieResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements resource.Resource.
func (*MovieResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"title": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
			},
		},
	}
}

// Update implements resource.Resource.
func (*MovieResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}
