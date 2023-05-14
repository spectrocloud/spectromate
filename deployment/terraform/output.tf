# Copyright (c) Spectro Cloud
# SPDX-License-Identifier: Apache-2.0

output "default_message" {
  value = "Use the following command to create the kubeconfig file lcoally:\n\nterraform output -raw kubeconfig > c1.config && export KUBECONFIG=$(pwd)/c1.config && \\\nkubectl get svc -n spectromate-app-spectromate-ns spectromate-app-spectromate-svc -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'\\\n"
}


output "kubeconfig" {
  value = spectrocloud_virtual_cluster.cluster-1.kubeconfig

  sensitive = true

  depends_on = [
    spectrocloud_virtual_cluster.cluster-1
  ]
}