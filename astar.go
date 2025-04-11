package astar

type IMapCell interface {
	// PathNeighbors 查找单元格周围上下左右4个格子（排除障碍的）
	PathNeighbors() []IMapCell
	// PathGCost 计算到一个节点的确切消耗，里面可以加入路程的消耗权重，例如平路1、水路2、山路3
	PathGCost(to IMapCell) float64
	// PathHCost 计算到终点节点的最小代价评估值
	PathHCost(to IMapCell) float64
}

type calcCell struct {
	cell   IMapCell
	parent *calcCell
	g      float64
	h      float64

	sortIndex int
}

// getPaths 寻路成功后，通过目标单元格倒退获取行走路线的所有单元格
func (cell *calcCell) getPaths() []IMapCell {
	tmpCell := cell
	paths := make([]IMapCell, 0)
	for tmpCell.parent != nil {
		paths = append(paths, tmpCell.parent.cell)
		tmpCell = tmpCell.parent
	}
	return paths
}

// calcCellMgr 记录计算单元格的管理器，主要寻路逻辑里会修改单元格g、h、parent数据，
// 所以需要通过IMapCell获取到同一个单元格指针
type calcCellMgr struct {
	records map[IMapCell]*calcCell
}

func newCalcCellMgr() *calcCellMgr {
	return &calcCellMgr{records: make(map[IMapCell]*calcCell)}
}

func (mgr *calcCellMgr) get(cell IMapCell) *calcCell {
	tmp := mgr.records[cell]
	if tmp == nil {
		tmp = &calcCell{cell: cell, sortIndex: -1}
		mgr.records[cell] = tmp
	}
	return tmp
}

func Paths(from, to IMapCell) ([]IMapCell, bool) {
	cellsMgr := newCalcCellMgr()

	// 准备探索的路线
	openList := newSortedList()
	// 探索过的路线
	closedList := make(map[*calcCell]bool)

	// 将起始格子加入探索列表
	openList.addCell(cellsMgr.get(from))

	var paths []IMapCell
	found := false
	for !openList.isEmpty() {
		// 从准备探索的路线选一个最小F值的单元格，值越小路线越优
		curCalcCell := openList.getMinFCell()

		if curCalcCell.cell == to {
			// 走到目标点了，返回找到的路线
			found = true
			paths = curCalcCell.getPaths()
			break
		}

		for _, neighborICell := range curCalcCell.cell.PathNeighbors() {
			neighborCalcCell := cellsMgr.get(neighborICell)

			// 邻居单元格已经找过了，跳过
			if _, find := closedList[neighborCalcCell]; find {
				continue
			}

			// 计算当前邻居格子距离g，是否小于邻居格子上一次记录的到起点距离g，
			// 小于说明是更优路线，更新邻居节点g和parent
			neighborOldCost := neighborCalcCell.g
			neighborNewCost := curCalcCell.g + curCalcCell.cell.PathGCost(neighborCalcCell.cell)
			if !openList.hasCell(neighborCalcCell) || neighborNewCost < neighborOldCost {
				// 当前邻居路线不在openList就可以作为接下来探索的路径，或者当前路线更优旧路线作废
				if openList.hasCell(neighborCalcCell) {
					openList.delCell(neighborCalcCell)
				}
				neighborCalcCell.g = neighborNewCost
				neighborCalcCell.h = neighborCalcCell.cell.PathHCost(to)
				// 切换邻居节点的父节点为当前节点，因为当前路线更优
				neighborCalcCell.parent = curCalcCell
				openList.addCell(neighborCalcCell)
			}
		}

		openList.delCell(curCalcCell)
		closedList[curCalcCell] = true
	}

	return paths, found
}
