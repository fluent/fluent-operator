PROJECT_ROOT=$PWD
E2E_DIR=$(realpath $(dirname $0)/..)

function build_ginkgo_test() {
  cd $E2E_DIR
  ginkgo build -r e2e/fluentd/
}

function cleanup() {
  cd $PROJECT_ROOT
  kubectl delete -f manifests/setup/setup.yaml
  kubectl delete ns kubesphere-logging-system
  kind delete cluster --name test
}

function prepare_cluster() {
  kind create cluster --name test 
  kubectl create ns kubesphere-logging-system

  echo "wait the control-plane ready..."
  kubectl wait --for=condition=Ready node/test-control-plane --timeout=60s

  # kubectl create clusterrolebinding system:anonymous --clusterrole=cluster-admin --user=system:anonymous
}

function start_fluent_operator() {
  export KUBECONFIG=$HOME/.kube/config

  cd $PROJECT_ROOT && kubectl apply -f manifests/setup/setup.yaml

  while true; do
      sleep 3
      kubectl get po -nkubesphere-logging-system 2>/dev/null | grep -q fluent-operator && break
  done
}

function run_test() {
  # inspired by github.com/kubeedge/kubeedge/tests/e2e/scripts/helm_keadm_e2e.sh
  :> /tmp/testcase.log
  $E2E_DIR/e2e/fluentd/fluentd.test $debugflag 2>&1 | tee -a /tmp/testcase.log
  
  #stop the edgecore after the test completion
  grep  -e "Running Suite" -e "SUCCESS\!" -e "FAIL\!" /tmp/testcase.log | sed -r 's/\x1B\[([0-9];)?([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]//g' | sed -r 's/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]//g'
  echo "Integration Test Final Summary Report"
  echo "======================================================="
  echo "Total Number of Test cases = `grep "Ran " /tmp/testcase.log | awk '{sum+=$2} END {print sum}'`"
  passed=`grep -e "SUCCESS\!" -e "FAIL\!" /tmp/testcase.log | awk '{print $3}' | sed -r "s/\x1B\[([0-9];)?([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]//g" | awk '{sum+=$1} END {print sum}'`
  echo "Number of Test cases PASSED = $passed"
  fail=`grep -e "SUCCESS\!" -e "FAIL\!" /tmp/testcase.log | awk '{print $6}' | sed -r "s/\x1B\[([0-9]{1,2}(;[0-9]{1,2})?)?[mGK]//g" | awk '{sum+=$1} END {print sum}'`
  echo "Number of Test cases FAILED = $fail"
  echo "==================Result Summary======================="

  if [ "$fail" != "0" ];then
      echo "Integration suite has failures, Please check !!"
      exit 1
  else
      echo "Integration suite successfully passed all the tests !!"
      exit 0
  fi
}

set -Ee
trap cleanup EXIT
trap cleanup ERR

echo -e "\nBuilding testcases..."
build_ginkgo_test

export KUBECONFIG=$HOME/.kube/config

echo -e "\nPreparing cluster..."
prepare_cluster

echo -e "\nStart fluent operator..."
start_fluent_operator

echo -e "\nRunning test..."
run_test
