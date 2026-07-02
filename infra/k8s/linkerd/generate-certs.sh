#!/bin/bash
set -e

# 1. Create temporary directory for certs
mkdir -p certs

echo "Generating Linkerd Trust Anchor Certificate (CA)..."
# Generate Trust Anchor (CA)
openssl req -x509 -new -newkey rsa:4096 -keyout certs/ca.key -out certs/ca.crt -nodes -subj "/CN=root.linkerd.cluster.local" -days 3650

echo "Generating Linkerd Identity Issuer CSR..."
# Generate Issuer Key & CSR
openssl req -new -newkey rsa:2048 -keyout certs/issuer.key -out certs/issuer.csr -nodes -subj "/CN=identity.linkerd.cluster.local"

# Linkerd requires the basicConstraints extension set to CA:true, pathlen:0
echo "basicConstraints=CA:TRUE,pathlen:0" > certs/ext.txt
echo "keyUsage=critical,keyCertSign,cRLSign" >> certs/ext.txt

echo "Signing Issuer Certificate using Trust Anchor..."
# Sign Issuer Certificate using CA
openssl x509 -req -in certs/issuer.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/issuer.crt -days 365 -extfile certs/ext.txt

# Create namespace manifest first to make it a kustomize resource
cat <<EOF > namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: linkerd
EOF

echo "Creating issuer-secret.yaml..."
# Generate issuer-secret.yaml
cat <<EOF > issuer-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: linkerd-identity-issuer
  namespace: linkerd
type: kubernetes.io/tls
data:
  tls.crt: $(base64 -w0 < certs/issuer.crt)
  tls.key: $(base64 -w0 < certs/issuer.key)
EOF

echo "Creating values.yaml..."
# Generate values.yaml containing Trust Anchor PEM
cat <<EOF > values.yaml
identityTrustAnchorsPEM: |
$(cat certs/ca.crt | sed 's/^/  /')

identity:
  issuer:
    scheme: kubernetes.io/tls
EOF

# Cleanup temp folder
rm -rf certs

echo "================================================="
echo "Success! Linkerd certificates generated."
echo "Created: namespace.yaml"
echo "Created: issuer-secret.yaml"
echo "Created: values.yaml"
echo "================================================="
echo "You can now deploy Linkerd by running:"
echo "  kubectl apply -k ."
echo "================================================="
