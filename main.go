package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
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
			"texplate_execute": resourceServer(),
		},
	}
}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Read: func(data *schema.ResourceData, _ interface{}) error {
			return execute(data, nil, generateID, defaultTemplate())
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

type Template interface {
	New(name string) *template.Template
	Parse(string) (*template.Template, error)
	Execute(io.Writer, interface{}) error
}

func execute(data ResourceData, _ interface{}, generateID func(string) string, tmpl Template) error {
	templateString := data.Get("template").(string)
	vars, _ := data.GetOk("vars")

	t, _ := tmpl.Parse(templateString)
	// if err != nil {
	// 	return errors.New("invalid template")
	// }
	var buffer bytes.Buffer
	t.Execute(&buffer, vars)

	data.SetId(generateID(""))
	data.Set("output", string(buffer.Bytes()))
	return nil
}

func defaultTemplate() *template.Template {
	tmpl := template.New("template")
	// tmpl = tmpl.Option("missingkey=error")
	// tmpl = tmpl.Funcs(sprig.FuncMap())
	// tmpl = tmpl.Funcs(map[string]interface{}{"cidrhost": cidrhost})
	return tmpl
}

func generateID(str string) string {
	h := sha256.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// funcs

func cidrhost(cidrVal string, hostIndex int) (string, error) {
	// adapted from https://github.com/hashicorp/terraform/blob/fe0cc3b0db0d1a5676c3d1a92ea8c5ff829b4233/config/interpolate_funcs.go#L253-L264
	_, network, err := net.ParseCIDR(cidrVal)
	if err != nil {
		return "", fmt.Errorf("invalid CIDR expression: %s", err)
	}

	ip, err := cidr.Host(network, hostIndex)
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}
