/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 18:04:02
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 19:35:21
 * @FilePath: /joyshop_srvs/goods_srv/tests/category_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// ... existing code ...
// 直接引用全局goodsClient，无需TestMain。
// ... existing code ...

package tests

import (
	"context"
	"goods_srv/proto"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"
)

// 递归创建多级分类
func createCategoryTree(t *testing.T, parentId int32, categories []map[string]interface{}) {
	for _, cat := range categories {
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
		// 递归创建子分类
		if children, ok := cat["SubCategories"].([]map[string]interface{}); ok && len(children) > 0 {
			createCategoryTree(t, resp.Id, children)
		}
	}
}

func TestCategory(t *testing.T) {
	// 1. 创建淘宝风格三级分类
	t.Run("创建淘宝风格三级分类", func(t *testing.T) {
		categories := GenerateTaobaoCategories()
		createCategoryTree(t, 0, categories)
	})

	// 2. 获取所有分类
	t.Run("获取所有分类", func(t *testing.T) {
		resp, err := goodsClient.GetAllCategoriesList(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Fatalf("获取所有分类失败: %v", err)
		}
		if resp.Total == 0 {
			t.Error("分类列表为空")
		}
	})

	// 3. 获取某一级的子分类
	t.Run("获取子分类", func(t *testing.T) {
		// 假设第一个大类的ID为1
		req := &proto.CategoryListRequest{
			Id:    1,
			Level: 1,
		}
		resp, err := goodsClient.GetSubCategory(context.Background(), req)
		if err != nil {
			t.Fatalf("获取子分类失败: %v", err)
		}
		if resp.Info == nil || resp.Info.Id != 1 {
			t.Error("未找到正确的父分类信息")
		}
	})

	// 4. 更新分类
	t.Run("更新分类", func(t *testing.T) {
		req := &proto.CategoryInfoRequest{
			Id:    1,
			Name:  "家用电器-更新",
			Level: 1,
			IsTab: true,
		}
		_, err := goodsClient.UpdateCategory(context.Background(), req)
		if err != nil {
			t.Fatalf("更新分类失败: %v", err)
		}
	})

	// 5. 删除分类
	t.Run("删除分类", func(t *testing.T) {
		req := &proto.DeleteCategoryRequest{
			Id: 1, // 假设删除第一个大类
		}
		_, err := goodsClient.DeleteCategory(context.Background(), req)
		if err != nil {
			t.Fatalf("删除分类失败: %v", err)
		}
	})
}
