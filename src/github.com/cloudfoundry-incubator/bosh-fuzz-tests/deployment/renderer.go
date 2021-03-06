package deployment

import (
	"bytes"
	"text/template"

	bftinput "github.com/cloudfoundry-incubator/bosh-fuzz-tests/input"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Renderer interface {
	Render(input bftinput.Input, manifestPath string, cloudConfigPath string) error
}

type renderer struct {
	fs boshsys.FileSystem
}

func NewRenderer(fs boshsys.FileSystem) Renderer {
	return &renderer{
		fs: fs,
	}
}

func (g *renderer) Render(input bftinput.Input, manifestPath string, cloudConfigPath string) error {
	deploymentTemplate := template.Must(template.New("deployment").Parse(DeploymentTemplate))

	buffer := bytes.NewBuffer([]byte{})
	err := deploymentTemplate.Execute(buffer, input)
	if err != nil {
		return bosherr.WrapErrorf(err, "Generating deployment manifest")
	}

	err = g.fs.WriteFile(manifestPath, buffer.Bytes())
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving generated manifest")
	}

	cloudTemplate := template.Must(template.New("cloud-config").Parse(CloudTemplate))

	buffer = bytes.NewBuffer([]byte{})

	err = cloudTemplate.Execute(buffer, input)
	if err != nil {
		return bosherr.WrapErrorf(err, "Generating cloud config")
	}

	err = g.fs.WriteFile(cloudConfigPath, buffer.Bytes())
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving generated cloud config")
	}

	return nil
}
