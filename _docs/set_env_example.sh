#!/usr/bin/env bash

export TCTEST_SERVER="ci.katbyte.me"
export TCTEST_BUILDTYPEID="TCBuildTypeId"
export TCTEST_USER="katbyte"
export TCTEST_PASS="nope"

export TCTEST_REPO="terraform-providers/terraform-provider-azurerm"

#handy way to set password without it ending up in your bash history
ZZ_tcpass() {
    echo -n "Set \$TCTEST_PASS to: "
    read -s INPUT
    echo

    export TCTEST_PASS="$INPUT"
}
alias tcpass=ZZ_tcpass