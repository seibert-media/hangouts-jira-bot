#!/bin/bash

# replace the pubsub project, topic, and subscription with their respective substitutions

template=`cat "kubernetes-manifests/bot/configmap.yaml" | 
sed "s/{{_PUBSUB_PROJECT}}/$_PUBSUB_PROJECT/g" |
sed "s/{{_PUBSUB_TOPIC}}/$_PUBSUB_TOPIC/g" |
sed "s/{{_PUBSUB_SUBSCRIPTION}}/$_PUBSUB_SUBSCRIPTION/g"`

echo "$template" > kubernetes-manifests/bot/configmap.yaml