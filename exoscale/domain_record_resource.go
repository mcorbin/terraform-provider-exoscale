package exoscale

import (
	"fmt"
	"strconv"

	"github.com/exoscale/egoscale"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func domainRecordResource() *schema.Resource {
	return &schema.Resource{
		Create: createRecord,
		Read:   readRecord,
		Exists: existsRecord,
		Update: updateRecord,
		Delete: deleteRecord,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"record_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A", "AAAA", "ALIAS", "CNAME", "HINFO", "MX", "NAPTR",
					"NS", "POOL", "SPF", "SRV", "SSHFP", "TXT", "URL",
				}, true),
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"prio": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createRecord(d *schema.ResourceData, meta interface{}) error {
	client := GetDNSClient(meta)

	record, err := client.CreateRecord(d.Get("domain").(string), egoscale.DNSRecord{
		Name:       d.Get("name").(string),
		Content:    d.Get("content").(string),
		RecordType: d.Get("record_type").(string),
		TTL:        d.Get("ttl").(int),
		Prio:       d.Get("prio").(int),
	})

	if err != nil {
		return err
	}

	return applyRecord(d, *record)
}

func existsRecord(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := GetDNSClient(meta)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	_, err := client.GetRecord(d.Get("domain").(string), id)

	return err == nil, err
}

func readRecord(d *schema.ResourceData, meta interface{}) error {
	client := GetDNSClient(meta)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	record, err := client.GetRecord(d.Get("domain").(string), id)
	if err != nil {
		return err
	}

	return applyRecord(d, *record)
}

func updateRecord(d *schema.ResourceData, meta interface{}) error {
	client := GetDNSClient(meta)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	record, err := client.UpdateRecord(d.Get("domain").(string), egoscale.DNSRecord{
		ID:         id,
		Name:       d.Get("name").(string),
		Content:    d.Get("content").(string),
		RecordType: d.Get("record_type").(string),
		TTL:        d.Get("ttl").(int),
		Prio:       d.Get("prio").(int),
	})

	if err != nil {
		return err
	}

	return applyRecord(d, *record)
}

func deleteRecord(d *schema.ResourceData, meta interface{}) error {
	client := GetDNSClient(meta)

	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	err := client.DeleteRecord(d.Get("domain").(string), id)
	if err != nil {
		d.SetId("")
	}

	return err
}

func applyRecord(d *schema.ResourceData, record egoscale.DNSRecord) error {
	d.SetId(strconv.FormatInt(record.ID, 10))
	d.Set("name", record.Name)
	d.Set("content", record.Content)
	d.Set("record_type", record.RecordType)
	d.Set("ttl", record.TTL)
	d.Set("prio", record.Prio)

	domain := d.Get("domain").(string)
	if record.Name == "" {
		d.Set("hostname", domain)
	} else {
		d.Set("hostname", fmt.Sprintf("%s.%s", record.Name, domain))
	}

	return nil
}
