package core

import (
	"fmt"
)

//定义地图结构体
type AOIManager struct {
	//区域的左边边界
	MinX int
	//区域的右边边界
	MaxX int
	//X轴方向格子的数量
	CntsX int
	//区域的上边边界
	MinY int
	//区域的下边边界
	MaxY int
	//Y轴方向的格子的数量
	CntsY int
	//整体区域(地图中)拥有哪些格子map:key格子ID,value:格子对象
	grids map[int]*Grid
}

//得到每个格子在X轴的宽度
func (m *AOIManager) GridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在Y轴的高度
func (m *AOIManager) GridHeight() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//初始化一个地图
//minX, maxX, cntsX, minY, maxY, cntsY int
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	//隶属于当前地图的全部格子 也一并进行初始化
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//初始化一个格子
			//格子ID := cntsX * y + x
			gid := y*cntsX + x

			//给aoiManager(地图)添加一个格子
			aoiMgr.grids[gid] = NewGid(gid,
				aoiMgr.MinX+x*aoiMgr.GridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.GridWidth(),
				aoiMgr.MinY+y*aoiMgr.GridHeight(),
				aoiMgr.MinY+(y+1)*aoiMgr.GridHeight())
		}
	}

	return aoiMgr
}

//打印当前地图信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n Minx:%d , MaxX:%d , cntsX:%d ,minY:%d, maxY:%d,cntsY:%d,Grids inManager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)

	//打印全部格子信息
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

//添加一个playerID到一个AOI格子中
func (m *AOIManager) AddPidToGrid(PID, gID int) {
	m.grids[gID].Add(PID, nil)
}

//移除一个PlayerID 从一个AOI区域中
func (m *AOIManager) RemovePidFromGrid(PID, gID int) {
	m.grids[gID].Remove(PID)
}

//通过格子ID获取当前格子的全部PlayerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

//通过一个格子ID得到当前格子的周边九宫格的格子集合
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//判断gid是否在AOI中
	if _, ok := m.grids[gID]; !ok {
		return
	}

	//将当前中心GID放入九宫格
	grids = append(grids, m.grids[gID])

	//==== > 判读GID左边是否有格子？ 右边是否有格子

	//通过格子ID 得到x轴编号 idx = gID % cntsX
	idx := gID % m.CntsX
	//判断idx编号左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	//判断idx编号右边是否还有格子
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}

	// ===> 得到一个x轴的格子集合，遍历这个格子集合
	// for ... 依次判断  格子的上面是否有格子？下面是否有格子

	//将X轴全部的Grid ID 放到一个slice中 ，遍历整个slice
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	for _, gid := range gidsX {
		//10,11,12
		//通过Gid得到当前Gid的Y轴编号
		//idy = gID / cntsX
		idy := gid / m.CntsX

		//上方是否有格子
		if idy > 0 {
			grids = append(grids, m.grids[gid-m.CntsX])
		}
		//下方是否有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[gid+m.CntsX])
		}
	}

	//按照gid大小排序九宫格gird
	for i := 0; i < len(grids)-1; i++ {
		for j := 0; j < len(grids)-i-1; j++ {
			if grids[j].GID > grids[j+1].GID {
				grids[j], grids[j+1] = grids[j+1], grids[j]
			}
		}
	}

	return
}

//通过x，y坐标得到对应的格子ID
func (m *AOIManager) GetGidByPos(x, y float32) int {
	if x < 0 || int(x) >= m.MaxX {
		return -1
	}
	if y < 0 || int(y) >= m.MaxY {
		return -1
	}
	//根据坐标 得到当前玩家所在格子ID
	idx := (int(x) - m.MinX) / m.GridWidth()
	idy := (int(y) - m.MinY) / m.GridHeight()
	gID := idy*m.CntsX + idx

	return gID
}

//根据坐标,得到周边的九宫格中的全部玩家的 ID集合
func (m *AOIManager) GetSurroundPIDsByPos(x, y float32) (playerIDs []int) {

	//通过x, y得到一个格子对应的ID
	gid := m.GetGidByPos(x, y)

	//通过格子ID获取周边九宫格 集合
	grids := m.GetSurroundGridsByGid(gid)

	//fmt.Println("gid = ", gid)
	//fmt.Println("grids=", grids)

	//将分别将九宫格内的全部的玩家 放在 playerIDs
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
	}

	//fmt.Println("playerIDs =", playerIDs)

	return
}

//通过坐标 将玩家ID--->pid  添加到 当前坐标所定位的格子grid中
func (m *AOIManager) AddToGridByPos(PID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	//取出当前的格子
	grid := m.grids[gID]
	//给格子添加玩家
	grid.Add(PID, nil)
}

//通过坐标 把一个palyer 从一个个格子中删除
func (m *AOIManager) RemoteFromGridbyPos(PID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	//取出当前的格子
	grid := m.grids[gID]
	//给格子添加玩家
	grid.Remove(PID)
}
