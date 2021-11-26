package main

import (
	"fmt"
	"os"
	"time"
)

const monsterRadius = 3

var monsterRadiusDagger = 1

var dagger int

var bonus int

var myGold int
var gold int
var enamyGold int

var daggerTimer int

var pastBaseInfo *base
var monsterPath [][]int

func main() {
	for {
		if daggerTimer != 0 {
			daggerTimer--
		}
		if dagger != 0 {
			dagger--
		}
		if bonus != 0 {
			bonus--
		}
		start := time.Now()
		var baseInfo base

		baseInfo.initmaps()
		baseInfo.addneighbors()
		// baseInfo.maps[0][0].name = "q"

		if len(monsterPath) == 0 {
			baseInfo.createMonsterBool()
		}

		var mobs Mobs
		mobs.initMobs(&baseInfo)
		baseInfo.mob = &mobs
		score(baseInfo.mob.me, baseInfo.mob.enamy)
		pastBaseInfo = &baseInfo
		fmt.Fprintf(os.Stderr, fmt.Sprintf("my= %v, enamy = %v \n", myGold, enamyGold))

		for _, v := range mobs.monster {
			if mobs.me.param1 == 0 && mobs.me.param2 != 2 {
				baseInfo = mobsAura(baseInfo, *v, monsterRadius)
			}
			if mobs.me.param2 == 2 {
				baseInfo.maps[v.y][v.x].touch = true
			}
		}

		if mobs.enamy.param2 == 2 && mobs.me.param2 != 3 {
			baseInfo.enamyAura()
		}

		baseInfo.bfs()
		baseInfo.OptimalRoad()

		if baseInfo.path == nil {
			baseInfo.OptimalRoad2()
		}

		baseInfo.PrintResult()
		duration := time.Since(start)

		fmt.Fprintf(os.Stderr, fmt.Sprintf("\n%v\n", duration))

	}
}

func score(me, enamy *Mob) {
	if pastBaseInfo == nil {
		return
	}
	if enamy == nil {
		return
	}
	if pastBaseInfo.maps[me.y][me.x].name == "#" {
		if me.param2 != 0 {
			myGold += 2
		} else {
			myGold++
		}
	}
	if pastBaseInfo.maps[enamy.y][enamy.x].name == "#" {
		if enamy.param2 != 0 {
			enamyGold += 2
		} else {
			enamyGold++
		}
	}
	if pastBaseInfo.maps[me.y][me.x].name == "d" {
		daggerTimer = 14
	}
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
			monsterPath[mob.y][mob.x]++
		} else if baseinfo.playerID == mob.pID {
			self.me = &mob
		} else {
			self.enamy = &mob
		}
	}
}

type base struct {
	maps                 [][]*box
	w, h, playerID, tick int
	mob                  *Mobs
	path                 *box
}

func (self *base) createMonsterBool() {
	for i := 0; i < self.h; i++ {
		list := []int{}
		for j := 0; j < self.w; j++ {
			list = append(list, 0)
		}
		monsterPath = append(monsterPath, list)
	}
}

func (self *base) OptimalRoad() {
	if gold == 1 && enamyGold < myGold-2 {
		fmt.Fprintf(os.Stderr, "more gold than the enemy \n")
		self.path = self.mob.me.dagger
		if self.path == nil {
			self.path = self.mob.me.frost
		}
		if self.path == nil {
			self.path = self.mob.me.bonus
		}
		return
	}
	if self.mob.me.bonus != nil && self.mob.me.bonus.findDis < bonus {
		self.path = self.mob.me.bonus
		return
	}
	if self.mob.enamy.param2 == 2 {
		if self.mob.me.imun != nil {
			self.path = self.mob.me.imun
			return
		}
	}
	if self.mob.me.frost != nil {
		self.path = self.mob.me.frost
		return
	}
	if self.mob.me.dagger != nil && self.mob.me.dagger.findDis < dagger {
		if self.mob.me.gold != nil && self.mob.me.dagger.findDis <= self.mob.me.gold.findDis+2 {
			if len(self.mob.monster) != 0 {
				self.path = self.mob.me.dagger
				return
			}
		}
	}
	self.path = self.mob.me.gold
}

