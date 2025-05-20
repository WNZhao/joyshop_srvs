package tests

import (
	"context"
	"fmt"
	"goods_srv/proto"
	"goods_srv/util"
	"testing"
)

// 测试获取商品列表
func TestGoodsList(t *testing.T) {
	// 测试基本分页
	t.Run("基本分页", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 2,
		})
		if err != nil {
			t.Errorf("获取商品列表失败: %v", err)
			return
		}
		t.Logf("商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("商品信息 - ID: %d, 名称: %s, 价格: %.2f", goods.Id, goods.Name, goods.ShopPrice)
		}
	})

	// 测试按品牌筛选
	t.Run("按品牌筛选", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			BrandId:     1, // 假设品牌ID为1
		})
		if err != nil {
			t.Errorf("按品牌获取商品列表失败: %v", err)
			return
		}
		t.Logf("品牌商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("品牌商品信息 - ID: %d, 名称: %s, 品牌ID: %d", goods.Id, goods.Name, goods.BrandId)
		}
	})

	// 测试按分类筛选
	t.Run("按分类筛选", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			CategoryId:  1, // 假设分类ID为1
		})
		if err != nil {
			t.Errorf("按分类获取商品列表失败: %v", err)
			return
		}
		t.Logf("分类商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("分类商品信息 - ID: %d, 名称: %s, 分类IDs: %v", goods.Id, goods.Name, goods.CategoryIds)
		}
	})

	// 测试按关键词搜索
	t.Run("按关键词搜索", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			Keywords:    "测试", // 假设搜索关键词为"手机"
		})
		if err != nil {
			t.Errorf("按关键词获取商品列表失败: %v", err)
			return
		}
		t.Logf("搜索结果总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("搜索结果 - ID: %d, 名称: %s", goods.Id, goods.Name)
		}
	})

	// 测试热门商品
	t.Run("热门商品", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			IsHot:       true,
		})
		if err != nil {
			t.Errorf("获取热门商品列表失败: %v", err)
			return
		}
		t.Logf("热门商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("热门商品信息 - ID: %d, 名称: %s, 是否热门: %v", goods.Id, goods.Name, goods.IsHot)
		}
	})

	// 测试新品
	t.Run("新品", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			IsNew:       true,
		})
		if err != nil {
			t.Errorf("获取新品列表失败: %v", err)
			return
		}
		t.Logf("新品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("新品信息 - ID: %d, 名称: %s, 是否新品: %v", goods.Id, goods.Name, goods.IsNew)
		}
	})

	// 测试在售商品
	t.Run("在售商品", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			OnSale:      true,
		})
		if err != nil {
			t.Errorf("获取在售商品列表失败: %v", err)
			return
		}
		t.Logf("在售商品总数: %d", rsp.Total)
		for _, goods := range rsp.Data {
			t.Logf("在售商品信息 - ID: %d, 名称: %s, 是否在售: %v", goods.Id, goods.Name, goods.OnSale)
		}
	})
}

