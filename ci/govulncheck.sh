#!/usr/bin/env bash

# Copyright (c) 2026 Proton AG
#
# This file is part of Proton Mail Bridge.
#
# Proton Mail Bridge is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Proton Mail Bridge is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Proton Mail Bridge.  If not, see <https://www.gnu.org/licenses/>.

set -eo pipefail

# Absolute: the script cd's into go-lib before reading it.
CONFIG="$PWD/.proton/reviewbot/config.yaml"

echo "Using Go version:"
go version
echo

cd ./go-lib

GOTOOLCHAIN=auto go run golang.org/x/vuln/cmd/govulncheck@latest -json ./... > vulns.json

jq -r '.finding | select((.osv != null) and (.trace[0].function != null)) | .osv' < vulns.json > vulns_osv_ids.txt

# Read GO- prefixed ignore rules from reviewbot config (single source of truth)
if [ -f "$CONFIG" ]; then
    grep -oE 'GO-[0-9]{4}-[0-9]+' "$CONFIG" | while read -r id; do
        echo "ignoring $id (tracked in $CONFIG)"
        grep -v "$id" < vulns_osv_ids.txt > tmp || true
        mv tmp vulns_osv_ids.txt
    done
fi

# Fail if any unignored vulns remain
if [ -s vulns_osv_ids.txt ]; then
    while read -r osv; do
        jq --arg osvid "$osv" \
            '.osv | select(.id == $osvid) | {"id":.id, "ranges": .affected[0].ranges, "import": .affected[0].ecosystem_specific.imports[0].path}' \
            < vulns.json
    done < vulns_osv_ids.txt
    echo
    echo "Vulnerability found"
    exit 1
fi

echo
echo "No new vulnerabilities found."
