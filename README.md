# miyuki
Miyuki bot on GKE

# Setup
1. download gcloud command
2. add a hubot configuration to you slack

```
# Setup a project configuration
$ gcloud auth login
$ gcloud components update kubectl
$ gcloud config set project $PROJECT_ID
$ gcloud config set compute/region us-west1
$ gcloud config set compute/zone us-west1-b
$ gcloud container clusters create miyuki-cluster \
      --machine-type f1-micro \
      --disk-size=30 \
      --num-nodes=3
$ kubectl get nodes
$ gcloud container clusters resize miyuki-cluster --size=1
$ gcloud container clusters get-credentials miyuki-cluster
$ gcloud container clusters describe miyuki-cluster

# CloudBuild
$ gcloud container builds submit --config cloudbuild.yaml

# Deployment
$ kubectl run pod-miyuki \
      --image=gcr.io/$PROJECT_ID/miyuki:latest \
      --env="HUBOT_SLACK_TOKEN=$MIYUKI_HUBOT_SLACK_TOKEN" \
      --env="HUBOT_SLACK_TEAM=$MIYUKI_HUBOT_SLACK_TEAM" \
      --env="HUBOT_SLACK_BOTNAME=miyuki" \
      --port=8080 \
      --restart='Always'
$ kubectl get deployments
$ kubectl get pods
$ kubectl exec -it $MIYUKI_POD_NAME /bin/bash
$ kubectl get service
$ kubectl expose deployment $MIYUKI_POD_NAME --type="LoadBalancer"
$ kubectl get service
```