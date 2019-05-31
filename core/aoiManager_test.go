package core

import (
	"testing"
	"fmt"
)

func TestAOIManager_init(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	//打印信息
	fmt.Println(aoiMgr)
}

func TestAOIManagerSurround(t *testing.T) {
	//初始化AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	//求出周边九宫格子信息
	for gid, _ := range aoiMgr.grids {
		grids := aoiMgr.GetSurroundGridsByGid(gid)
		fmt.Println("gid : ", gid, " grids num = ", len(grids))

		//当前九宫格的ID集合
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}

		fmt.Println("grids IDs are ", gIDs)
	}

	//根据坐标却定gid
	gid := aoiMgr.GetGidByPos(175, 68)
	fmt.Println("gid = ",gid)

	//根据坐标确定周围九宫格gid

}
