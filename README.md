# khomer
kubernetes self register homer like project

Running this example

First compile the application for Linux:

    GOOS=linux go build -o ./app .

Then package it to a docker image using the provided Dockerfile to run it on Kubernetes.

Build the docker image:

    docker build -t danielcosta1643360/in-cluster:2 .
    docker push danielcosta1643360/in-cluster:2

If you have RBAC enabled on your cluster, use the following snippet to create role binding which will grant the default service account view permissions.

    kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

Then, run the image in a Pod with a single instance Deployment:

    kubectl run --rm -i demo --image=danielcosta1643360/in-cluster:2

Port fwd to test:

    kubectl port-forward demo 1323:1323