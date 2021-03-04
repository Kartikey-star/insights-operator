package clusterconfig

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	dynamicfake "k8s.io/client-go/dynamic/fake"
)

func TestContainerRuntimeConfig(t *testing.T) {
	var machineconfigpoolYAML = `
apiVersion: machineconfiguration.openshift.io/v1
kind: ContainerRuntimeConfig
metadata:
    name: test-ContainerRC
`
	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme())
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	testContainerRuntimeConfigs := &unstructured.Unstructured{}

	_, _, err := decUnstructured.Decode([]byte(machineconfigpoolYAML), nil, testContainerRuntimeConfigs)
	if err != nil {
		t.Fatal("unable to decode machineconfigpool ", err)
	}
	_, err = client.Resource(containerRuntimeConfigGVR).Create(context.Background(), testContainerRuntimeConfigs, metav1.CreateOptions{})
	if err != nil {
		t.Fatal("unable to create fake machineconfigpool ", err)
	}

	ctx := context.Background()
	records, errs := gatherContainerRuntimeConfig(ctx, client)
	if len(errs) > 0 {
		t.Errorf("unexpected errors: %#v", errs)
		return
	}
	if len(records) != 1 {
		t.Fatalf("unexpected number or records %d", len(records))
	}
	if records[0].Name != "config/containerruntimeconfigs/test-ContainerRC" {
		t.Fatalf("unexpected containerruntimeconfig name %s", records[0].Name)
	}
}
