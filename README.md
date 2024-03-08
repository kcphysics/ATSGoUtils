# ATSGoUtils
Utilities for command line written in Go for ATS Trekmush!

## Best Route Algorithm

The Best Route Algorithm goes as follows:

1. Find the Direct Route (Unless either the Source or Target resides in DQ and the other does not)
2. Calculate First Legs (Source to Gates), find minimum by distance
3. Calculate Second Legs (Target to Gates), find minimum by distance

IF direct route is possible, compare against the sum of the distances of the shortest legs
IF direct route is not possible, return the shortest legs