#!/usr/bin/env bash
echo
echo "Start build SessionService"
echo

# Execute tests
go test -v ./
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
~/.go/bin/gometalinter
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