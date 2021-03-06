#!/usr/bin/env bash
echo
echo -en "\E[40;36m\033[1mStart build SessionService\033[0m"
echo
if [ ! -z "$1" ]
then
	export REDISADDR=$1
	echo "Redis Host: " $REDISADDR
fi
echo

# Execute tests
echo -en "\E[40;35m\033[1mExecute tests\033[0m"
echo
go test -v ./...
EXITCODE=$?
if [ $EXITCODE != 0 ]
then
	echo
	echo -en "\E[40;31m\033[1mTests failed with exit code: \033[0m" $EXITCODE
	echo
	echo
	exit $EXITCODE
fi

# Execute linter
echo
echo -en "\E[40;35m\033[1mExecute linters\033[0m"
echo
gometalinter --config=gometalinter.json
EXITCODE=$?
if [ $EXITCODE != 0 ]
then
	echo
	echo -en "\E[40;31m\033[1mLinting failed with exit code: \033[0m" $EXITCODE
	echo
	echo
	exit $EXITCODE
fi

# Success
echo
echo -en "\E[40;32m\033[1mBuild successful :-)\033[0m"
echo
echo