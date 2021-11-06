package main

import (
	"fmt"
	"os"
)

func main() {
	var maps [][]*coord
	var maps2 [][]*coord
	daggerTime := 0
	for true {
		// fmt.Fprintf(os.Stderr, fmt.Sprintf("daggerTime = %v \n", daggerTime))
		if daggerTime != 0 {
			daggerTime--
		}
		var w, h, playerID, tick int
		fmt.Scan(&w, &h, &playerID, &tick)
		//init map
		// read map
		addNeig := false
		dagger := 0
		for i := 0; i < h; i++ {
			var line []*coord
			var line2 []*coord
			var l string
			fmt.Scan(&l)
			for j := 0; j < w; j++ {
				var c rune
				c = rune(l[j])
				if c == 'd' {
					dagger++
				}
				if len(maps) < h {
					var cord coord
					cord.name = c
					cord.x = i
					cord.y = j
					var cop coord
					cop = cord
					line = append(line, &cord)
					line2 = append(line2, &cop)
				} else {
					maps[i][j].name = c
					maps[i][j].touch = false
					maps[i][j].depth = 0
					maps[i][j].parent = nil
					maps2[i][j].name = c
					maps2[i][j].touch = false
					maps2[i][j].depth = 0
					maps2[i][j].parent = nil
				}
			}
			if len(maps) < h {
				maps = append(maps, line)
				maps2 = append(maps2, line2)
				addNeig = true
			}
		}

		// add neighbors
		if addNeig {
			maps = addneighbors(maps, w, h)
		}
		// number of entities
		var n int
		fmt.Scan(&n)

		// read entities
		player := user{}
		for i := 0; i < n; i++ {
			var entType string
			var pID, x, y, param1, param2 int
			fmt.Scan(&entType, &pID, &x, &y, &param1, &param2)
			if entType != "p" {
				maps[y][x].name = 'm'
				if daggerTime == 0 {
					maps = monsterAgreZone(maps, x, y, 2)
					maps2 = monsterAgreZone(maps2, x, y, 1)
				}
				continue
			}
			player.name = entType
			player.pID = pID
			player.x = x
			player.y = y
			player.param1 = param1
			player.param2 = param2
			fmt.Fprintf(os.Stderr, fmt.Sprintf("entType = %v pID %v x %v y %v param1 %v param2 %v \n", entType, pID, x, y, param1, param2))

		}

		fmt.Fprintf(os.Stderr, fmt.Sprintf("player.x = %v, player.y = %v \n", player.x, player.y))

		//find path

		var path *coord
		if dagger == 0 || daggerTime != 0 {
			path = bfs(maps, player.x, player.y, '#')
		} else {
			path = bfs(maps, player.x, player.y, 'd')
			if path != nil {
				if path.parent == maps[player.y][player.x] {
					daggerTime = 15
				}
			}
		}

		if path == nil {
			fmt.Fprintf(os.Stderr, "BFS nil pointer\n")
			// m := false
			// p := maps[player.y][player.x]
			// k := p
			// for j := 0; j <= len(p.neighbors); j++ {
			// for i := 0; i < len(p.neighbors); i++ {
			// 	if p.neighbors[i].name == 'm' {
			// 		// m = true
			// 		continue
			// 	}
			// 	path = p.neighbors[i]
			// 	path.parent = p
			// }
			// 	if j < len(k.neighbors) {
			// 		p = k.neighbors[j]
			// 	}
			// }

			path = bfs(maps2, player.x, player.y, '#')

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

		if finalx == 1 {
			fmt.Println("up")
			fmt.Fprintf(os.Stderr, "up\n")
			continue
		}
		if finalx == -1 {
			fmt.Println("down")
			fmt.Fprintf(os.Stderr, "down\n")

			continue
		}
		if finaly == -1 {
			fmt.Println("right")
			fmt.Fprintf(os.Stderr, "right\n")

			continue
		}
		if finaly == 1 {
			fmt.Println("left")
			fmt.Fprintf(os.Stderr, "left\n")

			continue
		}
		fmt.Println("stay")
		fmt.Fprintf(os.Stderr, "stay\n")

	}
}

func bfs(maps [][]*coord, x, y int, triger rune) *coord {
	// var close []*coord
	var open []*coord
	open = append(open, maps[y][x])
	for len(open) != 0 {
		now := open[0]
		open = open[1:]
		if now.touch && len(open) != 0 {
			continue
		}
		for _, i := range now.neighbors {

			if i.name == 'm' {
				continue
			}
			if i.touch {
				continue
			}
			i.parent = now

			if i.name == triger {
				return i
			}
			i.depth = now.depth + 1
			open = append(open, i)
		}
		now.touch = true
	}
	return nil
}

func monsterAgreZone(maps [][]*coord, x, y, aur int) [][]*coord {
	for i := -1 * aur; i <= aur; i++ {
		for j := -1 * aur; j <= aur; j++ {
			if (i*i == 4 && j != 0) || (j*j == 4 && i != 0) {
				continue
			}
			xn := x + i
			yn := y + j
			if xn >= 0 && xn < len(maps[0]) && yn >= 0 && yn < len(maps) {
				maps[yn][xn].touch = true
				maps[yn][xn].coin = 100

			}
		}
	}
	return maps
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

type coord struct {
	x         int
	y         int
	name      rune
	touch     bool
	depth     int
	coin      int
	parent    *coord
	neighbors []*coord
}

type user struct {
	name   string
	pID    int
	x      int
	y      int
	param1 int
	param2 int
}
