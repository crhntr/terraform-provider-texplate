package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"texplate_interpolate": resourceServer(),
		},
	}
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Read: func(data *schema.ResourceData, _ interface{}) error {
			return interpolate(data, nil, generateID)
		},

		Schema: map[string]*schema.Schema{
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"output": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type ResourceData interface {
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	Set(string, interface{}) error
	SetId(string)
}

func interpolate(data ResourceData, _ interface{}, generateID func(string) string) error {
	templateString := data.Get("template").(string)
	vars, _ := data.GetOk("vars")

	t, _ := template.New("template").Parse(templateString)
	var buffer bytes.Buffer
	t.Execute(&buffer, vars)

	// // if !hasVars {
	// h := sha256.New()
	// io.WriteString(h, templateString)
	// id := fmt.Sprintf("%x", h.Sum(nil))

	// data.SetId(id)
	data.Set("output", string(buffer.Bytes()))
	return nil
	// }

	// return errors.New("not implemented yet")
	// interpolate
}

func generateID(str string) string {
	h := sha256.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/*

data "texplate" "director_config" {
	template = ${data.file.director_config.file_contents} // some yaml we loaded from a file
	vars {


	}
}


*/
