PROJECT_ROOT=$PWD
E2E_DIR=$(realpath $(dirname $0)/..)
LOGGING_NAMESPACE=fluent
IMAGE_TAG=`date "+%Y-%m-%d-%H-%M-%S"`

function build_ginkgo_test() {
  cd $E2E_DIR
  ginkgo build -r e2e/fluentd/
}

function cleanup() {
  cd $PROJECT_ROOT
#  helm uninstall fluent-operator -n $LOGGING_NAMESPACE
#  kubectl delete ns $LOGGING_NAMESPACE
  kind delete cluster --name test-helm && exit 0
}

function prepare_cluster() {
  kind create cluster --name test-helm
  kubectl create ns $LOGGING_NAMESPACE

  echo "wait the control-plane ready..."
  kubectl wait --for=condition=Ready node/test-helm-control-plane --timeout=60s
}

function build_image() {
  cd $PROJECT_ROOT
  make build-op-amd64 -e FO_IMG=kubesphere/fluent-operator:$IMAGE_TAG
  kind load docker-image kubesphere/fluent-operator:$IMAGE_TAG --name test-helm
}

function start_fluent_operator() {
  cd $PROJECT_ROOT && helm install --wait --timeout 30s fluent-operator  --create-namespace -n $LOGGING_NAMESPACE charts/fluent-operator/ --set operator.container.tag=$IMAGE_TAG
  kubectl -n $LOGGING_NAMESPACE wait --for=condition=available deployment/fluent-operator --timeout=60s
}

function run_test() {
  # inspired by github.com/kubeedge/kubeedge/tests/e2e/scripts/helm_keadm_e2e.sh
  :> /tmp/testcase.log
  $E2E_DIR/e2e/fluentd/fluentd.test $debugflag 2>&1 | tee -a /tmp/testcase.log

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

echo -e "\nPreparing cluster..."
prepare_cluster

echo -e "\nBuilding image..."
build_image

echo -e "\nStart fluent operator..."
start_fluent_operator

echo -e "\nRunning test..."
run_test