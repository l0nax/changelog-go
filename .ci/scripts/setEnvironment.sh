#!/bin/bash -
#===============================================================================
#
#          FILE: setEnvironment.sh
#
#         USAGE: ./setEnvironment.sh
#
#   DESCRIPTION: This script checks the current environment and finalizes/
#                deploys the sentry release with the correct env.
#
#       OPTIONS: ---
#  REQUIREMENTS: sentry-cli
#          BUGS: ---
#         NOTES: ---
#        AUTHOR: Francesco Emanuel Bennici (l0nax), benniciemanuel78@gmail.com
#  ORGANIZATION: FABMation GmbH
#       CREATED: 03/12/2020 02:02:52 PM
#      REVISION: 001
#===============================================================================

set -o nounset                              # Treat unset variables as an error

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# we check if this release is an production release or an dev comment
# there are some cases where we can say that the current CI/CD process is for
# a MR/ normal commit on develop or a feature branch, etc...
#
# Available Environments
# ======================
#
# - production
#     Any deployed release (tagged).
#
# - staging
#     Builds which are merged on `develop` – so they are already reviewed and
#     checked.
#
# - develop (default)
#     All changes which are done in (ie.) feature branches, this separate
#     environment is required to be able to filter them out in the Sentry
#     Dashboard. Because while development the application could crash many
#     times but this bug COULD not be committed – so we will never see the bad
#     code, and also Sentry can not recognize the buggy code – and the
#     developer will resolve it ASAP. So we are not interested in such error
#     reports. And the `changelog-go` application will also NOT send crashes to
#     the server if this environment is set.
if [[ "$VERSION" == *"-dev"* ]]; then
        if [ "$CI_COMMIT_BRANCH" == "develop" ]; then
          export ENV="staging"
        else
          export ENV="develop"
        fi
else
        export ENV="production"
fi

## write environment to file
echo ${ENV} > ${DIR}/../../ENVIRONMENT
