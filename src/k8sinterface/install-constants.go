package k8sinterface

import "time"

const repositoryName = "get-zen"
const repositoryLink = "https://github.com/get-zen-dev/tapp_store_rep"

const microk8sVersion = "--channel=1.25/stable"

var microk8sAddons = []string{"dns", "ingress", "cert-manager"}

const commandCore = "microk8s"

const awaitKubeInitializeForCertTimeout = time.Second * 10
const iterationRetryTimeout = time.Second

const cfgCertsYml = `apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-certs
spec:
  acme:
    server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-certs
    solvers:
    - http01:
        ingress:
          class: public
`
