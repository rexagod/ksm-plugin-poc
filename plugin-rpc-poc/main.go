// OUT_OF_TREE_POC_IMPL //

package main

import (
	"os"
	"os/exec"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/rexagod/ksm-rpc-plugin-poc/shared"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "KSM_PLUGIN",
			MagicCookieValue: "deployment_collector",
		},
		Plugins: map[string]plugin.Plugin{
			"deployment_collector": &shared.MetricFamiliesPlugin{},
		},
		Cmd:    exec.Command("./plugin/ksm-deployment-collector"),
		Logger: logger,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		panic(err)
	}

	raw, err := rpcClient.Dispense("deployment_collector")
	if err != nil {
		panic(err)
	}

	deploymentCollector := raw.(shared.MetricFamilies)
	f := deploymentCollector.XMetricFamilies(nil, nil)
	spew.Dump(f)
}