func (self *base) OptimalRoad2() {
	if gold == 1 && enamyGold < myGold {
		fmt.Fprintf(os.Stderr, "more gold than the enemy \n")
		self.path = self.mob.me.ndagger
		if self.path == nil {
			self.path = self.mob.me.nbonus
		}
		return
	}
	if self.mob.me.nbonus != nil && self.mob.me.nbonus.findDis < bonus {
		self.path = self.mob.me.nbonus
		return
	}
	if self.mob.me.ndagger != nil && self.mob.me.ndagger.findDis < dagger {
		if self.mob.me.ngold != nil && self.mob.me.ndagger.findDis <= self.mob.me.ngold.findDis+2 {
			if len(self.mob.monster) != 0 {
				self.path = self.mob.me.ndagger
				return
			}
		}
	}
	self.path = self.mob.me.ngold
}

func (self *base) PrintResult() {
	var x, y int
	actions := []string{"left", "right", "up", "down", "stay"}
	if self.path == nil {
		self.bfsRun()
		self.OptimalRoad()
	}
	if self.path != nil {
		if self.path.parent != nil {
			for self.path.parent.parent != nil {
				self.path = self.path.parent
			}
		} else {
			self.path = self.maps[self.mob.me.y][self.mob.me.x]
		}
	} else {
		fmt.Fprintf(os.Stderr, "Monster!!! Run\n")
		me := self.maps[self.mob.me.y][self.mob.me.x]
		for j := len(me.neighbors) - 1; j >= 0; j-- {
			var i *box
			if self.tick%2 == 0 {
				i = me.neighbors[j]
			} else {
				i = me.neighbors[j]
			}
			if me.x == 6 && me.y < self.h-2 && me.y > 2 && monsterPath[me.y][me.x] == 0 {
				if monsterPath[0][6] == 0 {
					self.path = me
					break
				}
			}
			if !me.monsaura && len(me.neighbors) != 1 {
				self.path = me
			}

			nei := i.neighbors
			if len(nei) == 2 {
				if nei[0].x != nei[1].x {
					if nei[0].y != nei[1].y {
						continue
					}
				}
			}

			if len(i.neighbors) == 1 {
				continue
			}
			if self.path == nil {
				self.path = i
			}
			if self.path.price > i.price {
				self.path = i
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

func (self *base) initmaps() {
	var w, h, playerID, tick int
	fmt.Scan(&w, &h, &playerID, &tick)
	// fmt.Fprintf(os.Stderr, fmt.Sprintf("%v %v %v %v\n", w, h, playerID, tick))
	self.h = h
	self.w = w
	self.playerID = playerID
	self.tick = tick
	// read map
	gold = 0
	bon := 0
	dag := 0
	for i := 0; i < h; i++ {
		line := ""
		var l []*box
		fmt.Scan(&line)
		// fmt.Fprint(os.Stderr, line, "\n")
		for x, name := range line {
			var b box
			b.name = string(name)
			b.x = x
			b.y = i
			if b.name == "d" {
				dag++
				if dagger == 0 {
					dagger = 14
				}
				b.price = -10
			}
			if b.name == "b" {
				bon++
				if bonus == 0 {
					bonus = 14
				}
			}
			if b.name == "#" {
				gold++
				// b.price = -1
			}
			l = append(l, &b)
		}
		self.maps = append(self.maps, l)
	}
	if bon == 0 {
		bonus = 0
	}
	if dag == 0 {
		dag = 0
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

			if y2 < self.h {
				if self.maps[y2][x].name != "!" {
					v.neighbors = append(v.neighbors, self.maps[y2][x])
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
			if x1 >= 0 {
				if self.maps[y][x1].name != "!" {
					v.neighbors = append(v.neighbors, self.maps[y][x1])
				}
			}

		}
	}
}

func (self *base) enamyAura() {
	open := []*box{}
	closed := []*box{}
	x := self.mob.enamy.x
	y := self.mob.enamy.y
	open = append(open, self.maps[y][x])
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if findV(closed, now) {
			continue
		}
		if now.radius == 3 {
			continue
		}
		for _, v := range now.neighbors {
			v.radius = 0
			v.radius = now.radius + 1
			open = append(open, v)
		}
		now.touch = true
		closed = append(closed, now)
	}
}

func mobsAura(baseinfo base, mob Mob, radius int) base {
	open := []*box{}
	closed := []*box{}
	baseinfo.maps[mob.y][mob.x].radius = 0
	open = append(open, baseinfo.maps[mob.y][mob.x])
	baseinfo.maps[mob.y][mob.x].touch = true
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if findV(closed, now) {
			continue
		}
		if now.radius == radius {
			now.price += (radius - now.radius + 1)
			me := baseinfo.mob.me
			if me.x == now.x && me.y == now.y {
				now.monsaura = true
			}
			continue
		}
		now.monsaura = true
		now.price += (radius - now.radius + 1)
		closed = append(closed, now)
		for _, v := range now.neighbors {
			v.radius = 0
			v.radius = now.radius + 1
			open = append(open, v)
		}
	}
	return baseinfo
}

func findV(closed []*box, v *box) bool {
	for _, value := range closed {
		if value == v {
			return true
		}
	}
	return false
}

type Mob struct {
	name                                      string
	pID, x, y, param1, param2                 int
	gold, bonus, dagger, monster, frost, imun *box
	ngold, nbonus, ndagger, nmonster          *box
}

type box struct {
	name         string
	y            int
	x            int
	parent       *box
	neighbors    []*box
	monsaura     bool
	radius       int
	runtouch     bool
	touch        bool
	touchenamy   bool
	touchMons    bool
	site         bool
	distans      int
	findDis      int
	findDisenamy int
	price        int
}

func (self *base) bfs() {
	fmt.Fprintf(os.Stderr, "activate BFS\n")
	var open []*box
	open = append(open, self.maps[self.mob.me.y][self.mob.me.x])
	m := 0
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		// fmt.Fprintf(os.Stderr, now.name)
		if now.touch {
			continue
		}
		if now.monsaura {
			now.touch = true
			m++
			continue
		}
		if now.name == "d" {
			// fmt.Fprintf(os.Stderr, "find d\n")
			if self.mob.me.dagger == nil {
				self.mob.me.dagger = now
			}
		}
		if now.name == "b" {
			if m != 0 && now.findDis < 2 {
				now.touch = true
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				now.touch = true
				continue
			}
			// fmt.Fprintf(os.Stderr, "find b\n")
			if self.mob.me.bonus == nil {
				self.mob.me.bonus = now
			}
		}
		if now.name == "f" {
			if m != 0 && now.findDis < 2 {
				now.touch = true
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				now.touch = true
				continue
			}
			// fmt.Fprintf(os.Stderr, "find b\n")
			if self.mob.me.frost == nil {
				self.mob.me.frost = now
			}
		}
		if now.name == "i" {
			if m != 0 && now.findDis < 2 {
				now.touch = true
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				now.touch = true
				continue
			}
			// fmt.Fprintf(os.Stderr, "find b\n")
			if self.mob.me.imun == nil {
				self.mob.me.imun = now
			}
		}
		if now.name == "#" {
			if m != 0 && now.findDis < 2 {
				now.touch = true
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				now.touch = true
				continue
			}
			if now.findDis <= now.findDisenamy && gold != 1 {
				now.touch = true
				if self.mob.me.ngold == nil {
					self.mob.me.ngold = now
				}
				continue
			}
			// fmt.Fprintf(os.Stderr, "find #\n")
			if self.mob.me.gold == nil {
				self.mob.me.gold = now
			}
		}
		now.touch = true
		for _, i := range now.neighbors {
			// if len(i.neighbors) == 2 && (i.neighbors[0].monsaura || i.neighbors[0].monsaura) {
			// 	continue
			// }
			if i.touch {
				continue
			}
			if i.monsaura {
				m++
				continue
			}
			if now.parent != i {
				i.parent = now
			} else {
				continue
			}
			i.findDis = now.findDis + 1

			if i.name == "d" {
				// fmt.Fprintf(os.Stderr, "find d\n")
				if self.mob.me.dagger == nil && len(self.mob.monster) != 0 {
					self.mob.me.dagger = i
				}
			}
			if i.name == "b" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				// fmt.Fprintf(os.Stderr, "find b\n")
				if self.mob.me.bonus == nil {
					self.mob.me.bonus = i
				}
			}
			if i.name == "f" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				// fmt.Fprintf(os.Stderr, "find b\n")
				if self.mob.me.frost == nil {
					self.mob.me.frost = i
				}
			}
			if i.name == "i" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				// fmt.Fprintf(os.Stderr, "find b\n")
				if self.mob.me.imun == nil {
					self.mob.me.imun = i
				}
			}

			if i.name == "#" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				if i.findDis <= i.findDisenamy && gold != 1 {
					i.touch = true
					if self.mob.me.ngold == nil {
						self.mob.me.ngold = i
					}
					continue
				}
				if self.mob.me.gold == nil {
					self.mob.me.gold = i
				}
			}
			open = append(open, i)
		}
	}
}

