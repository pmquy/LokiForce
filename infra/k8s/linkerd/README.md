# Linkerd Service Mesh Local Deployment

This directory contains the configurations and scripts to deploy Linkerd on your local Kubernetes cluster.

## Deployment Steps

1. **Generate Certificates:**
   Linkerd requires mutual TLS (mTLS) to secure traffic between services. Run the helper script to generate the CA and issuer certificates using `openssl`:
   ```bash
   chmod +x generate-certs.sh
   ./generate-certs.sh
   ```
   This will generate:
   - `namespace.yaml`: Standard namespace manifest.
   - `issuer-secret.yaml`: Securely holds the private keys (should NOT be committed to git).
   - `values.yaml`: Connects the CA to the Helm deployment.

2. **Deploy Linkerd:**
   Ensure Kustomize Helm integration is enabled on your cluster client, then run:
   ```bash
   kubectl apply -k .
   ```

3. **Verify Installation:**
   Check that all pods in the `linkerd` namespace are running:
   ```bash
   kubectl get pods -n linkerd
   ```
