package resourceapply

import (
	"k8s.io/klog"

	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	storageclientv1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	storageclientv1beta1 "k8s.io/client-go/kubernetes/typed/storage/v1beta1"

	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
)

// ApplyStorageClass merges objectmeta, tries to write everything else
func ApplyStorageClass(client storageclientv1.StorageClassesGetter, recorder events.Recorder, required *storagev1.StorageClass) (*storagev1.StorageClass, bool,
	error) {
	existing, err := client.StorageClasses().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.StorageClasses().Create(required)
		reportCreateEvent(recorder, required, err)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := resourcemerge.BoolPtr(false)
	existingCopy := existing.DeepCopy()

	resourcemerge.EnsureObjectMeta(modified, &existingCopy.ObjectMeta, required.ObjectMeta)

	// Now that we copied everything that matters from required.ObjectMeta, we
	// should reset required.ObjectMeta so a DeepEqual() comparison is possible
	required.ObjectMeta = *existingCopy.ObjectMeta.DeepCopy()

	// We also need to reset required.TypeMeta because existingCopy.TypeMeta,
	// that's comming from the apiserver, isn't set (see https://issues.k8s.io/3030)
	required.TypeMeta = existingCopy.TypeMeta

	contentSame := equality.Semantic.DeepEqual(existingCopy, required)
	if contentSame && !*modified {
		return existingCopy, false, nil
	}

	objectMeta := existingCopy.ObjectMeta.DeepCopy()
	existingCopy = required.DeepCopy()
	existingCopy.ObjectMeta = *objectMeta

	if klog.V(4) {
		klog.Infof("StorageClass %q changes: %v", required.Name, JSONPatchNoError(existing, existingCopy))
	}

	// TODO if provisioner, parameters, reclaimpolicy, or volumebindingmode are different, update will fail so delete and recreate
	actual, err := client.StorageClasses().Update(existingCopy)
	reportUpdateEvent(recorder, required, err)
	return actual, true, err
}

// ApplyCSIDriverV1Beta1 merges objectmeta, does not worry about anything else
func ApplyCSIDriverV1Beta1(client storageclientv1beta1.CSIDriversGetter, recorder events.Recorder, required *storagev1beta1.CSIDriver) (*storagev1beta1.CSIDriver, bool,
	error) {
	existing, err := client.CSIDrivers().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.CSIDrivers().Create(required)
		reportCreateEvent(recorder, required, err)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}

	modified := resourcemerge.BoolPtr(false)
	existingCopy := existing.DeepCopy()

	resourcemerge.EnsureObjectMeta(modified, &existingCopy.ObjectMeta, required.ObjectMeta)
	if !*modified {
		return existingCopy, false, nil
	}

	if klog.V(4) {
		klog.Infof("CSIDriver %q changes: %v", required.Name, JSONPatchNoError(existing, existingCopy))
	}

	actual, err := client.CSIDrivers().Update(existingCopy)
	reportUpdateEvent(recorder, required, err)
	return actual, true, err
}
