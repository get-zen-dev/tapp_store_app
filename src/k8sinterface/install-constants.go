package k8sinterface

var microk8s_addons = []string{"dns", "ingress", "cert-manager"}

var cfg_certs_yml = `apiVersion: cert-manager.io/v1
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
