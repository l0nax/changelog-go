#!/bin/bash -
#===============================================================================
#
#          FILE: getReleases.sh
#
#         USAGE: ./getReleases.sh <VERSION>
#
#   DESCRIPTION: Creates the 'public' dir with 'go-selfupdate'
#
#       OPTIONS: ---
#  REQUIREMENTS: go-selfupdate, sed, tr, cut, echo
#          BUGS: ---
#         NOTES: ---
#        AUTHOR: Francesco Emanuel Bennici (l0nax), benniciemanuel78@gmail.com
#  ORGANIZATION: FABMation GmbH
#       CREATED: 03/11/2020 06:38:38 PM
#      REVISION:  ---
#===============================================================================

set -o nounset                              # Treat unset variables as an error

for i in $(ls -d dist/*/* | grep -v prime); do
  dist=$(echo ${i} | cut -d '/' -f2 | sed -e 's/default_//g' | tr '_' '-')

  go-selfupdate -platform ${dist} ${i} ${1}
done

