#!/bin/env -S bash -l

set -eu

is_codeclimate_setup() {
    set +u
    if [ -z "${CC_TEST_REPORTER_ID}" ]; then
        set -u
        return 1
    fi
    set -u
    return 0
}

codeclimate_precheck() {
    # This function requires CC_TEST_REPORTER_ID env var, containing
    # CodeClimate Token, to be present in the environment
    if is_codeclimate_setup; then
        echo "*** CodeClimate is setup"
        CC_TEST_REPORTER_URL=https://codeclimate.com/downloads/test-reporter/
        CC_TEST_REPORTER_BIN=test-reporter-latest-linux-amd64
        curl -L ${CC_TEST_REPORTER_URL}/${CC_TEST_REPORTER_BIN} > ./cc-test-reporter
        chmod +x ./cc-test-reporter
        ./cc-test-reporter before-build
    else
        echo "*** CodeClimate is not setup"
    fi
}

convert_coverage_report_gcov2lcov() {
    go install github.com/jandelgado/gcov2lcov@latest
    gcov2lcov -infile=${1} -outfile=${2}
}

tests_check() {
    echo "*** Run unit tests"
    EXTRA_TEST_PARAMS=""
    if is_codeclimate_setup; then
        mkdir -p coverage
        EXTRA_TEST_PARAMS="${EXTRA_TEST_PARAMS}-coverprofile=coverage/gcov.out "
    fi
    go test -v ./... -race -covermode=atomic ${EXTRA_TEST_PARAMS}
    if is_codeclimate_setup; then
        convert_coverage_report_gcov2lcov coverage/gcov.out coverage/lcov.info
        ./cc-test-reporter after-build -t lcov
    fi
}

complexity_check() {
    echo "*** Complexity check"
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
    gocyclo -avg -top 5 -over 15 .
}

format_check() {
    echo "*** Format check"
    gofmt -s -e -d -l . | tee /tmp/gofmt.output && [ $(cat /tmp/gofmt.output | wc -l) -eq 0 ]
}

inefficiencies_check() {
    echo "*** Inefficiencies check"
    go install github.com/gordonklaus/ineffassign@latest
    go mod tidy
    ineffassign ./...
}

smells_check() {
    echo "*** Smells check"
    go mod tidy
    go vet ./...
}

spelling_check() {
    echo "*** Spelling check"
    go install github.com/client9/misspell/cmd/misspell@latest
    misspell -error .
}

static_check() {
    echo "*** Static check"
    go install honnef.co/go/tools/cmd/staticcheck@latest
    go mod download
    staticcheck ./...
}

style_check() {
    echo "*** Style check"
    go install golang.org/x/lint/golint@latest
    golint ./...
}

find_functions() {
    HOOK=${1:-check}

    declare -F | awk '{print $3}' | grep -E "_${HOOK}$" | sed -e 's/_${HOOK}$//'
}

try() {
    if find_functions "check" | grep -w ${1} >/dev/null; then
	${1} && echo "=== OK!" || (echo "=== NOK!" && return -1)
    else
        echo No ${1} available
	return 255
    fi
}

if [ "${1}" == "all" ]; then
    failure=0
    set +e
    codeclimate_precheck
    for check in $(find_functions "check"); do
        try ${check} || failure=1
    done
    set -e
    exit ${failure}
else
    try ${1}
fi
