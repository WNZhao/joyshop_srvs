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
	t.Log(TestScenarios.GoodsDetail)
	
	// 测试获取高价商品详情（iPhone）
	iphone := GetTestGoodsByID(1)
	if iphone != nil {
		t.Run("iPhone商品详情", func(t *testing.T) {
			rsp, err := goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: iphone.ID})
			if err != nil {
				t.Errorf("获取iPhone商品详情失败: %v", err)
				return
			}
			
			// 验证商品信息
			if rsp.Id != iphone.ID {
				t.Errorf("商品ID不匹配: 期望 %d, 实际 %d", iphone.ID, rsp.Id)
			}
			if rsp.Name != iphone.Name {
				t.Errorf("商品名称不匹配: 期望 %s, 实际 %s", iphone.Name, rsp.Name)
			}
			if rsp.ShopPrice != iphone.ShopPrice {
				t.Errorf("商品价格不匹配: 期望 %.2f, 实际 %.2f", iphone.ShopPrice, rsp.ShopPrice)
			}
			
			t.Logf("iPhone商品详情 - ID: %d, 名称: %s, 价格: ￥%.2f, 品牌: %s", 
				rsp.Id, rsp.Name, rsp.ShopPrice, iphone.BrandName)
		})
	}
	
	// 测试获取低价商品详情（小米手环）
	miBand := GetTestGoodsByID(9)
	if miBand != nil {
		t.Run("小米手环商品详情", func(t *testing.T) {
			rsp, err := goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: miBand.ID})
			if err != nil {
				t.Errorf("获取小米手环商品详情失败: %v", err)
				return
			}
			
			// 验证商品信息
			if rsp.Id != miBand.ID {
				t.Errorf("商品ID不匹配: 期望 %d, 实际 %d", miBand.ID, rsp.Id)
			}
			if rsp.IsHot != miBand.IsHot {
				t.Errorf("热销标记不匹配: 期望 %v, 实际 %v", miBand.IsHot, rsp.IsHot)
			}
			if rsp.IsNew != miBand.IsNew {
				t.Errorf("新品标记不匹配: 期望 %v, 实际 %v", miBand.IsNew, rsp.IsNew)
			}
			
			t.Logf("小米手环商品详情 - ID: %d, 名称: %s, 价格: ￥%.2f, 热销: %v, 新品: %v", 
				rsp.Id, rsp.Name, rsp.ShopPrice, rsp.IsHot, rsp.IsNew)
		})
	}
}

// 测试品牌相关功能
func TestBrandOperations(t *testing.T) {
	t.Log(TestScenarios.BrandList)
	
	// 测试品牌列表查询
	t.Run("品牌列表查询", func(t *testing.T) {
		rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
			Pages:       1,
			PagePerNums: 10,
		})
		if err != nil {
			t.Errorf("获取品牌列表失败: %v", err)
			return
		}
		
		t.Logf("品牌列表总数: %d", rsp.Total)
		for _, brand := range rsp.Data {
			t.Logf("品牌信息 - ID: %d, 名称: %s", brand.Id, brand.Name)
		}
	})
	
	// 测试特定品牌查询
	for _, testBrand := range TestBrands[:5] { // 测试前5个品牌
		t.Run(fmt.Sprintf("品牌_%s", testBrand.Name), func(t *testing.T) {
			// 获取该品牌下的商品
			goodsRsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
				Pages:       1,
				PagePerNums: 10,
				BrandId:     testBrand.ID,
			})
			if err != nil {
				t.Errorf("获取品牌 %s 下的商品失败: %v", testBrand.Name, err)
				return
			}
			
			t.Logf("品牌 %s 下的商品数量: %d", testBrand.Name, goodsRsp.Total)
			for _, goods := range goodsRsp.Data {
				if goods.BrandId != testBrand.ID {
					t.Errorf("商品品牌不匹配: 商品ID %d, 期望品牌ID %d, 实际品牌ID %d", 
						goods.Id, testBrand.ID, goods.BrandId)
				}
			}
		})
	}
}

