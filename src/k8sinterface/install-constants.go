package k8sinterface

var repositoryName = "get-zen"
var repositoryLink = "https://github.com/get-zen-dev/tapp_store_rep"

var microk8sAddons = []string{"dns", "ingress", "cert-manager"}

var cfgCertsYml = `apiVersion: cert-manager.io/v1
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
