#!/bin/bash
kubectl create ns mango-admin
kubectl create configmap settings-admin --from-file=../../config/settings.yml -n mango-admin
