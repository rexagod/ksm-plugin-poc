package main

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	generator "k8s.io/kube-state-metrics/v2/pkg/metric_generator"

	"github.com/rexagod/ksm-rpc-plugin-poc/shared"
)

type DeploymentInterface struct {
	logger hclog.Logger
}

func (d *DeploymentInterface) XMetricFamilies(allowAnnotationsList, allowLabelsList []string) []generator.FamilyGenerator {
	d.logger.Info("DeploymentInterface.XMetricFamilies called")
	out := deploymentMetricFamilies(allowAnnotationsList, allowLabelsList)
	return out
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "KSM_PLUGIN",
	MagicCookieValue: "deployment_collector",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	deploymentCollector := &DeploymentInterface{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"deployment_collector": &shared.MetricFamiliesPlugin{Impl: deploymentCollector},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins:         pluginMap,
	})
}
