## Totemo Solver

This is a brute-force solver, written because one of the puzzles in
Totemo was challenging enough to make me want to write it :-)

To use it, you'll need Go (http://www.golang.org) or to port it to
a different language. 

Then, build with [gb](http://github.com/skelterjohn/go-gb), run it with
a totem height set and pipe one of the puzzles to STDIN:

	git clone http://github.com/fluffle/totemo
	cd totemo
	gb .
	# 4 signifies the totem height
	./totemo 4 < hard.to

The "hard" puzzle takes 14.5s to find a solution on my crappy 1.4GHz
Core2 U9400. I'm sure this could be optimised.
