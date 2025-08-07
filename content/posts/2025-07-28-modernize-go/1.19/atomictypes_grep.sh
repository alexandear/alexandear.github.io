#! /bin/sh

# << snippet begin >>
# No auto fix: use atomic types instead of atomic operations
grep -rn 'atomic\.\(Store\|Load\|CompareAndSwap\)[A-Za-z0-9]*' . | sed 's/$/ # use atomic types instead of atomic operations/'
# << snippet end >>
