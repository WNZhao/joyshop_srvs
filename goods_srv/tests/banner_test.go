package tests

import (
	"context"
	"testing"

	"goods_srv/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

func TestBanner(t *testing.T) {
	// 测试创建轮播图
	t.Run("创建轮播图", func(t *testing.T) {
		req := &proto.BannerRequest{
			Image: "http://example.com/banner1.jpg",
			Url:   "http://example.com/page1",
			Index: 1,
		}

		resp, err := goodsClient.CreateBanner(context.Background(), req)
		if err != nil {
			t.Fatalf("创建轮播图失败: %v", err)
		}

		if resp.Image != req.Image {
			t.Errorf("期望图片URL为 %s, 实际为 %s", req.Image, resp.Image)
		}
		if resp.Url != req.Url {
			t.Errorf("期望链接URL为 %s, 实际为 %s", req.Url, resp.Url)
		}
		if resp.Index != req.Index {
			t.Errorf("期望索引为 %d, 实际为 %d", req.Index, resp.Index)
		}

		// 保存创建的轮播图ID用于后续测试
		bannerId = resp.Id
	})

	// 测试获取轮播图列表
	t.Run("获取轮播图列表", func(t *testing.T) {
		resp, err := goodsClient.BannerList(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Fatalf("获取轮播图列表失败: %v", err)
		}

		if resp.Total == 0 {
			t.Error("轮播图列表为空")
		}

		// 验证列表中的轮播图数据
		found := false
		for _, banner := range resp.Data {
			if banner.Id == bannerId {
				found = true
				if banner.Image != "http://example.com/banner1.jpg" {
					t.Errorf("期望图片URL为 http://example.com/banner1.jpg, 实际为 %s", banner.Image)
				}
				break
			}
		}
		if !found {
			t.Error("未找到创建的轮播图")
		}
	})

	// 测试更新轮播图
	t.Run("更新轮播图", func(t *testing.T) {
		req := &proto.BannerRequest{
			Id:    bannerId,
			Image: "http://example.com/banner1_updated.jpg",
			Url:   "http://example.com/page1_updated",
			Index: 2,
		}

		_, err := goodsClient.UpdateBanner(context.Background(), req)
		if err != nil {
			t.Fatalf("更新轮播图失败: %v", err)
		}

		// 验证更新后的数据
		listResp, err := goodsClient.BannerList(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Fatalf("获取更新后的轮播图列表失败: %v", err)
		}

		found := false
		for _, banner := range listResp.Data {
			if banner.Id == bannerId {
				found = true
				if banner.Image != req.Image {
					t.Errorf("期望更新后的图片URL为 %s, 实际为 %s", req.Image, banner.Image)
				}
				if banner.Url != req.Url {
					t.Errorf("期望更新后的链接URL为 %s, 实际为 %s", req.Url, banner.Url)
				}
				if banner.Index != req.Index {
					t.Errorf("期望更新后的索引为 %d, 实际为 %d", req.Index, banner.Index)
				}
				break
			}
		}
		if !found {
			t.Error("未找到更新后的轮播图")
		}
	})

	// 测试删除轮播图
	t.Run("删除轮播图", func(t *testing.T) {
		req := &proto.BannerRequest{
			Id: bannerId,
		}

		_, err := goodsClient.DeleteBanner(context.Background(), req)
		if err != nil {
			t.Fatalf("删除轮播图失败: %v", err)
		}

		// 验证删除后的数据
		listResp, err := goodsClient.BannerList(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Fatalf("获取删除后的轮播图列表失败: %v", err)
		}

		// 检查是否还能找到已删除的轮播图
		for _, banner := range listResp.Data {
			if banner.Id == bannerId {
				t.Error("轮播图未被成功删除")
				break
			}
		}
	})
}

// 用于存储测试过程中创建的轮播图ID
var bannerId int32
