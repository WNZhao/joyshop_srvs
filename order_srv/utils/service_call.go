package utils

import (
	"context"
	"fmt"
	"order_srv/global"
	goodspb "order_srv/proto/goods"
	inventorypb "order_srv/proto/inventory"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetGoodsByIds 批量获取商品信息
func GetGoodsByIds(ctx context.Context, goodsIds []int32) (map[int32]*goodspb.GoodsInfoResponse, error) {
	if global.GoodsClient == nil {
		return nil, fmt.Errorf("商品服务未连接")
	}

	goodsClient := goodspb.NewGoodsClient(global.GoodsClient)
	
	// 调用商品服务批量获取商品信息
	resp, err := goodsClient.BatchGetGoods(ctx, &goodspb.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		global.Logger.Errorf("调用商品服务失败: %v", err)
		return nil, fmt.Errorf("获取商品信息失败: %w", err)
	}

	// 转换为map便于查找
	goodsMap := make(map[int32]*goodspb.GoodsInfoResponse)
	for _, good := range resp.Data {
		goodsMap[good.Id] = good
	}

	return goodsMap, nil
}

// SellInventory 扣减库存
func SellInventory(ctx context.Context, sellItems []*inventorypb.GoodsInvInfo) error {
	if global.InventoryClient == nil {
		return fmt.Errorf("库存服务未连接")
	}

	inventoryClient := inventorypb.NewInventoryServiceClient(global.InventoryClient)
	
	// 调用库存服务扣减库存
	_, err := inventoryClient.Sell(ctx, &inventorypb.SellInfo{
		GoodsInvInfo: sellItems,
	})
	if err != nil {
		// 检查是否是库存不足的错误
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.ResourceExhausted {
				global.Logger.Warnf("库存不足: %v", err)
				return fmt.Errorf("库存不足")
			}
		}
		global.Logger.Errorf("库存扣减失败: %v", err)
		return fmt.Errorf("库存扣减失败: %w", err)
	}

	global.Logger.Info("库存扣减成功")
	return nil
}

// RebackInventory 归还库存（用于订单取消等场景）
func RebackInventory(ctx context.Context, rebackItems []*inventorypb.GoodsInvInfo) error {
	if global.InventoryClient == nil {
		return fmt.Errorf("库存服务未连接")
	}

	inventoryClient := inventorypb.NewInventoryServiceClient(global.InventoryClient)
	
	// 调用库存服务归还库存
	_, err := inventoryClient.Reback(ctx, &inventorypb.SellInfo{
		GoodsInvInfo: rebackItems,
	})
	if err != nil {
		global.Logger.Errorf("库存归还失败: %v", err)
		return fmt.Errorf("库存归还失败: %w", err)
	}

	global.Logger.Info("库存归还成功")
	return nil
}

// ValidateGoodsAvailability 验证商品是否可用（上架、有库存等）
func ValidateGoodsAvailability(goodsInfo *goodspb.GoodsInfoResponse, requiredNum int32) error {
	// 检查商品是否上架
	if !goodsInfo.OnSale {
		return fmt.Errorf("商品 %s 已下架", goodsInfo.Name)
	}
	
	// 检查库存是否足够
	if goodsInfo.Stocks < requiredNum {
		return fmt.Errorf("商品 %s 库存不足，当前库存: %d，需要: %d", 
			goodsInfo.Name, goodsInfo.Stocks, requiredNum)
	}

	return nil
}