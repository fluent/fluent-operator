set -o errexit
set -o nounset
set -o pipefail

FLUENT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
FLUENT_HELM_DIR=${FLUENT_ROOT}/charts/fluent-operator
_tmpdir=/tmp/fluent-operator

cd $FLUENT_HELM_DIR && helm package . -d ${_tmpdir} > /dev/null && mv ${_tmpdir}/*.tgz $FLUENT_HELM_DIR/../fluent-operator.tgz && rm -rf "${_tmpdir}"
