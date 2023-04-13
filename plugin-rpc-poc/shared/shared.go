package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"

	"k8s.io/kube-state-metrics/v2/pkg/metric_generator"
)

type MetricFamilies interface {
	XMetricFamilies([]string, []string) []generator.FamilyGenerator
}

type MetricFamiliesRPC struct {
	client *rpc.Client
}

func (m *MetricFamiliesRPC) XMetricFamilies(allowAnnotationsList, allowLabelsList []string) []generator.FamilyGenerator {
	var resp []generator.FamilyGenerator
	err := m.client.Call("Plugin.XMetricFamilies", [][]string{allowAnnotationsList, allowLabelsList}, &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

type MetricFamiliesRPCServer struct {
	Impl MetricFamilies
}

func (s *MetricFamiliesRPCServer) XMetricFamilies(args [][]string, resp *[]generator.FamilyGenerator) error {
	*resp = s.Impl.XMetricFamilies(args[0], args[1])
	return nil
}

type MetricFamiliesPlugin struct {
	Impl MetricFamilies
}

func (p *MetricFamiliesPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &MetricFamiliesRPCServer{Impl: p.Impl}, nil
}

func (*MetricFamiliesPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &MetricFamiliesRPC{client: c}, nil
}
