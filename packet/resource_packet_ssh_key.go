package packet

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func packetSSHKeyCommonFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"public_key": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"fingerprint": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"updated": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

}

func resourcePacketSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketSSHKeyCreate,
		Read:   resourcePacketSSHKeyRead,
		Update: resourcePacketSSHKeyUpdate,
		Delete: resourcePacketSSHKeyDelete,

		Schema: packetSSHKeyCommonFields(),
	}
}

func resourcePacketSSHKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.SSHKeyCreateRequest{
		Label: d.Get("name").(string),
		Key:   d.Get("public_key").(string),
	}

	projectID, isProjectKey := d.GetOk("project_id")
	if isProjectKey {
		createRequest.ProjectID = projectID.(string)
	}

	key, _, err := client.SSHKeys.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(key.ID)
	if isProjectKey {
		return resourcePacketProjectSSHKeyRead(d, meta)
	}

	return resourcePacketSSHKeyRead(d, meta)
}

func resourcePacketSSHKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	key, _, err := client.SSHKeys.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)

		// If the key is somehow already destroyed, mark as
		// succesfully gone
		if isNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("id", key.ID)
	d.Set("name", key.Label)
	d.Set("public_key", key.Key)
	d.Set("fingerprint", key.FingerPrint)
	d.Set("created", key.Created)
	d.Set("updated", key.Updated)

	return nil
}

func resourcePacketSSHKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	updateRequest := &packngo.SSHKeyUpdateRequest{}

	if d.HasChange("name") {
		kName := d.Get("name").(string)
		updateRequest.Label = &kName
	}

	if d.HasChange("public_key") {
		kKey := d.Get("public_key").(string)
		updateRequest.Key = &kKey
	}

	_, _, err := client.SSHKeys.Update(d.Id(), updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourcePacketSSHKeyRead(d, meta)
}

func resourcePacketSSHKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.SSHKeys.Delete(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
