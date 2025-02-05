# Install ansible role
ansible-galaxy collection install kubernetes.core

# Add your Bitbucket details in role/bitbucket-runner/vars/main.yml

# Make sure your aws user has proper access on eks cluster as it will be used in kubernetes by ansible.

# Set kubeconfig path as envrionment varaible.
export K8S_AUTH_KUBECONFIG=~/.kube/config

# Deploy to an env.
# To dev
ansible-playbook -i inventory/dev playbook.yaml -e @roles/bitbucket-runner/vars/main.yaml

# To staging
ansible-playbook -i inventory/staging playbook.yaml -e @roles/bitbucket-runner/vars/staging.yaml

# To prod
ansible-playbook -i inventory/staging playbook.yaml -e @roles/bitbucket-runner/vars/staging.yaml



# Deploy ecr role

## Update roles/aws_ecr/vars/main.yml and add ecr repo names and related policies if required.

# Run pipeline.
ansible-playbook playbook-ecr.yaml