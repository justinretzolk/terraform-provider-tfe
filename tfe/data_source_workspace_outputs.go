package tfe

import (
	"fmt"
	"log"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTFEWorkspaceOutputs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTFEWorkspaceOutputsRead,

		Schema: map[string]*schema.Schema{
			"outputs": {
				Type:     schema.TypeMap,
				Computed: true,
				//				Elem: &schema.Schema{
				//					Type:     schema.TypeMap,
				//					Computed: true,
				//Elem: &schema.Schema{Type: schema.TypeString},
				//				},
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceTFEWorkspaceOutputsRead(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the workspace ID.
	workspaceID := d.Get("workspace_id").(string)
	d.Set("workspace_id", workspaceID)

	// Get the current state version of the workspace.
	sv, err := tfeClient.StateVersions.Current(ctx, workspaceID)
	if err != nil {
		return fmt.Errorf(
			"Error retrieving current state version of workspace %s: %v", workspaceID, err)
	}

	moutputs := make(map[string]string, 0)

	// Get the value of each output found in the most recent state version
	if sv.Outputs != nil {
		for _, o := range sv.Outputs {
			wo, err := tfeClient.StateVersionOutputs.Read(ctx, o.ID)
			if err != nil {
				return fmt.Errorf(
					"Error retrieving state version output %s: %v", o.ID, err)
			}
			moutputs[wo.Name] = wo.Value
		}
	}
	log.Printf("outputs: %#v", moutputs)

	d.SetId(workspaceID)
	if err := d.Set("outputs", moutputs); err != nil {
		return fmt.Errorf(
			"Error setting workspace outputs: %v", err)
	}

	return nil
}
