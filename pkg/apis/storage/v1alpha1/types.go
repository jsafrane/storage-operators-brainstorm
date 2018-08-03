package v1alpha1

import (
	openshiftapi "github.com/openshift/api/operator/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Common storage operator status inlined in most of the status types below
type StorageStatus struct {
	//// Open question: Use whole openshiftapi.OperatorStatus???
	//// Why not:
	////  - unclear Conditions. What SyncSuccessful means?
	////  - CurrentAvailability vs TargetAvailability
	////    - Why both have GenerationHistory? There is only one Deployment with single GenerationHistory
	////    - What CurrentAvailability means? We update version of an image in a deployment. It starts rolling new pods.
	////      What is Current and what is Target? Is the old version still "current", if it has only one pod and 99 other
	////      pods are at the new version?

	// ObservedGeneration is the last generation of this object that the operator has acted on.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Generation of API objects that the operator has created / updated. We store it here, because the operator
	// needs to persist the last generations of objects across operator restart.
	ChildrenGenerations []openshiftapi.GenerationHistory
}

//
// EFS Provisioner Operator
//

type EFSProvisioner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              EFSProvisionerSpec   `json:"spec"`
	Status            EFSProvisionerStatus `json:"status,omitempty"`
}

type EFSProvisionerSpec struct {
	// Name of image (incl. version) with EFS provisioner. Used to override global image and version from cluster config.
	// It should be empty in the usual case.
	//// TODO: find out where the cluster config is.
	ProvisionerImage string

	// Name of storage class to create. If the storage class already exists, it will not be updated.
	//// This allows users to create their storage classes in advance.
	// Mandatory, no default.
	StorageClassName string

	// Location of AWS credentials. Used to override global AWS credential from cluster config.
	// It should be empty in the usual case.
	//// TODO: Where can we find the cluster credentials?
	AWSSecrets v1.SecretReference

	// ID of the EFS to use as base for dynamically provisioned PVs.
	// Such EFS must be created by admin before starting a provisioner!
	// Mandatory, no default
	FSID string

	// Subdirectory on the EFS specified by FSID that should be used as base
	// of all dynamically provisioner PVs. Root of the EFS will be used when
	// BasePath is not set.
	// Optional, no default.
	BasePath string

	// Group that can write to the EFS. The provisioner will run with this
	// supplemental group to be able to create new PVs.
	// Optional, no default.
	SupplementalGroup int64
}

type EFSProvisionerStatus struct {
	StorageStatus `json:",inline"`

	//// TODO: add UpdatedReplicas, ReadyReplicas and such?
}

//
// Manila provisioner
//

type ManilaProvisioner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ManilaProvisionerSpec   `json:"spec"`
	Status            ManilaProvisionerStatus `json:"status,omitempty"`
}

type ManilaProvisionerSpec struct {
	// Name of image (incl. version) with manila provisioner. Used to override global image and version from cluster config.
	// It should be empty in the usual case.
	//// TODO: find out where the cluster config is.
	ProvisionerImage string

	// Name of storage class to create. If the storage class already exists, it will not be updated.
	//// This allows users to create their storage classes in advance.
	// Mandatory, no default.
	StorageClassName string

	// Location of OpenStack credentials. Mandatory, no default.
	//// Open question: should it be optional and we'll use the same credentials as the cluster?
	//// Where can we find the cluster credentials?
	OpenStackSecrets v1.SecretReference

	// Parameters of the storage class.
	// See https://github.com/kubernetes/cloud-provider-openstack/blob/master/docs/using-manila-provisioner.md#share-options
	//// TODO: decompose to individual fields! We already have a field for osSecretName.
	//Parameters map[string]string
}

type ManilaProvisionerStatus struct {
	StorageStatus `json:",inline"`

	//// TODO: add UpdatedReplicas, ReadyReplicas and such?
}

//
// CephFS provisioner
//

type CephFSProvisioner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CephFSProvisionerSpec   `json:"spec"`
	Status            CephFSProvisionerStatus `json:"status,omitempty"`
}

type CephFSProvisionerSpec struct {
	openshiftapi.OperatorSpec `json:",inline"`

	// Name of storage class to create. If the storage class already exists, it will not be updated.
	//// This allows users to create their storage classes in advance.
	// Mandatory, no default.
	StorageClassName string

	// Location of CephFS admin credentials. The secret must have two keys, "username" and "password".
	// Mandatory, no default.
	CephFSSecrets v1.SecretReference

	//// From https://github.com/kubernetes-incubator/external-storage/blob/fd2a6f5805d17b80172d2c89a61e881c3edfb2ab/ceph/cephfs/cephfs-provisioner.go#L280
	// Name of the Ceph cluster
	ClusterName string

	// List of Ceph monitors, i.e. their addresses (+ port numbers, if necessary).
	Monitors []string

	// Subdirectory on the Ceph volume that should be used as base
	// of all dynamically provisioner PVs. Root of the volume will be used when
	// BasePath is not set.
	// Optional, no default.
	BasePath string

	//// deterministicnames will be forced to true, no configuration.

	// Namespace where to put secrets for PVs. When empty, secrets will be put into the same namespaces where
	// PVC for the provisioned volume reside. Optional, no default.
	// TODO: isn't it insecure? Should we put the secrets to the namespace where the operator runs?
	CreatedSecretsNamespace string

	// Whether quota should be enforced on created PVs.
	//// TODO: Is there any reason why this should be configurable?
	EnableQuota bool
}

