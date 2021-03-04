package clusterconfig

import (
	"context"
	"fmt"

	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"github.com/openshift/insights-operator/pkg/record"
)

// GatherCRD collects the specified Custom Resource Definitions.
//
// The following CRDs are gathered:
// - volumesnapshots.snapshot.storage.k8s.io (10745 bytes)
// - volumesnapshotcontents.snapshot.storage.k8s.io (13149 bytes)
//
// The CRD sizes above are in the raw (uncompressed) state.
//
// Location in archive: config/crd/
func GatherCRD(g *Gatherer) func() ([]record.Record, []error) {
	return func() ([]record.Record, []error) {
		crdClient, err := apixv1beta1client.NewForConfig(g.gatherKubeConfig)
		if err != nil {
			return nil, []error{err}
		}
		return gatherCRD(g.ctx, crdClient)
	}
}

func gatherCRD(ctx context.Context, crdClient apixv1beta1client.ApiextensionsV1beta1Interface) ([]record.Record, []error) {
	toBeCollected := []string{
		"volumesnapshots.snapshot.storage.k8s.io",
		"volumesnapshotcontents.snapshot.storage.k8s.io",
	}
	records := []record.Record{}
	for _, crdName := range toBeCollected {
		crd, err := crdClient.CustomResourceDefinitions().Get(ctx, crdName, metav1.GetOptions{})
		// Log missing CRDs, but do not return the error.
		if errors.IsNotFound(err) {
			klog.V(2).Infof("Cannot find CRD: %q", crdName)
			continue
		}
		// Other errors will be returned.
		if err != nil {
			return []record.Record{}, []error{err}
		}
		records = append(records, record.Record{
			Name: fmt.Sprintf("config/crd/%s", crd.Name),
			Item: record.JSONMarshaller{Object: crd},
		})
	}
	return records, []error{}
}
