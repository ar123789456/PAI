package main

import (
	"fmt"
	"os"
)

var A *box

func main() {
	for true {
		var Maps maps
		Maps = info(Maps)
		// if A != nil {
		// 	Maps.Maps[A.y][A.x].prise += -20
		// }
		// Maps.Maps[0][0].name = "q"
		Maps.mapsAddNeighbors()
		// number of entities
		allMobs := mobsInfo(&Maps)
		Maps.transvaluation(allMobs.me)

		a := bfs(&Maps, allMobs)
		// for _, l := range a {
		// 	fmt.Println(l.sumTotal)
		// }
		// fmt.Println(a)
		var path *box
		for _, k := range a {
			if path == nil {
				path = k
				continue
			}
			// fmt.Println(k.sumTotal)
			if k.sumTotal >= path.sumTotal && k.lenth < path.lenth {
				path = k
			}
		}
		// fmt.Println(path.sumTotal)
		// use `os.Stderr` to print for debugging
		fmt.Fprintf(os.Stderr, "debug code\n")

		PrintResult(path, &allMobs)
	}
}

func bfs(Maps *maps, allMobs mobs) []*box {
	graphs := []*box{}
	open := []*box{}
	open = append(open, Maps.Maps[allMobs.me.y][allMobs.me.x])
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if now.touch && len(open) != 0 {
			continue
		}
		for _, neig := range now.neighbors {
			if neig.touch {
				// if len(open) == 0 {
				// 	if neig.parent != nil && neig.lenth != 0 {
				// 		continue
				// 	}
				// } else {
				continue
				// }
			}
			// if neig.prise < -40 {
			// 	continue
			// }
			neig.sumTotal = now.sumTotal + neig.prise
			neig.parent = now
			neig.lenth = now.lenth + 1
			if neig.name != "." {
				graphs = append(graphs, neig)
			}
			open = append(open, neig)
		}
		now.touch = true
	}
	return graphs
}

func info(Maps maps) maps {
	var w, h, playerID, tick int
	fmt.Scan(&w, &h, &playerID, &tick)
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%v %v %v %v\n", w, h, playerID, tick))
	Maps.mapsInit(w, h, playerID, tick)
	// read map
	for i := 0; i < h; i++ {
		line := ""
		var lineBox []*box
		fmt.Scan(&line)
		fmt.Fprint(os.Stderr, line, "\n")
		for x, j := range line {
			var bOx box
			bOx.addBox(string(j), x, i)
			lineBox = append(lineBox, &bOx)
		}
		Maps.Maps = append(Maps.Maps, lineBox)
	}
	return Maps
}

func mobsInfo(Maps *maps) mobs {
	var n int
	fmt.Scan(&n)
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", n))

	var allmobs mobs
	var me mob
	var opponent mob
	monster := []*mob{}
	// read entities
	for i := 0; i < n; i++ {
		var entType string
		var pID, x, y, param1, param2 int
		fmt.Scan(&entType, &pID, &x, &y, &param1, &param2)
		fmt.Fprintf(os.Stderr, entType)
		fmt.Fprintf(os.Stderr, fmt.Sprintf(" %v %v %v %v %v\n", pID, x, y, param1, param2))
		if pID == Maps.playerID {
			me.mobInit(entType, pID, x, y, param1, param2)
		} else if pID != 0 {
			opponent.mobInit(entType, pID, x, y, param1, param2)
		} else {
			var m mob
			m.mobInit(entType, pID, x, y, param1, param2)
			monster = append(monster, &m)
			Maps.Maps[y][x].name = entType
		}
	}
	allmobs.mobsInit(me, opponent, monster)
	return allmobs
}

func PrintResult(path *box, mObs *mobs) {
	if path == nil {

		fmt.Println("stay")
		fmt.Fprintf(os.Stderr, "stay\n")
		return
	}

	var finalx, finaly int
	for path.parent != nil {
		finalx = mObs.me.y - path.y
		finaly = mObs.me.x - path.x
		path = path.parent
	}
	A = path
	// fmt.Println(finalx, finaly)
	if finalx == 1 {
		fmt.Println("up")
		fmt.Fprintf(os.Stderr, "up\n")
		return
	}
	if finalx == -1 {
		fmt.Println("down")
		fmt.Fprintf(os.Stderr, "down\n")
		return
	}
	if finaly == -1 {
		fmt.Println("right")
		fmt.Fprintf(os.Stderr, "right\n")
		return
	}
	if finaly == 1 {
		fmt.Println("left")
		fmt.Fprintf(os.Stderr, "left\n")
		return
	}
	fmt.Println("stay")
	fmt.Fprintf(os.Stderr, "stay\n")
}

type mobs struct {
	me       *mob
	opponent *mob
	monster  []*mob
}

func (self *mobs) mobsInit(me, opponent mob, monster []*mob) {
	self.me = &me
	self.opponent = &opponent
	self.monster = monster
}

type mob struct {
	name   string
	pID    int
	x      int
	y      int
	dagger int
	bonus  int
}

func (self *mob) mobInit(entType string, pID, x, y, param1, param2 int) {
	self.name = entType
	self.pID = pID
	self.x = x
	self.y = y
	self.dagger = param1
	self.bonus = param2
}

type maps struct {
	Maps     [][]*box
	w        int
	h        int
	playerID int
	tick     int
}

func (self *maps) mapsInit(w, h, playerID, tick int) {
	self.w = w
	self.h = h
	self.playerID = playerID
	self.tick = tick
}

func (self *maps) mapsAddNeighbors() {
	for y, i := range self.Maps {
		for x, j := range i {
			if j.name == "!" {
				continue
			}
			y1 := y - 1
			x1 := x - 1
			y2 := y + 1
			x2 := x + 1
			if x1 >= 0 {
				if self.Maps[y][x1].name != "!" {
					j.neighbors = append(j.neighbors, self.Maps[y][x1])
				}
			}
			if y1 >= 0 {
				if self.Maps[y1][x].name != "!" {
					j.neighbors = append(j.neighbors, self.Maps[y1][x])
				}
			}
			if x2 < self.w {
				if self.Maps[y][x2].name != "!" {
					j.neighbors = append(j.neighbors, self.Maps[y][x2])
				}
			}
			if y2 < self.h {
				if self.Maps[y2][x].name != "!" {
					j.neighbors = append(j.neighbors, self.Maps[y2][x])
				}
			}
			self.Maps[y][x] = j
		}
	}
}

func (self *maps) transvaluation(mOb *mob) {
	for _, i := range self.Maps {
		for _, j := range i {
			if j.name == "." {
				continue
			}
			if j.name == "m" {
				if mOb.dagger != 0 {
					// aura(j, 10)s
					continue
				}
				aura(j, -20)
			} else if j.name == "#" {
				j.prise += 20
			} else {
				if j.name != "q" {
					j.prise += 15
				}
			}
		}

	}
}

func aura(point *box, prise int) {
	open := []*box{}
	point.prise += prise * 5
	point.touch = true
	open = append(open, point)
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if now.radius < 4 {
			now.touch = true
		}
		if now.radius > 4 {
			continue
		}
		now.prise += prise * (5 - now.radius)
		for _, neig := range now.neighbors {
			neig.radius = now.radius + 1
			open = append(open, neig)
		}
	}
}

type box struct {
	name      string
	touch     bool
	prise     int
	radius    int
	sumTotal  int
	lenth     int
	x         int
	y         int
	parent    *box
	neighbors []*box
}

func (self *box) addBox(s string, x, y int) {
	self.name = s
	self.x = x
	self.y = y
}
