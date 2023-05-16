#!/bin/bash

_rainbowcal() {
	unbuffer cal --color | lolcat -f -F '0.5' | ansifilter --html | tail -n +20 | head -n -3
}
_web3help() {
	unbuffer go run web3.go --help --help | ansifilter --html | tail -n +20 | head -n -3
}
_dayscalc() {
	printf 'There are %s days in the month of %s.\n' "$(cal "$(date +%m)" "$(date +%Y)" | awk 'NF {DAYS = $NF}; END {print DAYS}')" "$(date +%B)"
	printf 'Today is %s %s.\n' "$(date +%B)"  "$(date +%d)"
	printf '%s days remain in the month of %s.\n' "$(echo "$(cal "$(date +%m)" "$(date +%Y)" | awk 'NF {DAYS = $NF}; END {print DAYS}')" - "$(date +%d)" | bc -l)" "$(date +%B)"
	printf '%s days in the year %s.\n' "$(seq  01 12 | while read -r _month ; do cal "$_month" "$(date +%Y)" | awk 'NF {DAYS = $NF}; END {print DAYS}' ; done | paste -sd+ | bc -l)" "$(date +%Y)"
	printf 'Today is day %s.\n' "$(date +%j)"
	printf 'There are %s  days remaining in %s\n'  "$(echo "$(echo "$(date --date="January 1 next year" +%s)" - "$(date +%s)"  | bc -l)" / 60 / 60 / 24 | bc)" "$(date +%Y)"
}
