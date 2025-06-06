# Storage Profiles

## Introduction

Storage Profile is the resource that serves the information about recommended parameters for the PVC. 

This can be used by CDI controllers when creating a PVC for DV. That way the DataVolume can be simplified and if the properties are missing,
defaults can be applied from the StorageProfile.

PVC created independently (without DV) with missing fields, can also be auto-completed with defaults from the StorageProfiles. This is achieved using [CDI PVC mutating webhook rendering](./pvc-mutating-webhook-rendering.md).

CDI provides a collection of Storage Profiles with default recommended values for some well known backends. 
If the storage provisioner defined in storage class does not have defaults configured in CDI the resulting StorageProfile 
has empty `claimPropertySets`.

CDI automatically creates the StorageProfile objects - one StorageProfile per 
one StorageClass that exists in the cluster. Below StorageProfile for hostpath-provisioner as an example.

```yaml
apiVersion: cdi.kubevirt.io/v1beta1
kind: StorageProfile
metadata: 
  name: hostpath-provisioner
spec:
  claimPropertySets: 
  - accessModes:
    - ReadWriteOnce
    volumeMode: 
      Filesystem
  cloneStrategy: snapshot
status:
    storageClass: hostpath-provisioner
    provisioner: kubevirt.io/hostpath-provisioner
    claimPropertySets: 
    - accessModes: 
      - ReadWriteOnce
      volumeMode: Filesystem
    cloneStrategy: snapshot
```

### Parameters
- `cloneStrategy` - defines the preferred method for performing a CDI clone
- `claimPropertySets` contains a list of `claimPropertySet`
  - `accessMode` - contains the desired access modes the volume should have
  - `volumeMode` - defines what type of volume is required by the claim  
  These are ordered by preference, and are considered against existing partial info in the storage stanza of the datavolume.  
  Some preference considerations to note are:  
      - Block is preferred over Filesystem for performance reasons (fewer layers)  
      - ReadWriteMany over ReadWriteOnce (live migration support)
- `dataImportCronSourceFormat` DataImportCron (recurring polling of golden registry sources) was originally designed to only maintain PVC sources, However, for certain storage types, we know that snapshots sources scale better. Some details and examples can be found in [clone-from-volumesnapshot-source](./clone-from-volumesnapshot-source.md).