// 测试分类相关功能
func TestCategoryOperations(t *testing.T) {
	t.Log(TestScenarios.CategoryList)
	
	// 测试获取所有分类
	t.Run("所有分类查询", func(t *testing.T) {
		rsp, err := goodsClient.GetAllCategorysList(context.Background(), &proto.Empty{})
		if err != nil {
			t.Errorf("获取所有分类失败: %v", err)
			return
		}
		
		t.Logf("分类总数: %d", rsp.Total)
		
		// 按级别统计分类
		levelCount := make(map[int32]int)
		for _, category := range rsp.Data {
			levelCount[category.Level]++
			t.Logf("分类信息 - ID: %d, 名称: %s, 级别: %d, 父级ID: %d", 
				category.Id, category.Name, category.Level, category.ParentId)
		}
		
		for level, count := range levelCount {
			t.Logf("%d级分类数量: %d", level, count)
		}
	})
	
	// 测试子分类查询
	t.Run("子分类查询", func(t *testing.T) {
		// 测试电子数码分类的子分类
		electronicCategory := GetTestCategoryByID(1) // 电子数码
		if electronicCategory != nil {
			rsp, err := goodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
				Id: electronicCategory.ID,
			})
			if err != nil {
				t.Errorf("获取电子数码子分类失败: %v", err)
				return
			}
			
			t.Logf("电子数码子分类数量: %d", len(rsp.SubCategorys))
			for _, subCategory := range rsp.SubCategorys {
				if subCategory.ParentId != electronicCategory.ID {
					t.Errorf("子分类父级ID不正确: 子分类ID %d, 期望父级ID %d, 实际父级ID %d",
						subCategory.Id, electronicCategory.ID, subCategory.ParentId)
				}
				t.Logf("子分类信息 - ID: %d, 名称: %s", subCategory.Id, subCategory.Name)
			}
		}
	})
}

// 测试价格区间商品查询
func TestPriceRangeQuery(t *testing.T) {
	t.Log("价格区间商品查询测试")
	
	priceRanges := []struct {
		Name     string
		MinPrice float32
		MaxPrice float32
	}{
		{"低价商品", PriceRanges.Low.Min, PriceRanges.Low.Max},
		{"中价商品", PriceRanges.Medium.Min, PriceRanges.Medium.Max},
		{"高价商品", PriceRanges.High.Min, PriceRanges.High.Max},
	}
	
	for _, priceRange := range priceRanges {
		t.Run(priceRange.Name, func(t *testing.T) {
			rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
				Pages:       1,
				PagePerNums: 10,
				PriceMin:    priceRange.MinPrice,
				PriceMax:    priceRange.MaxPrice,
			})
			if err != nil {
				t.Errorf("查询%s失败: %v", priceRange.Name, err)
				return
			}
			
			t.Logf("%s数量: %d", priceRange.Name, rsp.Total)
			
			// 验证价格区间
			for _, goods := range rsp.Data {
				if goods.ShopPrice < priceRange.MinPrice || goods.ShopPrice > priceRange.MaxPrice {
					t.Errorf("商品价格不在指定区间: 商品ID %d, 价格 %.2f, 期望区间 [%.2f, %.2f]",
						goods.Id, goods.ShopPrice, priceRange.MinPrice, priceRange.MaxPrice)
				}
				t.Logf("%s - ID: %d, 名称: %s, 价格: ¥%.2f", 
					priceRange.Name, goods.Id, goods.Name, goods.ShopPrice)
			}
		})
	}
}

// 测试组合条件查询
func TestCombinedQuery(t *testing.T) {
	t.Log("组合条件查询测试")
	
	// 测试品牌 + 分类 + 价格区间的组合查询
	t.Run("Apple品牌电子产品", func(t *testing.T) {
		appleBrand := GetTestBrandByID(1)
		phoneCategory := GetTestCategoryByID(9)
		
		if appleBrand != nil && phoneCategory != nil {
			rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
				Pages:       1,
				PagePerNums: 10,
				BrandId:     appleBrand.ID,
				CategoryId:  phoneCategory.ID,
				PriceMin:    5000.0,
				PriceMax:    50000.0,
			})
			if err != nil {
				t.Errorf("组合查询失败: %v", err)
				return
			}
			
			t.Logf("Apple高端电子产品数量: %d", rsp.Total)
			for _, goods := range rsp.Data {
				// 验证条件
				if goods.BrandId != appleBrand.ID {
					t.Errorf("品牌不匹配: 商品ID %d", goods.Id)
				}
				if goods.ShopPrice < 5000 || goods.ShopPrice > 50000 {
					t.Errorf("价格不在范围内: 商品ID %d, 价格 %.2f", goods.Id, goods.ShopPrice)
				}
				
				t.Logf("Apple高端产品 - ID: %d, 名称: %s, 价格: ¥%.2f", 
					goods.Id, goods.Name, goods.ShopPrice)
			}
		}
	})
	
	// 测试热销 + 新品的组合查询
	t.Run("热销新品", func(t *testing.T) {
		rsp, err := goodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			IsHot:       true,
			IsNew:       true,
		})
		if err != nil {
			t.Errorf("热销新品查询失败: %v", err)
			return
		}
		
		t.Logf("热销新品数量: %d", rsp.Total)
		for _, goods := range rsp.Data {
			if !goods.IsHot {
				t.Errorf("非热销商品出现在结果中: 商品ID %d", goods.Id)
			}
			if !goods.IsNew {
				t.Errorf("非新品商品出现在结果中: 商品ID %d", goods.Id)
			}
			
			t.Logf("热销新品 - ID: %d, 名称: %s, 价格: ¥%.2f", 
				goods.Id, goods.Name, goods.ShopPrice)
		}
	})
}
