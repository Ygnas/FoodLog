# Prerequisites

- go version v1.20.0+
- docker
- kubectl
- Kubernetes cluster.

Before deploying the backend, ensure you have the following:

- `foodlog-credentials.json`: This file contains the necessary credentials for the backend to operate. You can find an example on GitHub, or you can download it from your Firebase project website.

- In the `foodlog-config.yaml`:
  Set up the following `JWT_SECRET` to your choosen secret for generating the JWT tokens and `DATABASE_URL` to your Firebase database url.

# Deployment

Deployment can be initiated using a single command in the terminal:

```
make
```

This command executes the Makefile, which performs the following steps:

Creates Kubernetes secrets using the provided credentials file.
Applies the Kustomization file, which generates the Config Map, Deployment, and a Service to make the application reachable from the outside.

With these steps completed, your backend application should be successfully deployed and accessible within your Kubernetes cluster.
