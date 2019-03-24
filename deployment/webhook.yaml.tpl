---
apiVersion: admissionregistration.k8s.io/v1beta1
kind: ValidatingWebhookConfiguration
metadata:
  name: admission-webhook
webhooks:
  - name: k8s-admission-webhook.mafr.github.com
    clientConfig:
      service:
        name: admission-webhook
        namespace: kube-admission
        path: "/validate"
      caBundle: ${CA_BUNDLE}
    rules:
      - operations: ["CREATE","UPDATE"]
        apiGroups: ["", "apps"]
        apiVersions: ["v1"]
        resources: ["deployments", "services"]
    failurePolicy: Ignore
