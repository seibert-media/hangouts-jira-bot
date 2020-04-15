#!/bin/bash

# replace the pubsub project, topic, and subscription with their respective substitutions

file=kubernetes-manifests/bot/configmap.yaml

project=$PUBSUB_PROJECT
topic=$PUBSUB_TOPIC
subscription=$PUBSUB_SUBSCRIPTION

echo "templating $file $project $topic $subscription"

template=`cat "$file" | 
sed "s/{{_PUBSUB_PROJECT}}/$project/g" |
sed "s/{{_PUBSUB_TOPIC}}/$topic/g" |
sed "s/{{_PUBSUB_SUBSCRIPTION}}/$subscription/g"`

echo "$template"
echo "$template" > $file