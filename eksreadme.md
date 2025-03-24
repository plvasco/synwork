all these commands

   ### Eksctl get clusters list
    eksctl get cluster --profile=<profile-name> --region=us-east-1

    # Get addons information
    eksctl get addons --cluster=<cluster-name> --profile=<profile-name> --region=us-east-1

    ### Check the worker node version:
    kubectl get nodes -o wide
    
    ### Get nodegroup information
    eksctl get nodegroup --cluster=<cluster-name> --profile=<profile-name> --region us-east-1

    ## Info about any additional installation.
    kubectl get deployment -n kube-system


