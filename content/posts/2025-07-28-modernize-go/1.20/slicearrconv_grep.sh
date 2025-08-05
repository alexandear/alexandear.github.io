#! /bin/sh

# << snippet begin >>
# No auto fix: replace *(*[N]T)(slice) with [N]T(slice)
grep -rn '\*(\*\[[0-9]\+\][^)]*)(.*)' | sed 's/$/ # replace \*\(\*\[N\]T\)\(slice\) with \[N\]\T\(slice\)/'
# << snippet end >>
