#! /bin/sh


# << snippet begin >>
# No auto fix: replace math.Min/math.Max with min/max
grep -r 'math.M\(in\|ax\)' . | sed 's/$/ # replace math\.Min\/math\.Max with min\/max/'
# << snippet end >>
