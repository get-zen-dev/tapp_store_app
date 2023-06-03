package main

import (
	he "handleException"
	"io"
	k8 "k8sinterface"
	"os/exec"
	"strconv"
)

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

func main() {

	//_, err = file.Write(data)
	//if err != nil {
	//	return err
	//}

	cmds := []*exec.Cmd{exec.Command("microk8s", "kubectl", "apply", "-f", "-"),
		exec.Command("microk8s", "kubectl apply", "-f", "-"),
		exec.Command("microk8s kubectl apply", "-f", "-")}

	for i, val := range cmds {
		println("==========#" + strconv.Itoa(i) + "#=============")
		reader, writer := io.Pipe()
		val.Stdin = reader
		//cmd.Stdout = os.Stdout

		go func() {
			defer writer.Close()
			writer.Write([]byte(cfg_certs_yml))
		}()

		out, err2 := val.Output()
		println(string(out))
		if err2 != nil {
			println("Error: " + string(out))
		}
		err := err2
		if err != nil {
			println("Other ERROR:" + err.Error())
		}
		println("==========#" + strconv.Itoa(i) + "#=============")
	}

	//cmd := exec.Command("microk8s", "kubectl", "apply", "-f", "-")
	//
	//err := invoke(cmd)
	//println("Other error:" + err.Error())

	//file, err := os.Create("temp.yml")
	//file.Chmod(0777)
	//data := []byte(cfg_certs_yml)
	//_, err = file.Write(data)
	//if err == nil {
	//	println("Success Write")
	//}
	//
	//cmd := exec.Command("cat", file.Name())
	//out, err := cmd.Output()
	//println(string(out))
	//if err == nil {
	//	println("Success 1")
	//}

	clientMicrok8s, err := k8.GetInterfaceProvider("domain")
	if err != nil {
		he.PrintErr(err)
	}

	err = clientMicrok8s.Start()
	if err != nil {
		he.PrintErr(err)
	} else {
		println("Succesfull install system")
	}

	//if !k8.CheckIsRootGranted() {
	//	he.PrintErr(fmt.Errorf("cannot user interface without root privileges"))
	//}
	//domain, err := env.GetDomain()
	//if err != nil {
	//	q, err := view.NewModelQuestion()
	//	he.PrintErrorIfNotNil(err)
	//	p := tea.NewProgram(q, tea.WithAltScreen())
	//	if _, err := p.Run(); err != nil {
	//		he.PrintErr(err)
	//	}
	//} else {
	//	clientMicrok8s, err := k8.GetInterfaceProvider(domain)
	//	he.PrintErrorIfNotNil(err)
	//	w := view.NewModelWaiting(clientMicrok8s, view.KubernetesLaunch)
	//	p := tea.NewProgram(w, tea.WithAltScreen())
	//	if _, err := p.Run(); err != nil {
	//		he.PrintErr(err)
	//	}
	//}
}
