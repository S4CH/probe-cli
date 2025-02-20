package registry

//
// Registers the `dnsping' experiment.
//

import (
	"github.com/ooni/probe-cli/v3/internal/experiment/dnsping"
	"github.com/ooni/probe-cli/v3/internal/model"
)

func init() {
	AllExperiments["dnsping"] = &Factory{
		build: func(config interface{}) model.ExperimentMeasurer {
			return dnsping.NewExperimentMeasurer(
				*config.(*dnsping.Config),
			)
		},
		config:      &dnsping.Config{},
		inputPolicy: model.InputOrStaticDefault,
	}
}
