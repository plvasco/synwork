# 🧩 Bitbucket Runner Deployment on AWS EKS (with Ansible)

This project uses **Ansible** to automate the deployment of [Bitbucket Pipelines Runners](https://support.atlassian.com/bitbucket-cloud/docs/manage-self-hosted-runners-for-linux/) onto an **Amazon EKS** cluster. It also includes a role to configure **ECR (Elastic Container Registry)** repositories.

---

## 📦 Prerequisites

- Python & Ansible installed:
  ```bash
  pip install ansible
  ```

- Install the Kubernetes Ansible collection:
  ```bash
  ansible-galaxy collection install kubernetes.core
  ```

- AWS CLI v2 installed and authenticated.

- AWS IAM user must have sufficient permissions (e.g. EKS access, ECR management).

---

## ⚙️ Configuration

### 🔧 Bitbucket Runner Setup

1. Go to your Bitbucket repository:
   - `Repository settings → Runners → Add Runner`
   - Copy the values shown (runner name, UUIDs, OAuth credentials).

2. update values in playbook.yaml:

   ```yaml
   runner_name: "runner-abc123"
   runner_namespace: "bitbucket-runner"
   account_uuid: "xxxxxxxxx"
   runner_uuid: "xxxxxxxxx"
   repository_uuid: "xxxxxxxxx"

   # Environment-specific secrets
   oauth_client_id: "your-client-id"
   oauth_client_secret: "your-client-secret"
   ```

---

## 🚀 Deployment Instructions

### 🧠 Set Kubeconfig

Ensure your local kubeconfig is pointing to the target EKS cluster:

```bash
export K8S_AUTH_KUBECONFIG=~/.kube/config
```

### ▶️ Run the Ansible Playbook

#### Deploy bitbucket runner

```bash
ansible-playbook playbook.yaml
```

---

## 🧪 Troubleshooting

- ❗ **Unauthorized (401)**:
  - Ensure your AWS profile has access to the EKS cluster.
  - Check that `~/.kube/config` is pointing to the correct cluster.

- ❗ **Playbook variables not overriding role defaults**:
  - Move variables from `roles/bitbucket-runner/vars/main.yml` ➝ `roles/bitbucket-runner/defaults/main.yml`.
  - Or override with `-e` on the CLI as shown above.

- 🔍 Debug EKS token value:
  ```yaml
  - name: Debug EKS token
    debug:
      var: eks_api_token
  ```

---

## 📚 References

- [Bitbucket Runner Docs](https://support.atlassian.com/bitbucket-cloud/docs/manage-self-hosted-runners-for-linux/)
- [Ansible Kubernetes Collection](https://docs.ansible.com/ansible/latest/collections/kubernetes/core/index.html)
- [AWS EKS Authentication](https://docs.aws.amazon.com/eks/latest/userguide/authenticate.html)

---

## 📄 License

MIT License – use freely and modify as needed.
