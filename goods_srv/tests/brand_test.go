package tests

import (
	"context"
	"goods_srv/proto"
	"testing"
)

// 测试品牌列表
func TestBrandList(t *testing.T) {
	// 测试基本分页
	t.Run("基本分页", func(t *testing.T) {
		rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
			Pages:       1,
			PagePerNums: 2,
		})
		if err != nil {
			t.Errorf("获取品牌列表失败: %v", err)
			return
		}
		t.Logf("品牌总数: %d", rsp.Total)
		for _, brand := range rsp.Data {
			t.Logf("品牌信息 - ID: %d, 名称: %s, 描述: %s", brand.Id, brand.Name, brand.Desc)
		}
	})

	// 测试按名称过滤
	t.Run("按名称过滤", func(t *testing.T) {
		// 先创建一个测试品牌
		name, desc := generateRandomBrand()
		createReq := &proto.BrandRequest{
			Name: name,
			Desc: desc,
		}
		createRsp, err := goodsClient.CreateBrand(context.Background(), createReq)
		if err != nil {
			t.Errorf("创建测试品牌失败: %v", err)
			return
		}

		// 使用创建的品牌名称进行查询
		rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			Name:        name,
		})
		if err != nil {
			t.Errorf("按名称查询品牌列表失败: %v", err)
			return
		}
		t.Logf("按名称查询结果总数: %d", rsp.Total)
		for _, brand := range rsp.Data {
			t.Logf("查询结果 - ID: %d, 名称: %s, 描述: %s", brand.Id, brand.Name, brand.Desc)
		}

		// 清理测试数据
		_, err = goodsClient.DeleteBrand(context.Background(), &proto.BrandRequest{Id: createRsp.Id})
		if err != nil {
			t.Errorf("清理测试品牌失败: %v", err)
		}
	})

	// 测试按描述过滤
	t.Run("按描述过滤", func(t *testing.T) {
		// 先创建一个测试品牌
		name, desc := generateRandomBrand()
		createReq := &proto.BrandRequest{
			Name: name,
			Desc: desc,
		}
		createRsp, err := goodsClient.CreateBrand(context.Background(), createReq)
		if err != nil {
			t.Errorf("创建测试品牌失败: %v", err)
			return
		}

		// 使用创建的品牌描述进行查询
		rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			Desc:        desc,
		})
		if err != nil {
			t.Errorf("按描述查询品牌列表失败: %v", err)
			return
		}
		t.Logf("按描述查询结果总数: %d", rsp.Total)
		for _, brand := range rsp.Data {
			t.Logf("查询结果 - ID: %d, 名称: %s, 描述: %s", brand.Id, brand.Name, brand.Desc)
		}

		// 清理测试数据
		_, err = goodsClient.DeleteBrand(context.Background(), &proto.BrandRequest{Id: createRsp.Id})
		if err != nil {
			t.Errorf("清理测试品牌失败: %v", err)
		}
	})

	// 测试组合条件过滤
	t.Run("组合条件过滤", func(t *testing.T) {
		// 先创建一个测试品牌
		name, desc := generateRandomBrand()
		createReq := &proto.BrandRequest{
			Name: name,
			Desc: desc,
		}
		createRsp, err := goodsClient.CreateBrand(context.Background(), createReq)
		if err != nil {
			t.Errorf("创建测试品牌失败: %v", err)
			return
		}

		// 使用名称和描述组合查询
		rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
			Pages:       1,
			PagePerNums: 10,
			Name:        name,
			Desc:        desc,
		})
		if err != nil {
			t.Errorf("组合条件查询品牌列表失败: %v", err)
			return
		}
		t.Logf("组合条件查询结果总数: %d", rsp.Total)
		for _, brand := range rsp.Data {
			t.Logf("查询结果 - ID: %d, 名称: %s, 描述: %s", brand.Id, brand.Name, brand.Desc)
		}

		// 清理测试数据
		_, err = goodsClient.DeleteBrand(context.Background(), &proto.BrandRequest{Id: createRsp.Id})
		if err != nil {
			t.Errorf("清理测试品牌失败: %v", err)
		}
	})
}

