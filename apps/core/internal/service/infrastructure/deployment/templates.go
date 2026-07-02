package deployment

var TemplatesMap = map[string]string{
	"argocd.yaml":        argocdTemplate,
	"deployment.yaml":    deploymentTemplate,
	"service.yaml":       serviceTemplate,
	"kustomization.yaml": kustomizationTemplate,
}

const argocdTemplate = `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{.ServiceName}}
  namespace: argocd
  annotations:
    argocd-image-updater.argoproj.io/image-list: app=ghcr.io/{{.Owner}}/{{.ServiceName}}
    argocd-image-updater.argoproj.io/app.update-strategy: latest
    argocd-image-updater.argoproj.io/write-back-method: git:kustomize
    argocd-image-updater.argoproj.io/write-back-path: manifests/{{.ServiceName}}
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  project: default
  source:
    repoURL: 'https://github.com/{{.Owner}}/{{.GitOpsRepo}}.git'
    targetRevision: HEAD
    path: manifests/{{.ServiceName}}
  destination:
    server: 'https://kubernetes.default.svc'
    namespace: {{.Namespace}}
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
`

const deploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ServiceName}}
  namespace: production
  labels:
    app: {{.ServiceName}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.ServiceName}}
  template:
    metadata:
      labels:
        app: {{.ServiceName}}
    spec:
      containers:
        - name: app
          image: ghcr.io/{{.Owner}}/{{.ServiceName}}:latest
          ports:
            - containerPort: 8080
          resources:
            limits:
              cpu: "500m"
              memory: "512Mi"
            requests:
              cpu: "100m"
              memory: "128Mi"
`

const serviceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: {{.ServiceName}}
  namespace: production
  labels:
    app: {{.ServiceName}}
spec:
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
  selector:
    app: {{.ServiceName}}
  type: ClusterIP
`

const kustomizationTemplate = `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - deployment.yaml
  - service.yaml
images:
  - name: ghcr.io/{{.Owner}}/{{.ServiceName}}
    newTag: latest
`

const networkPolicyTemplate = `apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: policy-{{.PolicyID}}
  namespace: {{.Namespace}}
spec:
  podSelector:
    matchLabels:
      app: {{.TargetID}}
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: {{.ClientID}}
      ports:
        - port: {{.TargetPort}}
          protocol: TCP
`