func (self *base) bfsRun() {
	fmt.Fprintf(os.Stderr, "BFSRun!!! Run\n")
	var open []*box
	open = append(open, self.maps[self.mob.me.y][self.mob.me.x])
	m := 0
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if now.touch {
			continue
		}
		if now.monsaura && len(open) != 0 {
			m++
			continue
		}
		if now.name == "d" {
			if now.radius == 1 && now.monsaura {
				continue
			}
			if self.mob.me.dagger == nil {
				self.mob.me.dagger = now
			}
		}
		if now.name == "b" {
			if m != 0 && now.findDis < 2 {
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				continue
			}
			if self.mob.me.bonus == nil {
				self.mob.me.bonus = now
			}
		}
		if now.name == "f" {
			if m != 0 && now.findDis < 2 {
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				continue
			}
			if self.mob.me.frost == nil {
				self.mob.me.frost = now
			}
		}
		if now.name == "i" {
			if m != 0 && now.findDis < 2 {
				continue
			}
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				continue
			}
			if self.mob.me.imun == nil {
				self.mob.me.imun = now
			}
		}
		if now.name == "#" {
			if len(now.neighbors) == 1 && now.neighbors[0].monsaura {
				continue
			}
			if self.mob.me.gold == nil {
				self.mob.me.gold = now
			}
		}
		now.touch = true

		for _, i := range now.neighbors {
			nei := i.neighbors
			if len(nei) == 2 {
				if nei[0].x != nei[1].x {
					if nei[0].y != nei[1].y {
						continue
					}
				}
			}
			if i.touch {
				continue
			}
			if len(i.neighbors) == 1 {
				continue
			}
			if i.monsaura {
				m++
				continue
			}
			if now.parent != i {
				i.parent = now
			} else {
				continue
			}
			i.findDis = now.findDis + 1

			if i.name == "d" {
				if self.mob.me.dagger == nil {
					self.mob.me.dagger = i
				}
			}
			if i.name == "b" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				if self.mob.me.bonus == nil {
					self.mob.me.bonus = i
				}
			}
			if i.name == "f" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				if self.mob.me.frost == nil {
					self.mob.me.frost = i
				}
			}
			if i.name == "i" {
				if m != 0 && i.findDis < 2 {
					continue
				}
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				if self.mob.me.imun == nil {
					self.mob.me.imun = i
				}
			}
			if i.name == "#" {
				if len(i.neighbors) == 1 && i.neighbors[0].monsaura {
					continue
				}
				if self.mob.me.gold == nil {
					self.mob.me.gold = i
				}
			}
			open = append(open, i)
		}
	}
}

func (self *base) enamybfs() {
	fmt.Fprintf(os.Stderr, "activate BFS\n")
	var open []*box
	open = append(open, self.maps[self.mob.me.y][self.mob.me.x])

	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if now.touchenamy {
			continue
		}
		now.touchenamy = true
		for _, i := range now.neighbors {
			if i.touchenamy {
				continue
			}
			i.findDisenamy = now.findDisenamy + 1

			open = append(open, i)
		}
	}
}
