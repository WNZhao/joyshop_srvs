/*
 * @Author: Will nanan_zhao@163.com
 * @Date: 2025-05-18 14:06:45
 * @LastEditors: Will nanan_zhao@163.com
 * @LastEditTime: 2025-05-18 18:04:07
 * @FilePath: /joyshop_srvs/goods_srv/tests/test_utils.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package tests

import (
	"fmt"
	"math/rand"
)

// 创建一个本地随机数生成器
var rnd = rand.New(rand.NewSource(rand.Int63()))

// 品牌名称前缀列表
var brandNamePrefixes = []string{
	"优品", "精选", "品质", "优选", "精品", "高端", "时尚", "经典", "创新", "卓越",
	"领先", "专业", "优质", "顶级", "尊享", "奢华", "精致", "完美", "独特", "非凡",
}

// 品牌名称后缀列表
var brandNameSuffixes = []string{
	"科技", "电子", "数码", "家居", "服饰", "美妆", "食品", "运动", "户外", "母婴",
	"箱包", "鞋靴", "配饰", "家电", "厨具", "文具", "玩具", "宠物", "园艺", "汽车",
}

// 品牌描述模板
var brandDescTemplates = []string{
	"专注于%s领域的%s品牌，致力于为消费者提供高品质的产品和服务。",
	"作为%s行业的%s品牌，我们始终坚持创新和品质。",
	"以%s为核心，打造%s领域的专业品牌形象。",
	"致力于%s领域，成为%s行业的领导品牌。",
	"专注于%s产品研发，打造%s领域的优质品牌。",
}

// 从切片中随机选择一个元素
func randomChoice(slice []string) string {
	return slice[rnd.Intn(len(slice))]
}

// 生成随机品牌名称
func generateRandomBrandName() string {
	prefix := randomChoice(brandNamePrefixes)
	suffix := randomChoice(brandNameSuffixes)
	return fmt.Sprintf("%s%s", prefix, suffix)
}

// 生成随机品牌描述
func generateRandomBrandDesc() string {
	template := randomChoice(brandDescTemplates)
	prefix := randomChoice(brandNamePrefixes)
	suffix := randomChoice(brandNameSuffixes)
	return fmt.Sprintf(template, suffix, prefix)
}

// 生成随机品牌信息
func generateRandomBrand() (name, desc string) {
	return generateRandomBrandName(), generateRandomBrandDesc()
}

// 生成淘宝风格的三级分类数据
func GenerateTaobaoCategories() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"Name":  "家用电器",
			"Level": 1,
			"IsTab": true,
			"SubCategories": []map[string]interface{}{
				{
					"Name":  "电视",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "智能电视", "Level": 3, "IsTab": false},
						{"Name": "曲面电视", "Level": 3, "IsTab": false},
					},
				},
				{
					"Name":  "空调",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "壁挂式空调", "Level": 3, "IsTab": false},
						{"Name": "柜式空调", "Level": 3, "IsTab": false},
					},
				},
				{
					"Name":  "洗衣机",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "滚筒洗衣机", "Level": 3, "IsTab": false},
						{"Name": "波轮洗衣机", "Level": 3, "IsTab": false},
					},
				},
			},
		},
		{
			"Name":  "手机数码",
			"Level": 1,
			"IsTab": true,
			"SubCategories": []map[string]interface{}{
				{
					"Name":  "手机",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "智能手机", "Level": 3, "IsTab": false},
						{"Name": "老人机", "Level": 3, "IsTab": false},
					},
				},
				{
					"Name":  "平板",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "安卓平板", "Level": 3, "IsTab": false},
						{"Name": "苹果iPad", "Level": 3, "IsTab": false},
					},
				},
				{
					"Name":  "智能穿戴",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "智能手表", "Level": 3, "IsTab": false},
						{"Name": "智能手环", "Level": 3, "IsTab": false},
					},
				},
			},
		},
		{
			"Name":  "食品酒水",
			"Level": 1,
			"IsTab": true,
			"SubCategories": []map[string]interface{}{
				{
					"Name":  "休闲零食",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "坚果炒货", "Level": 3, "IsTab": false},
						{"Name": "糖果巧克力", "Level": 3, "IsTab": false},
					},
				},
				{
					"Name":  "粮油调味",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "食用油", "Level": 3, "IsTab": false},
						{"Name": "大米", "Level": 3, "IsTab": false},
					},
				},
				{
					"Name":  "饮料冲调",
					"Level": 2,
					"IsTab": false,
					"SubCategories": []map[string]interface{}{
						{"Name": "牛奶乳品", "Level": 3, "IsTab": false},
						{"Name": "咖啡", "Level": 3, "IsTab": false},
					},
				},
			},
		},
	}
}