// 测试创建品牌
func TestCreateBrand(t *testing.T) {
	t.Run("创建新品牌", func(t *testing.T) {
		name, desc := generateRandomBrand()
		req := &proto.BrandRequest{
			Name: name,
			Desc: desc,
		}
		rsp, err := goodsClient.CreateBrand(context.Background(), req)
		if err != nil {
			t.Errorf("创建品牌失败: %v", err)
			return
		}
		t.Logf("创建品牌成功 - ID: %d, 名称: %s, 描述: %s", rsp.Id, rsp.Name, rsp.Desc)
	})
}

// 测试更新品牌
func TestUpdateBrand(t *testing.T) {
	t.Run("更新品牌", func(t *testing.T) {
		req := &proto.BrandRequest{
			Id:   1, // 假设要更新的品牌ID为1
			Name: "更新后的品牌名称",
			Desc: "更新后的品牌描述",
		}
		_, err := goodsClient.UpdateBrand(context.Background(), req)
		if err != nil {
			t.Errorf("更新品牌失败: %v", err)
			return
		}
		t.Logf("更新品牌成功 - ID: %d", req.Id)
	})
}

// 测试删除品牌
func TestDeleteBrand(t *testing.T) {
	t.Run("删除品牌", func(t *testing.T) {
		req := &proto.BrandRequest{
			Id: 1, // 假设要删除的品牌ID为1
		}
		_, err := goodsClient.DeleteBrand(context.Background(), req)
		if err != nil {
			t.Errorf("删除品牌失败: %v", err)
			return
		}
		t.Logf("删除品牌成功 - ID: %d", req.Id)
	})
}

// 测试品牌完整流程
func TestBrandFullProcess(t *testing.T) {
	var brandId int32

	// 1. 创建品牌
	t.Run("创建品牌", func(t *testing.T) {
		name, desc := generateRandomBrand()
		req := &proto.BrandRequest{
			Name: name,
			Desc: desc,
		}
		rsp, err := goodsClient.CreateBrand(context.Background(), req)
		if err != nil {
			t.Errorf("创建品牌失败: %v", err)
			return
		}
		brandId = rsp.Id
		t.Logf("创建品牌成功 - ID: %d, 名称: %s, 描述: %s", brandId, rsp.Name, rsp.Desc)
	})

	// 2. 更新品牌
	t.Run("更新品牌", func(t *testing.T) {
		name, desc := generateRandomBrand()
		req := &proto.BrandRequest{
			Id:   brandId,
			Name: name,
			Desc: desc,
		}
		_, err := goodsClient.UpdateBrand(context.Background(), req)
		if err != nil {
			t.Errorf("更新品牌失败: %v", err)
			return
		}
		t.Logf("更新品牌成功 - ID: %d, 新名称: %s, 新描述: %s", brandId, name, desc)
	})

	// 3. 查询品牌列表
	t.Run("查询品牌列表", func(t *testing.T) {
		rsp, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
			Pages:       1,
			PagePerNums: 10,
		})
		if err != nil {
			t.Errorf("获取品牌列表失败: %v", err)
			return
		}
		t.Logf("品牌总数: %d", rsp.Total)
		for _, brand := range rsp.Data {
			t.Logf("品牌信息 - ID: %d, 名称: %s, 描述: %s", brand.Id, brand.Name, brand.Desc)
		}
	})

	// 4. 删除品牌
	t.Run("删除品牌", func(t *testing.T) {
		req := &proto.BrandRequest{
			Id: brandId,
		}
		_, err := goodsClient.DeleteBrand(context.Background(), req)
		if err != nil {
			t.Errorf("删除品牌失败: %v", err)
			return
		}
		t.Logf("删除品牌成功 - ID: %d", brandId)
	})
}