// 商品增删改测试用例，mock品牌和分类数据严格复用test_utils.go规则
func TestGoodsCRUD(t *testing.T) {
	// 清理商品相关表数据
	if err := util.CleanGoodsRelatedTables(); err != nil {
		t.Fatalf("清理商品相关表数据失败: %v", err)
	}

	// 1. 批量插入品牌
	var brandIds []int32
	brandCount := 5

	for i := 0; i < brandCount; i++ {
		name, desc := generateRandomBrand()
		brandReq := &proto.BrandRequest{
			Name: name,
			Desc: desc,
		}
		brandResp, err := goodsClient.CreateBrand(context.Background(), brandReq)
		if err != nil {
			t.Fatalf("批量创建品牌失败: %v", err)
		}
		brandIds = append(brandIds, brandResp.Id)
	}
	t.Logf("批量创建品牌ID: %v", brandIds)

	// 2. 递归插入分类，并收集所有三级分类ID
	var categoryIds []int32
	var collectLevel3 func(parentId int32, cats []map[string]interface{})
	collectLevel3 = func(parentId int32, cats []map[string]interface{}) {
		for _, cat := range cats {
			req := &proto.CategoryInfoRequest{
				Name:     cat["Name"].(string),
				Level:    int32(cat["Level"].(int)),
				IsTab:    cat["IsTab"].(bool),
				ParentId: parentId,
			}
			resp, err := goodsClient.CreateCategory(context.Background(), req)
			if err != nil {
				t.Fatalf("创建分类失败: %v", err)
			}
			if req.Level == 3 {
				categoryIds = append(categoryIds, resp.Id)
			}
			if children, ok := cat["SubCategories"].([]map[string]interface{}); ok && len(children) > 0 {
				collectLevel3(resp.Id, children)
			}
		}
	}
	cats := GenerateTaobaoCategories()
	collectLevel3(0, cats)
	t.Logf("所有三级分类ID: %v", categoryIds)

	// 3. 批量添加商品
	goodsIds := make([]int32, 0)
	batchCount := 10
	for i := 0; i < batchCount; i++ {
		brandId := brandIds[i%len(brandIds)]
		catId := categoryIds[i%len(categoryIds)]
		goodsReq := &proto.CreateGoodsInfo{
			Name:            fmt.Sprintf("测试商品-%d", i+1),
			GoodsSn:         fmt.Sprintf("SN%04d", i+1),
			Stocks:          int32(100 + i*10),
			MarketPrice:     float32(1000 + i*100),
			ShopPrice:       float32(900 + i*80),
			GoodsBrief:      "这是测试商品简要描述",
			GoodsDesc:       "这是测试商品详细描述",
			Images:          []string{"http://example.com/img1.jpg", "http://example.com/img2.jpg"},
			DescImages:      []string{"http://example.com/desc1.jpg"},
			GoodsFrontImage: "http://example.com/front.jpg",
			Status:          1,
			IsHot:           i%2 == 0,
			IsNew:           i%3 == 0,
			OnSale:          true,
			ShipFree:        i%2 == 1,
			BrandId:         brandId,
			CategoryIds:     []int32{catId},
		}
		goodsResp, err := goodsClient.CreateGoods(context.Background(), goodsReq)
		if err != nil {
			t.Fatalf("批量添加商品失败: %v", err)
		}
		goodsIds = append(goodsIds, goodsResp.Id)
	}
	t.Logf("批量添加商品ID: %v", goodsIds)

	// 4. 单条添加商品
	singleGoodsReq := &proto.CreateGoodsInfo{
		Name:            "单条测试商品",
		GoodsSn:         "SN9999",
		Stocks:          50,
		MarketPrice:     1999.99,
		ShopPrice:       1888.88,
		GoodsBrief:      "单条商品简要描述",
		GoodsDesc:       "单条商品详细描述",
		Images:          []string{"http://example.com/single1.jpg"},
		DescImages:      []string{"http://example.com/singledesc.jpg"},
		GoodsFrontImage: "http://example.com/singlefront.jpg",
		Status:          1,
		IsHot:           true,
		IsNew:           true,
		OnSale:          true,
		ShipFree:        false,
		BrandId:         brandIds[0],
		CategoryIds:     []int32{categoryIds[0]},
	}
	singleGoodsResp, err := goodsClient.CreateGoods(context.Background(), singleGoodsReq)
	if err != nil {
		t.Fatalf("单条添加商品失败: %v", err)
	}
	goodsIds = append(goodsIds, singleGoodsResp.Id)
	t.Logf("单条添加商品ID: %d", singleGoodsResp.Id)

	// 5. 修改商品
	if len(goodsIds) > 0 {
		updateId := goodsIds[0]
		updateReq := &proto.CreateGoodsInfo{
			Id:              updateId,
			Name:            "修改后的商品名称",
			GoodsSn:         "SN0000",
			Stocks:          888,
			MarketPrice:     8888.88,
			ShopPrice:       7777.77,
			GoodsBrief:      "修改后的简要描述",
			GoodsDesc:       "修改后的详细描述",
			Images:          []string{"http://example.com/updated1.jpg"},
			DescImages:      []string{"http://example.com/updateddesc.jpg"},
			GoodsFrontImage: "http://example.com/updatedfront.jpg",
			Status:          2,
			IsHot:           false,
			IsNew:           false,
			OnSale:          false,
			ShipFree:        true,
			BrandId:         brandIds[0],
			CategoryIds:     []int32{categoryIds[0]},
		}

		_, err := goodsClient.UpdateGoods(context.Background(), updateReq)
		if err != nil {
			t.Errorf("修改商品失败: %v", err)
		} else {
			t.Logf("修改商品成功 - ID: %d", updateId)
		}
	}

	// 6. 删除商品
	if len(goodsIds) > 1 {
		deleteId := goodsIds[1]
		_, err := goodsClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: deleteId})
		if err != nil {
			t.Errorf("删除商品失败: %v", err)
		} else {
			t.Logf("删除商品成功 - ID: %d", deleteId)
		}
	}
}

// 测试获取商品详情
func TestGetGoodsDetail(t *testing.T) {
	// 测试获取商品详情
	rsp, err := goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: 14})
	if err != nil {
		t.Errorf("获取商品详情失败: %v", err)
		return
	}
	t.Logf("商品详情: %v", rsp)
}
