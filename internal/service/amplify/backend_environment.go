// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package amplify

import (
	"context"
	"log"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amplify"
	awstypes "github.com/aws/aws-sdk-go-v2/service/amplify/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_amplify_backend_environment")
func ResourceBackendEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceBackendEnvironmentCreate,
		ReadWithoutTimeout:   resourceBackendEnvironmentRead,
		DeleteWithoutTimeout: resourceBackendEnvironmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"deployment_artifacts": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexache.MustCompile(`^[0-9A-Za-z-]{1,100}$`), "should be not be more than 100 alphanumeric or hyphen characters"),
			},

			"environment_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexache.MustCompile(`^[a-z]{2,10}$`), "should be between 2 and 10 characters (only lowercase alphabetic)"),
			},

			"stack_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexache.MustCompile(`^[0-9A-Za-z-]{1,100}$`), "should be not be more than 100 alphanumeric or hyphen characters"),
			},
		},
	}
}

func resourceBackendEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AmplifyClient(ctx)

	appID := d.Get("app_id").(string)
	environmentName := d.Get("environment_name").(string)
	id := BackendEnvironmentCreateResourceID(appID, environmentName)

	input := &amplify.CreateBackendEnvironmentInput{
		AppId:           aws.String(appID),
		EnvironmentName: aws.String(environmentName),
	}

	if v, ok := d.GetOk("deployment_artifacts"); ok {
		input.DeploymentArtifacts = aws.String(v.(string))
	}

	if v, ok := d.GetOk("stack_name"); ok {
		input.StackName = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Creating Amplify Backend Environment: %+v", input)
	_, err := conn.CreateBackendEnvironment(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Amplify Backend Environment (%s): %s", id, err)
	}

	d.SetId(id)

	return append(diags, resourceBackendEnvironmentRead(ctx, d, meta)...)
}

func resourceBackendEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AmplifyClient(ctx)

	appID, environmentName, err := BackendEnvironmentParseResourceID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing Amplify Backend Environment ID: %s", err)
	}

	backendEnvironment, err := FindBackendEnvironmentByAppIDAndEnvironmentName(ctx, conn, appID, environmentName)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Amplify Backend Environment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Amplify Backend Environment (%s): %s", d.Id(), err)
	}

	d.Set("app_id", appID)
	d.Set("arn", backendEnvironment.BackendEnvironmentArn)
	d.Set("deployment_artifacts", backendEnvironment.DeploymentArtifacts)
	d.Set("environment_name", backendEnvironment.EnvironmentName)
	d.Set("stack_name", backendEnvironment.StackName)

	return diags
}

func resourceBackendEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AmplifyClient(ctx)

	appID, environmentName, err := BackendEnvironmentParseResourceID(d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "parsing Amplify Backend Environment ID: %s", err)
	}

	log.Printf("[DEBUG] Deleting Amplify Backend Environment: %s", d.Id())
	_, err = conn.DeleteBackendEnvironment(ctx, &amplify.DeleteBackendEnvironmentInput{
		AppId:           aws.String(appID),
		EnvironmentName: aws.String(environmentName),
	})

	if errs.IsA[*awstypes.NotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Amplify Backend Environment (%s): %s", d.Id(), err)
	}

	return diags
}
