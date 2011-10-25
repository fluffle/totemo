package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
)

var totem = flag.String("totem", "", "load totem from this file rather than stdin")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// A totemo level is an 8x6 grid with points having values
type grid [6][8]int

// A totemo move consists of two or more points
type p struct {
	r, c int
}
type move []p

// And we need a stack to push moves (and their associated grids) onto.
type stack struct {
	g *grid
	m []move
}

func (g *grid) output(w io.Writer) {
	for i := 0; i < 6; i++ {
		io.WriteString(w, fmt.Sprintf("%d %d %d %d %d %d %d %d\n", 
			g[i][0], g[i][1], g[i][2], g[i][3],
			g[i][4], g[i][5], g[i][6], g[i][7]))
	}
}

func (g *grid) load(r io.Reader) {
	for i := 0; i < 6; i++ {
		fmt.Fscanf(r, "%d %d %d %d %d %d %d %d",
			&g[i][0], &g[i][1], &g[i][2], &g[i][3],
			&g[i][4], &g[i][5], &g[i][6], &g[i][7])
	}
}

func (g *grid) empty() bool {
	for r := 0; r < 6; r++ {
		for c := 0; c < 8; c++ {
			if g[r][c] > 0 {
				return false
			}
		}
	}
	return true
}

func (g *grid) move(m move) *grid {
	n := new(grid)
	for r := 0; r < 6; r++ {
		for c := 0; c < 8; c++ {
			n[r][c] = g[r][c]
		}
	}
	for _, p := range m {
		n[p.r][p.c] = 0
	}
	return n
}

func (g *grid) possible(sz int) []move {
	s := make([]move, 0)
	for r := 0; r < 6; r++ {
		for c := 0; c < 8; c++ {
			if g[r][c] > 0 {
				if m := g.checkrow(sz, r, c); len(m) > 0 {
					s = append(s, m)
				}
				if m := g.checkcol(sz, r, c); len(m) > 0 {
					s = append(s, m)
				}
			}
		}
	}
	return s
}

func (g *grid) checkrow(sz, r, c int) move {
	t := g[r][c]
	m := move{p{r,c}}
	if t == sz {
		// this point on it's own satisfies a move of size sz
		return m
	}
	// check right along row
	for c2 := c+1; c2 < 8; c2++ {
		if g[r][c2] == 0 { continue }
		t += g[r][c2]
		if t > sz { return move{} }
		m = append(m, p{r,c2})
		if t == sz { return m }
	}
	// if we get here no moves are possible
	return move{}
}

func (g *grid) checkcol(sz, r, c int) move {
	t := g[r][c]
	m := move{p{r,c}}
	if t == sz {
		// We only want to return a single-point move once, so
		// do it in checkrow but not here in checkcol
		return move{}
	}
	// Check down along column
	for r2 := r+1; r2 < 6; r2++ {
		if g[r2][c] == 0 { continue }
		t += g[r2][c]
		if t > sz { return move{} }
		m = append(m, p{r2,c})
		if t == sz { return m }
	}
	// if we get here no moves are possible
	return move{}
}

func (m move) output(w io.Writer) {
	s := "move: "
	for _, p := range m {
		s += fmt.Sprintf("[%d,%d] ", p.r, p.c)
	}
	s += "\n"
	io.WriteString(w, s)
}

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f, err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal(err)
        }
		fmt.Println("Profiling CPU")
        pprof.StartCPUProfile(f)
        defer pprof.StopCPUProfile()
    }

	g := new(grid)
	if *totem != "" {
		f, err := os.Open(*totem)
		if err != nil {
			log.Fatal(err)
		}
		g.load(f)
		f.Close()
	} else {
		g.load(os.Stdin)
	}
	st := []stack{{g, make([]move, 0)}}
	sz := 2
	if len(flag.Args()) > 0 {
		sz, _ = strconv.Atoi(flag.Arg(0))
	}
	if moves := search(st, sz); len(moves) > 0 {
		fmt.Printf("Found a solution in %d moves.\n", len(moves))
		for i, m := range moves {
			fmt.Printf("Move %d:\n", i)
			g = g.move(m)
			g.output(os.Stdout)
			m.output(os.Stdout)
			fmt.Println()
		}
	} else {
		fmt.Printf("No solution found :-(\n")
	}
}

func search(st []stack, sz int) []move {
	var s stack
	for len(st) > 0 {
		// pull a stack frame from the end of the stack
		// (we want a depth-first search for memory reasons)
		st, s = st[:len(st)-1], st[len(st)-1]
		for _, m := range s.g.possible(sz) {
			// create a new grid copy taking this possible move
			g := s.g.move(m)
			// copy all the moves from the stack
			n := make([]move, len(s.m)+1)
			copy(n, s.m)
			// and append this one
			n[len(n)-1] = m
			if g.empty() {
				// fuck it, bail out when we have one solution :-)
				return n
			} else {
				// otherwise push the new grid and move set onto the stack
				st = append(st, stack{g, n})
			}
		}
	}
	return make([]move, 0)
}
