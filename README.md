# Kubernetes Admission Webhook

Kubernetes' Admission Controller supports webhooks for validating admission
requests. This repository contains some code that can be used to build such
a webhook. Use at your own risk, this is very much work in progress.


## Deployment

Deploying the webhook is quite a bit of work because Kubernetes will only
accept the webhook if it uses TLS with a certificate signed by the
api-server's own CA. Follow this process:

  # Create a key and a certificate signing request (CSR)
  # Submit the CSR to Kubernetes
  # Have Kubernetes approve the CSR
  # Retrieve the signed certificate from Kubernetes

Make sure to safeguard the certificate and key, they are highly security
critical. Next, set up the namespace:

  # Create a dedicated namespace for the webhook
  # Deploy a Secret into the namespace, containing both key and certificate

After the infrastructure setup, we deploy the webhook:

  # Create a Deployment, with Pods mounting the Secret created before
  # Create a Service for the Admission Controller to call
  # Register the webhook with Kubernetes' Admission Controller


## References

The code in this repository is heavily inspired by these blog posts:

  * [In-depth introduction to Kubernetes admission webhooks](https://banzaicloud.com/blog/k8s-admission-webhooks/)
  * [A Gentle Intro to Validation Admission Webhooks in Kubernetes](https://container-solutions.com/a-gentle-intro-to-validation-admission-webhooks-in-kubernetes/)
