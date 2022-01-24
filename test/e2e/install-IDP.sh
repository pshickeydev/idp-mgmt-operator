#!/bin/bash

# Copyright Red Hat

#set -e

echo "Install identity configuration management service for Kubernetes ..."

idp_dir=$(mktemp -d -t idp-XXXXX)
cd "$idp_dir" || exit 1
export HOME="$idp_dir"

# Set up repo URLs.
echo "Cloning branch idp-mgmt-operator ${PULL_BASE_REF}"
idp_mgmt_operator_url="https://${IDP_MGMT_OPERATOR_REPO}.git"
idp_mgmt_operator_dir="${idp_dir}/idp-mgmt-operator"
git clone -b "${PULL_BASE_REF}" "$idp_mgmt_operator_url" "$idp_mgmt_operator_dir" || {
    echo "ERROR Could not clone branch ${PULL_BAWE_REF} from idp-mgmt-operator repo $idp_mgmt_operator_url"
    exit 1
}


cd ${idp_mgmt_operator_dir}


#echo REPO_OWNER ${REPO_OWNER}
#echo REPO_NAME ${REPO_NAME}
#echo PULL_BASE_REF ${PULL_BASE_REF}
#echo PULL_BASE_SHA ${PULL_BASE_SHA}
#echo PULL_NUMBER ${PULL_NUMBER}
#echo PULL_PULL_SHA ${PULL_PULL_SHA}

export IMG="quay.io/identitatem/idp-mgmt-operator:0.1-PR${PULL_NUMBER}-${PULL_PULL_SHA}"
echo "Quay image is ${IMG}"

echo "Check namespace - before"
oc get namespaces

#oc create ns idp-mgmt-config

#echo "Check namespace - before 2"
#oc get namespaces

echo "Start deploy"
#export CATALOG_DEPLOY_NAMESPACE=idp-mgmt-config
make deploy

sleep 20

echo "Check namespace - after"
oc get namespaces

echo "Check that the pod is running"
oc wait --for=condition=ready pods --all --timeout=5m -n idp-mgmt-config

oc get pods -n idp-mgmt-config


echo "Create IDPConfig"
cat > e2e-IDPConfig.yaml <<EOF
apiVersion: identityconfig.identitatem.io/v1alpha1
kind: IDPConfig
metadata:
  name: idp-config
  namespace: idp-mgmt-config
spec:
EOF
oc create -f e2e-IDPConfig.yaml

sleep 20

echo "Check for more pods running"
oc wait --for=condition=ready pods --all --timeout=5m -n idp-mgmt-config
oc get pods -n idp-mgmt-config

echo "Done installing identity configuration management service for Kubernetes"