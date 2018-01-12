package lib

import (
	"fmt"
)

func TransFlowRedisDone(Account string, AimAccount string, WtAmt int64, OrderNo string) (err error) {
	redis := NewRedis("trans_flow")
	redis1 := NewRedis("xianshi_trans")
	redis2 := NewRedis(fmt.Sprintf("xianshi_users_%s", Account))
	// redis3 := NewRedis("xianshi_zonghui")
	// 删除redis，限时竞购表"xianshi_trans"
	err = redis.Delete(OrderNo)
	if err != nil {
		err = fmt.Errorf("[%s]删除redis结算流水表trans_flow错误[%v]", Account, err)
		return
	}

	// 删除redis，限时竞购表"xianshi_trans"
	err = redis1.DeleteSortStringList(ps("trans_flow_%s", OrderNo))
	if err != nil {
		err = fmt.Errorf("[%s]删除redis限时竞购表xianshi_trans数据失败[%v]", Account, err)
		return
	}

	// 删除redis，限时竞购表"xianshi_users_"
	redis2.DeleteSortStringList(ps("trans_flow_%s", OrderNo))
	if err != nil {
		err = fmt.Errorf("[%s]删除redis限时竞购表xianshi_users_数据失败[%v]", Account, err)
		return
	}

	// OnlineAmt, _ := redis3.GetInt64(AimAccount)
	// redis3.Put(AimAccount, OnlineAmt-WtAmt)
	return
}

func TransFlowRedisDelete(Account string, OrderNo string) (err error) {
	redis := NewRedis("trans_flow")
	redis1 := NewRedis("xianshi_trans")
	redis2 := NewRedis(fmt.Sprintf("xianshi_users_%s", Account))
	// 删除redis，限时竞购表"xianshi_trans"
	err = redis.Delete(OrderNo)
	if err != nil {
		err = fmt.Errorf("[%s]删除redis结算流水表trans_flow错误[%v]", Account, err)
		return
	}

	// 删除redis，限时竞购表"xianshi_trans"
	err = redis1.DeleteSortStringList(ps("trans_flow_%s", OrderNo))
	if err != nil {
		err = fmt.Errorf("[%s]删除redis限时竞购表xianshi_trans数据失败[%v]", Account, err)
		return
	}

	// 删除redis，限时竞购表"xianshi_users_"
	redis2.DeleteSortStringList(ps("trans_flow_%s", OrderNo))
	if err != nil {
		err = fmt.Errorf("[%s]删除redis限时竞购表xianshi_users_数据失败[%v]", Account, err)
		return
	}

	return
}
