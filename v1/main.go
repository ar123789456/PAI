package main

import (
	"fmt"
	"os"
)

func main() {
	var Maps mapsType
	paramBool := true
	for true {
		Maps := initMap(Maps)
		// add neighbors
		if Maps.addNeig {
			Maps.maps = addneighbors(Maps.maps, Maps.w, Maps.h)
			Maps.maps2 = addneighbors(Maps.maps2, Maps.w, Maps.h)
		}
		// number of entities
		var n int
		fmt.Scan(&n)

		// read entities
		player := user{}
		mobs := []user{}
		for i := 0; i < n; i++ {
			var entType string
			var pID, x, y, param1, param2 int
			fmt.Scan(&entType, &pID, &x, &y, &param1, &param2)
			if Maps.playerID == pID {
				paramBool = param1 == 0

				player.addParam(entType, pID, x, y, param1, param2)
			}
			if entType == "m" {
				Maps.maps[y][x].name = 'm'
				// Maps.maps2[0][0].name = 'q'

				g := user{}
				g.addParam(entType, pID, x, y, param1, param2)
				mobs = append(mobs, g)
			}

			fmt.Fprintf(os.Stderr, fmt.Sprintf("entType = %v pID %v x %v y %v param1 %v param2 %v \n", entType, pID, x, y, param1, param2))
		}

		fmt.Fprintf(os.Stderr, fmt.Sprintf("player.x = %v, player.y = %v \n", player.x, player.y))

		//add mobs agre zone

		for _, m := range mobs {
			if paramBool {
				Maps.maps = monsterAgreZone(Maps.maps, m.x, m.y, 3)
				Maps.maps2 = monsterAgreZone(Maps.maps2, m.x, m.y, 2)

			}
		}

		//find path

		var path *coord
		// var multiPath road
		ch := make(chan road)
		ch2 := make(chan road)
		go bfs(Maps.maps, player.x, player.y, ch)
		go bfs(Maps.maps2, player.x, player.y, ch2)
		multiPath, _ := <-ch
		multiPath2, _ := <-ch2
		// multiPath = bfs(Maps.maps, player.x, player.y)
		path = multiPath.optimal(&player, len(mobs))

		if path == nil {
			fmt.Fprintf(os.Stderr, "BFS nil pointer\n")
			path = multiPath2.optimal(&player, len(mobs))

			if path == nil {
				fmt.Println("stay")
				fmt.Fprintf(os.Stderr, "stay\n")
				continue

			}
		}
		var finalx, finaly int
		for path.parent != nil {
			finalx = player.y - path.x
			finaly = player.x - path.y
			path = path.parent
		}
		fmt.Fprintf(os.Stderr, fmt.Sprintf("finalx = %v, finaly = %v \n", finalx, finaly))
		// this will choose one of random actions
		PrintResult(finalx, finaly)
	}
}

func bfs(maps [][]*coord, x, y int, ch chan road) {
	var r road
	r.player = maps[y][x]
	var open []*coord
	open = append(open, maps[y][x])
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if now.touch && len(open) != 0 {
			continue
		}
		if len(now.neighbors) == 2 {
			if now.neighbors[0].touch && now.neighbors[1].touch {
				continue
			}
		}
		if now.name != '.' && now.name != 'q' {
			r.allroad = append(r.allroad, now)
		}
		if now.name == 'm' {
			if r.monster == nil {
				// fmt.Fprintf(os.Stderr, "m\n")
				r.monster = now
			}
		}
		if now.name == 'b' {
			if r.bonus == nil {
				// fmt.Fprintf(os.Stderr, "b\n")
				r.bonus = now
			}
		}
		if now.name == 'd' {
			if r.dagger == nil {
				// fmt.Fprintf(os.Stderr, "d\n")
				r.dagger = now
			}
		}
		if now.name == '#' {
			if r.gold == nil {
				// fmt.Fprintf(os.Stderr, "#\n")
				r.gold = now
			}
		}
		if now.name == 'q' {
			// fmt.Fprintf(os.Stderr, "q\n")
			r.quit = now
		}
		for _, i := range now.neighbors {
			if i.touch {
				continue
			}

			i.parent = now
			i.depth = now.depth + 1
			open = append(open, i)
		}
		now.touch = true
	}
	ch <- r
}

type road struct {
	player  *coord
	gold    *coord
	bonus   *coord
	dagger  *coord
	monster *coord
	quit    *coord
	allroad []*coord
}

func (self *road) retreat() *coord {
	if self.dagger != nil {
		return self.dagger
	}
	return self.quit
}

