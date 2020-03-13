// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/bertinatto/aws-ebs-csi-driver-operator/pkg/generated/clientset/versioned/typed/operator/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeCsiV1alpha1 struct {
	*testing.Fake
}

func (c *FakeCsiV1alpha1) EBSCSIDrivers() v1alpha1.EBSCSIDriverInterface {
	return &FakeEBSCSIDrivers{c}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeCsiV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
