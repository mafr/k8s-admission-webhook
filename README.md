# Kubernetes Admission Webhook

Kubernetes' Admission Controller supports webhooks for validating admission
requests. This repository contains some code that can be used to build such
a webhook. Use at your own risk, this is very much work in progress.


## Deployment

Deploying the webhook is quite a bit of work because Kubernetes will only
accept the webhook if it uses TLS with a certificate signed by the
api-server's own CA. Follow this process:

  1. Create a dedicated namespace for the webhook
  2. Create a key and a certificate signing request (CSR)
  3. Submit the CSR to Kubernetes
  4. Have Kubernetes approve the CSR
  5. Retrieve the signed certificate from Kubernetes
  6. Deploy a Secret into the namespace, containing both key and certificate

The first step we execute manually:

    $ kubectl create namespace kube-admission

For the rest we use [a script](https://github.com/istio/istio/blob/release-0.7/install/kubernetes/webhook-create-signed-cert.sh) from an old Istio release:

    $ ./webhook-create-signed-cert.sh --namespace kube-admission --service admission-webhook --secret kube-admission

Confirm that the secret has been created:

    $ kubectl -n kube-admission get secret kube-admission
    NAME             TYPE     DATA   AGE
    kube-admission   Opaque   2      2m
    $

At this point, you can delete the certificate and key that the script placed
inside a temporary directory in `/tmp/'.

After the infrastructure setup, we deploy the webhook itself:

  1. Create a Deployment, with Pods mounting the Secret created before
  2. Create a Service for the Admission Controller to call
  3. Register the webhook with Kubernetes' Admission Controller

Step one and two are simple:

    $ kubectl apply -f deployment/deployment.yaml
    $ kubectl apply -f deployment/service.yaml

The webhook configuration has to contain the CA's root certificate,
which we have to extract from Kubernetes:

    $ cat deployment/webhook.yaml.tpl | scripts/add-ca-bundle.sh | kubectl apply -f -

Now the webhook should be active. It will validate all resource types
specified in its configuration.


## References

The code in this repository is heavily inspired by these blog posts:

  * [In-depth Introduction to Kubernetes Admission Webhooks](https://banzaicloud.com/blog/k8s-admission-webhooks/)
  * [A Gentle Intro to Validation Admission Webhooks in Kubernetes](https://container-solutions.com/a-gentle-intro-to-validation-admission-webhooks-in-kubernetes/)
