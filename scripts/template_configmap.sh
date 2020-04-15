#!/bin/bash

# replace the pubsub project, topic, and subscription with their respective substitutions

project=$1
topic=$2
subscription=$3

template=`cat "kubernetes-manifests/bot/configmap.yaml" | 
sed "s/{{_PUBSUB_PROJECT}}/$project/g" |
sed "s/{{_PUBSUB_TOPIC}}/$topic/g" |
sed "s/{{_PUBSUB_SUBSCRIPTION}}/$subscription/g"`

echo "$template" > kubernetes-manifests/bot/configmap.yaml