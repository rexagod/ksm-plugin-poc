/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	basemetrics "k8s.io/component-base/metrics"

	"k8s.io/kube-state-metrics/v2/pkg/metric"
	generator "k8s.io/kube-state-metrics/v2/pkg/metric_generator"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var descDeploymentLabelsDefaultLabels = []string{"namespace", "deployment"}

func deploymentMetricFamilies(allowAnnotationsList, allowLabelsList []string) []generator.FamilyGenerator {
	fmt.Println("Hello from plugin/deployment.go")
	return []generator.FamilyGenerator{
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_created",
			"Unix creation timestamp",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				var ms []*metric.Metric

				if !d.CreationTimestamp.IsZero() {
					ms = append(ms, &metric.Metric{
						Value: float64(d.CreationTimestamp.Unix()),
					})
				}

				return &metric.Family{
					Metrics: ms,
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_status_replicas",
			"The number of replicas per deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.Replicas),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_status_replicas_ready",
			"The number of ready replicas per deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.ReadyReplicas),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_status_replicas_available",
			"The number of available replicas per deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.AvailableReplicas),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_status_replicas_unavailable",
			"The number of unavailable replicas per deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.UnavailableReplicas),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_status_replicas_updated",
			"The number of updated replicas per deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.UpdatedReplicas),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_status_observed_generation",
			"The generation observed by the deployment controller.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.ObservedGeneration),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_spec_replicas",
			"Number of desired pods for a deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(*d.Spec.Replicas),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_spec_strategy_rollingupdate_max_unavailable",
			"Maximum number of unavailable replicas during a rolling update of a deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				if d.Spec.Strategy.RollingUpdate == nil {
					return &metric.Family{}
				}

				maxUnavailable, err := intstr.GetScaledValueFromIntOrPercent(d.Spec.Strategy.RollingUpdate.MaxUnavailable, int(*d.Spec.Replicas), false)
				if err != nil {
					panic(err)
				}

				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(maxUnavailable),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_spec_strategy_rollingupdate_max_surge",
			"Maximum number of replicas that can be scheduled above the desired number of replicas during a rolling update of a deployment.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				if d.Spec.Strategy.RollingUpdate == nil {
					return &metric.Family{}
				}

				maxSurge, err := intstr.GetScaledValueFromIntOrPercent(d.Spec.Strategy.RollingUpdate.MaxSurge, int(*d.Spec.Replicas), true)
				if err != nil {
					panic(err)
				}

				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(maxSurge),
						},
					},
				}
			}),
		),
		*generator.NewFamilyGeneratorWithStability(
			"kube_deployment_metadata_generation",
			"Sequence number representing a specific generation of the desired state.",
			metric.Gauge,
			basemetrics.STABLE,
			"",
			wrapDeploymentFunc(func(d *v1.Deployment) *metric.Family {
				return &metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.ObjectMeta.Generation),
						},
					},
				}
			}),
		),
	}
}

func wrapDeploymentFunc(f func(*v1.Deployment) *metric.Family) func(interface{}) *metric.Family {
	return func(obj interface{}) *metric.Family {
		deployment := obj.(*v1.Deployment)

		metricFamily := f(deployment)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys, m.LabelValues = mergeKeyValues(descDeploymentLabelsDefaultLabels, []string{deployment.Namespace, deployment.Name}, m.LabelKeys, m.LabelValues)
		}

		return metricFamily
	}
}

// mergeKeyValues merges label keys and values slice pairs into a single slice pair.
// Arguments are passed as equal-length pairs of slices, where the first slice contains keys and second contains values.
// Example: mergeKeyValues(keys1, values1, keys2, values2) => (keys1+keys2, values1+values2)
func mergeKeyValues(keyValues ...[]string) (keys, values []string) {
	capacity := 0
	for i := 0; i < len(keyValues); i += 2 {
		capacity += len(keyValues[i])
	}

	// Allocate one contiguous block, then split it up to keys and values zero'd slices.
	keysValues := make([]string, 0, capacity*2)
	keys = (keysValues[0:capacity:capacity])[:0]
	values = (keysValues[capacity : capacity*2])[:0]

	for i := 0; i < len(keyValues); i += 2 {
		keys = append(keys, keyValues[i]...)
		values = append(values, keyValues[i+1]...)
	}

	return keys, values
}
