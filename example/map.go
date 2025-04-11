package main

import (
	"astar"
	"fmt"
)

const MapXSize = 12
const MapYSize = 5

func main() {
	graphRuneList := graphList2graphRuneList(graphList)
	for _, g := range graphRuneList {
		m := &Map{}
		m.RenderFromGraph(g[0])
		fmt.Printf("raw:\n")
		for _, r := range m.CellGrid {
			for _, c := range r {
				fmt.Printf("%v", string(c.graphRune))
			}
			fmt.Printf("\n")
		}
		paths, found := astar.Paths(m.startCell, m.endCell)
		if found {
			for _, v := range paths {
				m.getCell(v.(*Cell).X, v.(*Cell).Y).graphRune = '●'
				//fmt.Printf("%v", string(v.(*Cell).graphRune))
			}
		}
		fmt.Printf("route found result:%v\n", found)
		for _, r := range m.CellGrid {
			for _, c := range r {
				fmt.Printf("%v", string(c.graphRune))
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}
}

/*
	. - 路
	X - 障碍
	F - 起点
	T - 终点
	● - 寻路的路径
*/

var graphList = [][2]string{
	{
		`
.....X......
.....XX.....
.F........T.
....XXX.....
............
`,
		`
.....X......
.....XX.....
.F●●●●●●●●T.
....XXX.....
............
`,
	},
	{
		`
............
.....XX.....
.F..XXXX..T.
....XXX.....
............
`,
		`
............
.....XX.....
.●●●XXXX●●●.
...●XXX●●...
...●●●●●....
`,
	},
	{
		`
............      
.........XXX
.F.......XTX
.........XXX
............
`,
		`
............
.........XXX
.F.......XTX
.........XXX
............
`,
	},
	{
		`
FX.X........
.X...XXXX.X.
.X.X.X....X.
...X.X.XXXXX
.XX..X.....T
`,
		`
●X.X●●●●●●..
●X●●●XXXX●X.
●X●X.X●●●●X.
●●●X.X●XXXXX
.XX..X●●●●●●
`,
	},
}

func graphList2graphRuneList(srcList [][2]string) [][2][MapYSize][MapXSize]rune {
	dstList := make([][2][MapYSize][MapXSize]rune, 0, len(srcList))
	parseGraphString2RuneFun := func(str string) [MapYSize][MapXSize]rune {
		var retGraph [MapYSize][MapXSize]rune
		i := 0
		for _, v := range str {
			if v == '\n' || v == ' ' {
				continue
			}
			row := i / MapXSize
			col := i % MapXSize
			retGraph[row][col] = v
			i++
		}
		return retGraph
	}
	for _, graphInfo := range srcList {
		rawGraph := graphInfo[0]
		routeGraph := graphInfo[1]
		dstList = append(dstList, [2][MapYSize][MapXSize]rune{
			parseGraphString2RuneFun(rawGraph), parseGraphString2RuneFun(routeGraph)})
	}
	return dstList
}

// ===========================================================

// 模拟一次地图寻路
type Map struct {
	CellGrid  [MapYSize][MapXSize]*Cell
	startCell *Cell
	endCell   *Cell
}

// 地图单元格
type Cell struct {
	m         *Map
	X, Y      int
	graphRune rune
}

func newCell(m *Map, x, y int, r rune) *Cell {
	c := &Cell{
		m:         m,
		X:         x,
		Y:         y,
		graphRune: r,
	}
	return c
}

// PathNeighbors 查找单元格周围上下左右4个格子（排除障碍的）
func (c *Cell) PathNeighbors() []astar.IMapCell {
	return c.m.getPointNeighbors(c.X, c.Y)
}

// PathGCost 计算到一个邻居节点的确切消耗，里面可以加入路程的消耗权重，例如平路1、水路2、山路3
func (c *Cell) PathGCost(to astar.IMapCell) float64 {
	// 这个模拟只能上下、左右走，所以写死1步吧
	return 1
}

// PathHCost 计算到终点节点的最小代价评估值
func (t *Cell) PathHCost(to astar.IMapCell) float64 {
	// 曼哈顿距离
	toT := to.(*Cell)
	absX := toT.X - t.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Y - t.Y
	if absY < 0 {
		absY = -absY
	}
	r := float64(absX + absY)

	return r
}

// RenderFromGraph 从一个字符串图形渲染出二维地图
func (m *Map) RenderFromGraph(graph [MapYSize][MapXSize]rune) {
	for ri, r := range graph {
		for ci, c := range r {
			cell := newCell(m, ci, ri, c)
			if c == 'F' {
				m.startCell = cell
			} else if c == 'T' {
				m.endCell = cell
			}
			m.CellGrid[ri][ci] = cell
		}
	}
	return
}

func (m *Map) getCell(x, y int) *Cell {
	return m.CellGrid[y][x]
}

func (m *Map) getPointNeighbors(x, y int) []astar.IMapCell {
	if x < 0 || x >= MapXSize || y < 0 || y >= MapYSize {
		return nil
	}
	list := make([]astar.IMapCell, 0)
	if x > 0 {
		c := m.getCell(x-1, y)
		if c.graphRune != 'X' {
			list = append(list, c)
		}
	}
	if x < MapXSize-1 {
		c := m.getCell(x+1, y)
		if c.graphRune != 'X' {
			list = append(list, c)
		}
	}
	if y > 0 {
		c := m.getCell(x, y-1)
		if c.graphRune != 'X' {
			list = append(list, c)
		}
	}
	if y < MapYSize-1 {
		c := m.getCell(x, y+1)
		if c.graphRune != 'X' {
			list = append(list, c)
		}
	}

	return list
}