type CephFSProvisionerStatus struct {
	StorageStatus `json:",inline"`

	//// TODO: add UpdatedReplicas, ReadyReplicas and such?
}

//
// Snapshot.
//

//// TODO: better name. This CRD defines both external snapshot controller and provisioner, not just controller.
type SnapshotController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              SnapshotControllerSpec   `json:"spec"`
	Status            SnapshotControllerStatus `json:"status,omitempty"`
}

type SnapshotControllerSpec struct {
	openshiftapi.OperatorSpec `json:",inline"`
	//// TODO: fill it
	//// Should there be only one CR for whole cluster? In that case we need multiple storage classes, credentials
	//// (e.g. AWS and Gluster) and so on in one CR.
	//// Shall we move to CSI? It could be easier there.

	// Name of Kubernetes user group that will have access to snapshots, i.e. can create / restore / delete snapshots
	// of PVCs.
	// Empty value means that the operator won't authorize anyone and it's up to system admin to create own RBAC rules.
	Group string
}

type SnapshotControllerStatus struct {
	StorageStatus `json:",inline"`

	//// TODO: add UpdatedReplicas, ReadyReplicas and such?
}

//
// Local storage.
//

//// Local storage is more complicated. Admin must group nodes with the same storage devices and label them.
//// Then the admin create LocalStorageProvider for each of this group.
//// The operator will run separate DaemonSet with local storage provisioner for each LocalStorageProvider.

// Defines local devices on set of nodes grouped by NodeSelector.
type LocalStorageProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              LocalStorageProvider `json:"spec"`
	Status            LocalProviderStatus  `json:"status,omitempty"`
}

type LocalStorageProviderSpec struct {
	openshiftapi.OperatorSpec `json:",inline"`

	// Selector that matches nodes whose local storage is managed by the
	// LocalStorageProvider. All nodes are matched when NodeSelector is not
	// set.
	NodeSelector *v1.NodeSelector `json:"nodeSelector,omitempty"`

	// List of storage classes and their devices on the matched nodes.
	StorageClassDevices []StorageClassDevices `json:"storageClassDevices"`
}

type StorageClassDevices struct {
	// Storage class where the listed devices belong to.
	StorageClassName string `json:"storageClassName"`

	// List of kernel names of devices that belong the storage class. Without
	// any "/dev" prefix, e.g. "sda" or "vdb".
	// At least one of deviceNames or deviceStableNames must be set.
	DeviceNames []string `json:"deviceNames,omitempty"`

	// List of /dev/disk/by-id/ names of devices that belong to the storage
	// class. Without any "/dev/disk/by-id" prefix, e.g.
	// "ata-SAMSUNG_MZ7LN512HMJP-000L7_S2X9NX0H123456".
	// At least one of deviceNames or deviceStableNames must be set.
	DeviceStableNames []string `json:"deviceStableNames,omitempty"`

	// TODO
	// UdevRules []string `json:"udevRules"`
}

type LocalProviderStatus struct {
	StorageStatus `json:",inline"`

	//// TODO: add UpdatedReplicas, ReadyReplicas and such?
}

//
// CSI driver.
//

type CSIDriverDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              CSIDriverDeploymentSpec   `json:"spec"`
	Status            CSIDriverDeploymentStatus `json:"status,omitempty"`
}

type CSIDriverDeploymentSpec struct {
	// Name of the CSI driver.
	DriverName string

	// Template of pods that will run on every node. It must contain a container with the driver and all in volumes
	// it needs (Secrets, ConfigMaps). Sidecar with driver registrar and liveness probe will be added by the operator.
	NodeTemplate v1.PodTemplate

	// Selector of nodes where to run DaemonSet with the node driver. All nodes are used when the selector is empty.
	NodeSelector v1.NodeSelector

	// Template of pods that will run the controller parts. Nil when the driver does not require any attacher or
	// provisioner. Sidecar with provisioner and attacher will be added by the operator.
	ControllerTemplate *v1.PodTemplate

	// Path to CSI socket in the container with CSI driver. In case nodeTemplate or controllerTemplate have more
	// containers, this is the *first* container.
	DriverSocket string

	//// The operator will do with nodeTemplate:
	//// - add necessary driver registrar and liveness probe sidecar containers
	//// - add hostPath to the first container in nodeTemplate at Dir(DriverSocket) to export the driver socket
	//// somewhere to /var/lib/kubelet/plugins/<driver name>/.
	//// - runs Deploymnent with the controllers.
	//// - runs DaemonSet with the node driver.

	//// On update (NodeTemplate or ControllerTemplate), we update DaemonSet + Deployment and we expect rolling update.
	//// Old pod on a node will be deleted (its CSI driver socket closed), all CSI volume plugin for this driver will
	//// fail. A new pod with the driver will be started and will open new CSI driver socket.
	//// TODO: this will break gluster. When we stop the old pod, it will kill all gluster fuse mounts on the node.
	//// Shall we update the nodes one by one and drain the node first?

	//// TODO:
	// NodeUpdateStrategy string // "Rolling" or "Drain"?
}

type CSIDriverDeploymentStatus struct {
	StorageStatus `json:",inline"`

	//// TODO: add UpdatedReplicas, ReadyReplicas and such?
}
