package main

import (
	"fmt"
	"os"
)

const monsterRadius = 3

var dagger int

var lastBox *Mob

func main() {
	for true {
		var baseInfo base

		baseInfo.initmaps()
		baseInfo.addneighbors()

		// number of entities
		var mobs Mobs
		mobs.initMobs(&baseInfo)
		baseInfo.mob = &mobs
		if mobs.me.param1 != 0 && dagger == 0 {
			dagger = 14
		}
		if dagger != 0 {
			dagger--
		}
		for _, v := range mobs.monster {
			if mobs.me.param1 == 0 {
				baseInfo = mobsAura(baseInfo, *v, monsterRadius)
			}
		}
		// for _, v := range baseInfo.maps {
		// 	for _, i := range v {
		// 		fmt.Print(i.monsaura, " ")
		// 	}

		// 	fmt.Println("")
		// }
		baseInfo.bfs()
		// use `os.Stderr` to print for debugging
		// fmt.Fprintf(os.Stderr, "debug code\n")

		baseInfo.PrintResult()
		lastBox = baseInfo.mob.me

		// this will choose one of random actions

	}
}

type base struct {
	maps                 [][]*box
	w, h, playerID, tick int
	mob                  *Mobs
	path                 *box
}

func (self *base) PrintResult() {
	var x, y int

	actions := []string{"left", "right", "up", "down", "stay"}
	// var m bool
	if self.path != nil {
		if self.path.parent != nil {
			for self.path.parent.parent != nil {
				// fmt.Println("---------------------------------------------", self.path.name)
				self.path = self.path.parent
			}
		} else {
			fmt.Println(actions[4])
			fmt.Fprintf(os.Stderr, actions[4])

			return
		}
	} else {
		fmt.Fprintf(os.Stderr, "search Monster\n")
		lastme := self.maps[lastBox.y][lastBox.x]
		if !lastme.monsaura && (len(lastme.neighbors) != 2) {
			self.path = lastme
		} else {
			// m = true
			self.bfsMonter()
			// lastme.site = true
			// fmt.Fprintf(os.Stderr, fmt.Sprintf("x %v, y %v\n", self.path.x, self.path.y))
			if self.path != nil {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("x %v, y %v\n", self.path.x, self.path.y))
				if self.path.distans > 3 {
					fmt.Fprintf(os.Stderr, fmt.Sprintf("self.path.distanse %v\n", self.path.distans))
					fmt.Println(actions[4])
					fmt.Fprintf(os.Stderr, actions[4])

					return
				}
				for self.path.parent != nil {
					self.path.site = true
					self.path = self.path.parent
				}
			}
			p := self.path
			for _, i := range self.maps[self.mob.me.y][self.mob.me.x].neighbors {
				if i.site {
					continue
				}
				self.path = i
				if i.monsaura {
					continue
				}
				if len(i.neighbors) == 1 {
					continue
				}
				break
			}
			if self.path == p && len(self.maps[self.mob.me.y][self.mob.me.x].neighbors) != 1 {
				fmt.Println(actions[4])
				fmt.Fprintf(os.Stderr, actions[4])
			}
		}
	}
	x = self.mob.me.x - self.path.x
	y = self.mob.me.y - self.path.y
	fmt.Fprintf(os.Stderr, fmt.Sprintf("x := %v, y := %v\n", x, y))
	if x < 0 {
		fmt.Println(actions[1])
		fmt.Fprintf(os.Stderr, actions[1])
		return
	}
	if x > 0 {
		fmt.Println(actions[0])
		fmt.Fprintf(os.Stderr, actions[0])
		return
	}
	if y < 0 {
		fmt.Println(actions[3])
		fmt.Fprintf(os.Stderr, actions[3])
		return
	}
	if y > 0 {
		fmt.Println(actions[2])
		fmt.Fprintf(os.Stderr, actions[2])
		return
	}
	fmt.Println(actions[4])
	fmt.Fprintf(os.Stderr, actions[4])
}

type Mobs struct {
	n       int
	me      *Mob
	enamy   *Mob
	monster []*Mob
}

func (self *Mobs) initMobs(baseinfo *base) {
	var n int
	fmt.Scan(&n)
	fmt.Fprintf(os.Stderr, fmt.Sprintf("%v", n))
	self.n = n
	// read entities
	for i := 0; i < n; i++ {
		var mob Mob
		fmt.Scan(&mob.name, &mob.pID, &mob.x, &mob.y, &mob.param1, &mob.param2)
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%v %v %v %v %v %v\n", mob.name, mob.pID, mob.x, mob.y, mob.param1, mob.param2))
		if mob.pID == 0 {
			self.monster = append(self.monster, &mob)
			baseinfo.maps[mob.y][mob.x].name = "m"
		} else if baseinfo.playerID == mob.pID {
			self.me = &mob
		} else {
			self.enamy = &mob
		}
	}
}

func mobsAura(baseinfo base, mob Mob, radius int) base {
	open := []*box{}
	open = append(open, baseinfo.maps[mob.y][mob.x])
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if now.radius == radius {
			continue
		}
		if now.monsaura {
			continue
		}
		now.monsaura = true
		for _, v := range now.neighbors {
			if v == baseinfo.maps[baseinfo.mob.me.y][baseinfo.mob.me.x] {
				v.monsaura = true
				continue
			}
			v.radius = now.radius + 1
			open = append(open, v)
		}
	}
	return baseinfo
}

