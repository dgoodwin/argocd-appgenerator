# ArgoCD AppGenerator

An operator for dynamically generating Applications for ArgoCD clusters matching a given label selector.

ArgoCD Applications are defined with a 1-1 relationship to a cluster. In environments where clusters are treated as cattle, coming and going all the time, expecting a fixed Application to exist for each becomes more difficult to manage, yet ArgoCD is a great tool to deliver configuration to all clusters in your fleet.

Today, ArgoCD clusters are only represented as a Kubernetes Secret with a specific label. (argocd.argoproj.io/secret-type=cluster) There is no way to apply labels when running `argocd cluster add`, and as such you must manually add your labels to each cluster secret at this time.

An example cluster secret:

```yaml
apiVersion: v1
data:
  config: SNIP
  name: SNIP
  server: SNIP
kind: Secret
metadata:
  annotations:
    managed-by: argocd.argoproj.io
  creationTimestamp: "2020-04-09T16:31:06Z"
  labels:
    argocd.argoproj.io/secret-type: cluster
  name: cluster-kubernetes.default.svc-3396314289
  namespace: argocd
  resourceVersion: "100824"
```

The appgenerator CRD specifies the ApplicationSpec and which labels we should match on. Whenever a cluster secret is created/updated/deleted, we reconcile the appgenerator and create the desired ArgoCD Applications each tied to one specific cluster.

```yaml
apiVersion: appgenerator.rm-rf.ca/v1
kind: ApplicationGenerator
metadata:
  name: configmap-via-argocd-sample
  namespace: argocd
spec:
  matchLabels:
    argocd.argoproj.io/secret-type: cluster
  applicationSpec:
    destination:
      namespace: default
      # server will be overwritten for every cluster that matches this generator.
      server: https://kubernetes.default.svc
    project: default
    source:
      path: config/samples/argocd-application
      repoURL: https://github.com/dgoodwin/argocd-appgenerator.git
    syncPolicy:
      automated:
        prune: true
```

This repository is intended to be a proof-of-concept to eventually be proposed for inclusion in ArgoCD, potentially with discussion around a couple possible improvements:

  1. Represent Clusters with a concete ArgoCD CRD.
  1. Allow specifying labels on Clusters when creating via the CLI. (or UI if that is possible)

## Current Status

The AppGenerator will create applications for any cluster secrets which match your defined labels, update them if necessary, cleanup orphaned apps, and cleanup all apps if the owning appgenerator is deleted.