func (self *road) optimal(player *user, mon int) *coord {
	// if self.gold != nil {
	// 	k := 0
	// 	// for _, i := range self.gold.neighbors {
	// 	// 	if i.mons {
	// 	// 		k++
	// 	// 	}
	// 	// }
	// 	if k == len(self.gold.neighbors) {
	// 		self.gold = nil
	// 	}
	// }
	// if player.param1 == 0 {
	// 	if self.bonus != nil {
	// 		return self.bonus
	// 	} else if self.dagger != nil && mon != 0 {
	// 		if self.dagger.depth < 15 {
	// 			return self.dagger
	// 		} else {
	// 			return self.gold
	// 		}
	// 	} else if self.gold != nil {
	// 		return self.gold
	// 	} else if self.quit != nil {
	// 		return self.quit
	// 	}
	// } else {
	// 	if self.monster != nil {
	// 		if self.monster.depth < self.gold.depth {
	// 			return self.monster
	// 		}
	// 	}
	// 	return self.gold
	// }
	// return nil
	var rt *coord
	for _, i := range self.allroad {
		if rt == nil {
			rt = i
			continue
		}
		if i.name == 'b' || i.name == 'd' {
			if i.depth < rt.depth+2 {
				rt = i
				continue
			}
		}
		if i.depth < rt.depth {
			rt = i
		}
	}
	return rt

}

func initMap(Maps mapsType) mapsType {
	var w, h, playerID, tick int
	fmt.Scan(&w, &h, &playerID, &tick)
	//init map
	Maps.h = h
	Maps.w = w
	Maps.playerID = playerID
	Maps.tick = tick
	Maps.dagger = 0
	Maps.addNeig = false
	// read map
	for i := 0; i < h; i++ {
		var line []*coord
		var line2 []*coord
		var l string
		fmt.Scan(&l)
		for j := 0; j < w; j++ {
			var c rune
			c = rune(l[j])
			if c == 'd' {
				Maps.dagger++
			}
			if len(Maps.maps) < h { //add new cord
				var cord coord
				cord.make(c, i, j)
				var cop coord
				cop = cord
				line = append(line, &cord)
				line2 = append(line2, &cop)
			} else { // modification old map
				Maps.maps[i][j].modification(c)
				Maps.maps2[i][j].modification(c)
			}
		}
		if len(Maps.maps) < h {
			Maps.maps = append(Maps.maps, line)
			Maps.maps2 = append(Maps.maps2, line2)
			Maps.addNeig = true
		}
	}
	return Maps
}

func monsterAgreZone(maps [][]*coord, x, y, aur int) [][]*coord {
	open := []*coord{}
	maps[y][x].touch = true
	maps[y][x].mons = true
	open = append(open, maps[y][x])

	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		for _, i := range now.neighbors {
			if i.touch {
				continue
			}
			i.depth = now.depth + 1
			if i.depth == aur+1 {
				i.depth = 0
				continue
			}
			i.touch = true
			i.mons = true
			open = append(open, i)
			// fmt.Fprintf(os.Stderr, fmt.Sprintf("x = %v, y = %v, name = %v \n", i.x, i.y, string(i.name)))

		}
		now.depth = 0
	}
	return maps
}

func PrintResult(finalx, finaly int) {
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

func addneighbors(maps [][]*coord, w, h int) [][]*coord {
	for i, line := range maps {
		for j, cord := range line {
			if cord.name == '!' {
				continue
			}
			if i+1 != h {
				if i != 0 {
					if maps[i+1][j].name != '!' {
						maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i+1][j])
					}
					if maps[i-1][j].name != '!' {
						maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i-1][j])
					}
				} else {
					if maps[i+1][j].name != '!' {
						maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i+1][j])
					}
				}
			} else {
				if maps[i-1][j].name != '!' {
					maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i-1][j])
				}
			}
			if j+1 != w {
				if j != 0 {
					if maps[i][j+1].name != '!' {
						maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i][j+1])
					}
					if maps[i][j-1].name != '!' {
						maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i][j-1])
					}
				} else {
					if maps[i][j+1].name != '!' {
						maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i][j+1])
					}
				}
			} else {
				if maps[i][j-1].name != '!' {
					maps[i][j].neighbors = append(maps[i][j].neighbors, maps[i][j-1])
				}
			}
		}
	}
	return maps
}

type mapsType struct {
	w, h        int
	playerID    int
	tick        int
	maps, maps2 [][]*coord
	addNeig     bool
	dagger      int
}

type coord struct {
	x         int
	y         int
	name      rune
	touch     bool
	mons      bool
	depth     int
	coin      int
	parent    *coord
	neighbors []*coord
}

func (self *coord) make(c rune, x, y int) {
	self.name = c
	self.x = x
	self.y = y
}

func (self *coord) modification(c rune) {
	self.name = c
	self.touch = false
	self.depth = 0
	self.parent = nil
}

type user struct {
	name   string
	pID    int
	x      int
	y      int
	param1 int
	param2 int
}

func (self *user) addParam(entType string, pID, x, y, param1, param2 int) {
	self.name = entType
	self.pID = pID
	self.x = x
	self.y = y
	self.param1 = param1
	self.param2 = param2
}
