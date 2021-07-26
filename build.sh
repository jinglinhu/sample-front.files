#!/bin/bash

TAG=1.0.4

REPOSITORY=jinglinhu/eks-workshop-x-ray-sample-front

docker build --tag $REPOSITORY:$TAG .

docker push $REPOSITORY:$TAG