Values for accessModes and volumeMode are exactly the same as for PVC: `accessModes` is a list of `[ReadWriteMany|ReadWriteOnce|ReadOnlyMany]`.  
We are aware of `ReadWriteOncePod` but [currently](https://github.com/kubevirt/containerized-data-importer/issues/2365) are not testing it.  
`volumeMode` is a single value `Filesystem` or `Block`.
Multiple claim property sets can be specified (`claimPropertySets` is a list).

The value for `cloneStrategy` can be one of:
- `copy` - copy blocks of data over the network
- `snapshot` - clones the volume by creating a temporary VolumeSnapshot and restoring it to a new PVC
- `csi-clone` - clones the volume using a CSI clone

When the value is not specified the CDI will try to use the `snapshot` if possible otherwise it falls back to `copy`. 
If the storage class (and its provider) is capable of doing CSI Volume Clone then the user may choose `csi-clone` as a preferred clone method.  
`csi-clone` is preferred in general, since it offloads the optimization responsibility to the storage provider.

StorageClass can be annotated with `cdi.kubevirt.io/clone-strategy`. The annotation value can be one of: `copy`,`snapshot`,`csi-clone`.
CDI is using this annotation value when configuring the clone strategy on storage profile. 
This is helpful for known provisioners that want different behavior for certain configurations in the storage class 


## Handling the DV with defaults from Storage Profiles 

The example uses the `hpp` (`kubevirt.io/hostpath-provisioner`) as the storage provisioner.
For brevity some fields managed by kubernetes, like managedFields or creationTimestamp, were removed from output. 

1. Given the `hostpath-provisioner` StorageClass, CDI creates a `hostpath-provisioner` StorageProfile

`kubectl get sc hostpath-provisioner -o yaml`
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  name: hostpath-provisioner
provisioner: kubevirt.io/hostpath-provisioner
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```

`k get storageprofile hostpath-provisioner -o yaml`
```yaml
apiVersion: cdi.kubevirt.io/v1beta1
kind: StorageProfile
metadata:
  labels:
    app: containerized-data-importer
    cdi.kubevirt.io: ""
  name: hostpath-provisioner
  ownerReferences: 
    ...
spec: {}
status:
  claimPropertySets:
  - accessModes:
    - ReadWriteOnce
    volumeMode: Filesystem
  provisioner: kubevirt.io/hostpath-provisioner
  storageClass: hostpath-provisioner
```

2. Now the user can create a new DV using new Storage type, without specifying accessModes or volumeMode for the PVC. 
The `storage` field is the direct replacement of the `pvc` field.

`cat dv.yaml`
```yaml
apiVersion: cdi.kubevirt.io/v1beta1
kind: DataVolume
metadata:
  name: blank-dv
spec:
  storage:
    resources:
      requests:
        storage: 1Gi
    storageClassName: hostpath-provisioner
  source:
    blank: {}
```
`kubectl create -f dv.yaml`

Notice `pvc` replaced with `storage`, and both accessModes and volumeMode missing from `*.spec.pvc` on DataVolume.
```yaml
accessModes:
- ReadWriteOnce
volumeMode: Filesystem
```

3. As a result the following pvc is created.
   
`kubectl get pvc blank-dv -o yaml`

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  annotations:
    cdi.kubevirt.io/storage.condition.running.reason: Completed
    ...
  creationTimestamp: "2021-03-04T14:19:33Z"
  finalizers:
  - kubernetes.io/pvc-protection
  labels:
    app: containerized-data-importer
  name: blank-dv
  namespace: default
  ownerReferences: 
    ...
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: hostpath-provisioner
  volumeMode: Filesystem
  volumeName: pvc-a1c62357-dbfd-4909-aed7-19a88fa1e643
status:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 2Gi
  phase: Bound
```
Notice how accessModes is ReadWriteOnce and volumeMode is Filesystem, exactly as configured in the Storageprofile.

## Empty Storage Profile

Not all provisioners have recommended parameters provided by CDI. In a case where no recommendation is available, CDI creates an empty Storage Profile.

`kubectl get storageprofile some-unknown-provisioner -o yaml`
```yaml
apiVersion: cdi.kubevirt.io/v1beta1
kind: StorageProfile
metadata:
  name: some-unknown-provisioner-class
    ...
spec: {}
status:
  provisioner: some-unknown-provisioner
  storageClass: some-unknown-provisioner-class
```
There are no recommended parameters on StorageProfile so it is not possible to create a PVC for a DV without accessModes configured.   

`kubectl describe dv blank-dv`
```yaml
Name:         blank-dv
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  cdi.kubevirt.io/v1beta1
Kind:         DataVolume
Metadata:
 ...
Spec:
  Pvc:
    Resources:
      Requests:
        Storage:         1Gi
    Storage Class Name:  local
  Source:
    Blank:
Events:
  Type     Reason            Age                From                   Message
  ----     ------            ----               ----                   -------
  Warning  ErrClaimNotValid  1s (x12 over 18s)  datavolume-controller  DataVolume.storage spec is missing accessMode and cannot get access mode from StorageProfile local
```

Notice the event on the DV.

## User defined Storage Profile

User with access rights to edit StorageProfile can configure recommended parameters. Edit spec section of StorageProfile by adding claimPropertySets with accessModes and volumeMode.
When editing volumeMode you must also configure accessModes.
Shortly, all provided parameters should be visible in the status section. User defined parameter has higher priority and overrides the one provided by CDI. 

## Priorities

1. Overrides (for example `cdi.Spec.CloneStrategyOverride`)
2. Parameter defined on DataVolume
3. User provided parameters - defined on StorageProfile spec section.
4. Parameters provided by CDI.
5. Empty or kubernetes defaults (if available).