type Mob struct {
	name                      string
	pID, x, y, param1, param2 int
}

func (self *base) initmaps() {
	mapsBool := false
	if len(self.maps) != 0 {
		mapsBool = true
	}
	var w, h, playerID, tick int
	// fmt.Fprintf(os.Stderr, fmt.Sprintf("%v %v %v %v\n", w, h, playerID, tick))
	fmt.Scan(&w, &h, &playerID, &tick)
	self.h = h
	self.w = w
	self.playerID = playerID
	self.tick = tick
	// read map
	for i := 0; i < h; i++ {
		line := ""
		var l []*box
		fmt.Scan(&line)
		// fmt.Fprint(os.Stderr, line, "\n")
		for x, name := range line {
			if mapsBool {
				self.maps[i][x].name = string(name)
			} else {
				var b box
				b.name = string(name)
				b.x = x
				b.y = i
				l = append(l, &b)
			}
		}
		if !mapsBool {
			self.maps = append(self.maps, l)
		}
	}
}

func (self *base) addneighbors() {
	for y, i := range self.maps {
		for x, v := range i {
			if len(v.neighbors) != 0 {
				return
			}
			x1 := x - 1
			y1 := y - 1
			x2 := x + 1
			y2 := y + 1
			if v.name == "!" {
				continue
			}
			if x1 >= 0 {
				if self.maps[y][x1].name != "!" {
					v.neighbors = append(v.neighbors, self.maps[y][x1])
				}
			}
			if y1 >= 0 {
				if self.maps[y1][x].name != "!" {
					v.neighbors = append(v.neighbors, self.maps[y1][x])
				}
			}
			if x2 < self.w {
				if self.maps[y][x2].name != "!" {
					v.neighbors = append(v.neighbors, self.maps[y][x2])
				}
			}
			if y2 < self.h {
				if self.maps[y2][x].name != "!" {
					v.neighbors = append(v.neighbors, self.maps[y2][x])
				}
			}
		}
	}
}

type box struct {
	name      string
	y         int
	x         int
	parent    *box
	neighbors []*box
	monsaura  bool
	radius    int
	touch     bool
	touchMons bool
	site      bool
	distans   int
	findDis   int
}

func (self *base) bfs() {
	var open []*box
	open = append(open, self.maps[self.mob.me.y][self.mob.me.x])
	m := 0
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		fmt.Fprintf(os.Stderr, now.name)
		if now.touch {
			continue
		}
		if now.monsaura {
			m++
			continue
		}
		// if len(now.neighbors) == 2 {
		// 	if now.neighbors[0].monsaura || now.neighbors[1].monsaura {
		// 		continue
		// 	}
		// }
		if now.name == "d" {
			fmt.Fprintf(os.Stderr, "find d\n")
			self.path = now
			return
		}
		if now.name == "b" {
			if m != 0 && now.findDis < 2 {
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				continue
			}
			fmt.Fprintf(os.Stderr, "find b\n")
			self.path = now
			return
		}
		if now.name == "#" {
			if m != 0 && now.findDis < 2 {
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				continue
			}
			fmt.Fprintf(os.Stderr, "find #\n")
			// fmt.Println(string(now.name))
			self.path = now
			return
		}

		if self.mob.me.param1 != 0 && dagger != 0 {
			if now.name == "m" {
				fmt.Fprintf(os.Stderr, "find m\n")
				self.path = now
				return
			}
		}

		for _, i := range now.neighbors {

			if i.touch {
				continue
			}
			if i.monsaura {
				m++
				continue
			}
			// if len(i.neighbors) == 2 {
			// 	if i.neighbors[0].monsaura || i.neighbors[1].monsaura {
			// 		continue
			// 	}
			// }
			i.parent = now
			i.findDis = now.findDis + 1

			if i.name == "d" {
				fmt.Fprintf(os.Stderr, "find d\n")
				self.path = i
				return
			}
			if i.name == "b" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				fmt.Fprintf(os.Stderr, "find b\n")
				self.path = i
				return
			}
			if i.name == "#" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				fmt.Fprintf(os.Stderr, "find #\n")
				// fmt.Println(string(now.name))
				self.path = i
				return
			}

			if self.mob.me.param1 != 0 && dagger != 0 {
				if now.name == "m" {
					fmt.Fprintf(os.Stderr, "find m\n")
					self.path = i
					return
				}
			}
			open = append(open, i)
		}
		now.touch = true
	}
}

func (self *base) bfsMonter() {
	var open []*box
	open = append(open, self.maps[self.mob.me.y][self.mob.me.x])
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		fmt.Fprintf(os.Stderr, now.name)
		if now.touchMons {
			continue
		}
		if now.name == "m" {
			fmt.Fprintf(os.Stderr, "find m\n")
			// fmt.Println(string(now.name))
			self.path = now
			return
		}

		for _, i := range now.neighbors {

			if i.touchMons {
				continue
			}
			i.parent = now
			i.distans = now.distans + 1
			if i.name == "m" {
				fmt.Fprintf(os.Stderr, "find m\n")
				// fmt.Println(string(now.name))

				self.path = i
				return
			}
			open = append(open, i)
		}
		now.touchMons = true
	}
}
