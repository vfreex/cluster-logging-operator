package helpers

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cl "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	k8shandler "github.com/openshift/cluster-logging-operator/pkg/k8shandler"
	elasticsearch "github.com/openshift/elasticsearch-operator/pkg/apis/logging/v1"
)

type LogComponentType string

const (
	ComponentTypeStore         LogComponentType = "LogStore"
	ComponentTypeVisualization LogComponentType = "Visualization"
	ComponentTypeCollector     LogComponentType = "Collector"
)

func NewClusterLogging(componentTypes ...LogComponentType) *cl.ClusterLogging {
	instance := &cl.ClusterLogging{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterLogging",
			APIVersion: cl.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ClusterLoggingName,
			Namespace: OpenshiftLoggingNS,
			Annotations: map[string]string{
				k8shandler.ForwardingAnnotation: "enabled",
			},
		},
		Spec: cl.ClusterLoggingSpec{
			ManagementState: cl.ManagementStateManaged,
		},
	}

	for _, compType := range componentTypes {
		switch compType {
		case ComponentTypeStore:
			instance.Spec.LogStore = &cl.LogStoreSpec{
				Type: cl.LogStoreTypeElasticsearch,
				ElasticsearchSpec: cl.ElasticsearchSpec{
					Resources: &v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceMemory: resource.MustParse("3Gi"),
							v1.ResourceCPU:    resource.MustParse("100m"),
						},
					},
					NodeCount:        1,
					RedundancyPolicy: elasticsearch.ZeroRedundancy,
				},
			}
		case ComponentTypeCollector:
			instance.Spec.Collection = &cl.CollectionSpec{
				Logs: cl.LogCollectionSpec{
					Type: cl.LogCollectionTypeFluentd,
					FluentdSpec: cl.FluentdSpec{
						Resources: &v1.ResourceRequirements{
							Requests: v1.ResourceList{
								v1.ResourceMemory: resource.MustParse("736Mi"),
								v1.ResourceCPU:    resource.MustParse("500m"),
							},
						},
					},
				},
			}
		}
	}
	return instance
}
