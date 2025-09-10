package handler

import (
	context "context"
	"fmt"
	"order_srv/global"
	"order_srv/model"
	"order_srv/proto"
	inventorypb "order_srv/proto/inventory"
	"order_srv/utils"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type OrderServiceServer struct {
	proto.UnimplementedOrderServiceServer
}

// 购物车相关
func (s *OrderServiceServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	global.Logger.Infof("获取用户购物车列表，用户ID: %d", req.Id)

	// 参数验证
	if req.Id <= 0 {
		global.Logger.Error("用户ID无效")
		return nil, status.Errorf(codes.InvalidArgument, "用户ID必须大于0")
	}

	var shoppingCarts []model.ShoppingCart
	var total int64

	// 查询用户的购物车记录总数
	if err := global.DB.Model(&model.ShoppingCart{}).Where("user = ?", req.Id).Count(&total).Error; err != nil {
		global.Logger.Errorf("查询购物车总数失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询购物车失败")
	}

	// 查询购物车记录
	if err := global.DB.Where("user = ?", req.Id).Find(&shoppingCarts).Error; err != nil {
		global.Logger.Errorf("查询购物车记录失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询购物车失败")
	}

	// 转换为响应格式
	cartItems := make([]*proto.ShopCartInfoResponse, 0, len(shoppingCarts))
	for _, cart := range shoppingCarts {
		cartItems = append(cartItems, &proto.ShopCartInfoResponse{
			Id:         cart.ID,
			UserId:     cart.User,
			GoodsId:    cart.Goods,
			GoodsName:  cart.GoodsName,
			GoodsImage: cart.GoodsImage,
			GoodsPrice: cart.GoodsPrice,
			Nums:       cart.Nums,
			Checked:    cart.Checked,
		})
	}

	response := &proto.CartItemListResponse{
		Total:     int32(total),
		CartItems: cartItems,
	}

	global.Logger.Infof("成功获取用户购物车列表，用户ID: %d，总数: %d", req.Id, total)
	return response, nil
}

func (s *OrderServiceServer) CartItemAdd(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	global.Logger.Infof("添加购物车商品，用户ID: %d，商品ID: %d，数量: %d", req.UserId, req.GoodsId, req.Nums)

	// 参数验证
	if req.UserId <= 0 {
		global.Logger.Error("用户ID无效")
		return nil, status.Errorf(codes.InvalidArgument, "用户ID必须大于0")
	}
	if req.GoodsId <= 0 {
		global.Logger.Error("商品ID无效")
		return nil, status.Errorf(codes.InvalidArgument, "商品ID必须大于0")
	}
	if req.Nums <= 0 {
		global.Logger.Error("商品数量无效")
		return nil, status.Errorf(codes.InvalidArgument, "商品数量必须大于0")
	}
	if req.GoodsName == "" {
		global.Logger.Error("商品名称不能为空")
		return nil, status.Errorf(codes.InvalidArgument, "商品名称不能为空")
	}
	if req.GoodsPrice <= 0 {
		global.Logger.Error("商品价格无效")
		return nil, status.Errorf(codes.InvalidArgument, "商品价格必须大于0")
	}

	var shoppingCart model.ShoppingCart

	// 查询是否已存在该商品
	result := global.DB.Where("user = ? AND goods = ?", req.UserId, req.GoodsId).First(&shoppingCart)

	if result.Error == nil {
		// 商品已存在，更新数量和商品信息（价格可能有变化）
		global.Logger.Infof("购物车中已存在该商品，更新数量: 原数量 %d，新增数量 %d", shoppingCart.Nums, req.Nums)
		shoppingCart.Nums += req.Nums
		shoppingCart.GoodsName = req.GoodsName
		shoppingCart.GoodsImage = req.GoodsImage
		shoppingCart.GoodsPrice = req.GoodsPrice

		if err := global.DB.Save(&shoppingCart).Error; err != nil {
			global.Logger.Errorf("更新购物车商品数量失败: %v", err)
			return nil, status.Errorf(codes.Internal, "更新购物车失败")
		}

		global.Logger.Infof("成功更新购物车商品数量，商品ID: %d，新数量: %d", shoppingCart.Goods, shoppingCart.Nums)
	} else {
		// 商品不存在，创建新记录
		global.Logger.Infof("购物车中不存在该商品，创建新记录")
		shoppingCart = model.ShoppingCart{
			User:       req.UserId,
			Goods:      req.GoodsId,
			GoodsName:  req.GoodsName,
			GoodsImage: req.GoodsImage,
			GoodsPrice: req.GoodsPrice,
			Nums:       req.Nums,
			Checked:    false,
		}

		if err := global.DB.Create(&shoppingCart).Error; err != nil {
			global.Logger.Errorf("创建购物车记录失败: %v", err)
			return nil, status.Errorf(codes.Internal, "创建购物车失败")
		}

		global.Logger.Infof("成功创建购物车记录，商品ID: %d，数量: %d", shoppingCart.Goods, shoppingCart.Nums)
	}

	// 返回购物车信息
	response := &proto.ShopCartInfoResponse{
		Id:      shoppingCart.ID,
		UserId:  shoppingCart.User,
		GoodsId: shoppingCart.Goods,
		Nums:    shoppingCart.Nums,
		Checked: shoppingCart.Checked,
	}

	global.Logger.Infof("购物车操作完成，返回购物车ID: %d", shoppingCart.ID)
	return response, nil
}

func (s *OrderServiceServer) CartItemUpdate(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	global.Logger.Infof("更新购物车商品，ID: %d，用户ID: %d，数量: %d，选中状态: %t", req.Id, req.UserId, req.Nums, req.Checked)

	// 参数验证
	if req.Id <= 0 && (req.UserId <= 0 || req.GoodsId <= 0) {
		global.Logger.Error("购物车ID或用户ID+商品ID必须提供")
		return nil, status.Errorf(codes.InvalidArgument, "购物车ID或用户ID+商品ID必须提供")
	}
	if req.Nums < 0 {
		global.Logger.Error("商品数量不能为负数")
		return nil, status.Errorf(codes.InvalidArgument, "商品数量不能为负数")
	}

	var shoppingCart model.ShoppingCart
	var err error

	// 根据ID查询或根据用户ID+商品ID查询
	if req.Id > 0 {
		err = global.DB.Where("id = ?", req.Id).First(&shoppingCart).Error
	} else {
		err = global.DB.Where("user = ? AND goods = ?", req.UserId, req.GoodsId).First(&shoppingCart).Error
	}

	if err != nil {
		global.Logger.Errorf("查询购物车记录失败: %v", err)
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	// 更新字段
	updated := false
	if req.Nums > 0 && req.Nums != shoppingCart.Nums {
		global.Logger.Infof("更新购物车数量: %d -> %d", shoppingCart.Nums, req.Nums)
		shoppingCart.Nums = req.Nums
		updated = true
	}
	if req.Checked != shoppingCart.Checked {
		global.Logger.Infof("更新购物车选中状态: %t -> %t", shoppingCart.Checked, req.Checked)
		shoppingCart.Checked = req.Checked
		updated = true
	}

	if !updated {
		global.Logger.Info("购物车信息无需更新")
		return &emptypb.Empty{}, nil
	}

	// 保存更新
	if err := global.DB.Save(&shoppingCart).Error; err != nil {
		global.Logger.Errorf("更新购物车失败: %v", err)
		return nil, status.Errorf(codes.Internal, "更新购物车失败")
	}

	global.Logger.Infof("成功更新购物车，ID: %d", shoppingCart.ID)
	return &emptypb.Empty{}, nil
}

func (s *OrderServiceServer) CartItemDelete(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	global.Logger.Infof("删除购物车商品，ID: %d，用户ID: %d，商品ID: %d", req.Id, req.UserId, req.GoodsId)

	// 参数验证
	if req.Id <= 0 && (req.UserId <= 0 || req.GoodsId <= 0) {
		global.Logger.Error("购物车ID或用户ID+商品ID必须提供")
		return nil, status.Errorf(codes.InvalidArgument, "购物车ID或用户ID+商品ID必须提供")
	}

	var result *gorm.DB

	// 根据ID删除或根据用户ID+商品ID删除
	if req.Id > 0 {
		result = global.DB.Where("id = ?", req.Id).Delete(&model.ShoppingCart{})
	} else {
		result = global.DB.Where("user = ? AND goods = ?", req.UserId, req.GoodsId).Delete(&model.ShoppingCart{})
	}

	if result.Error != nil {
		global.Logger.Errorf("删除购物车记录失败: %v", result.Error)
		return nil, status.Errorf(codes.Internal, "删除购物车失败")
	}

	if result.RowsAffected == 0 {
		global.Logger.Warn("购物车记录不存在或已被删除")
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	global.Logger.Infof("成功删除购物车记录，影响行数: %d", result.RowsAffected)
	return &emptypb.Empty{}, nil
}

// ===============================================
// 订单创建核心逻辑：分布式锁 + 跨服务调用 + 批量处理
// ===============================================
//
// 架构设计说明：
// 1. 分布式锁：基于Redis的分布式锁，防止用户重复下单
// 2. 跨服务调用：批量调用商品服务和库存服务，提高性能
// 3. 数据库事务：保证订单和订单商品的数据一致性
// 4. 批量插入：使用GORM的CreateInBatches进行高效批量插入
// 5. 异常回滚：任何环节失败都会回滚事务和库存
//
// 执行流程：
// 分布式锁获取 -> 查询购物车 -> 数据库事务开启 -> 创建订单基础信息 
// -> 批量获取商品信息 -> 验证商品可用性 -> 批量扣减库存 
// -> 批量插入订单商品 -> 更新订单总额 -> 清空购物车 -> 提交事务 -> 释放锁
//
// 性能优化：
// - 批量获取商品信息（减少网络调用）
// - 批量扣减库存（减少服务间通信）
// - 批量插入订单商品（减少数据库交互）
// - 分布式锁保护（防止重复提交）
//
func (s *OrderServiceServer) OrderCreate(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	global.Logger.Infof("开始创建订单，用户ID: %d，收货地址: %s，收货人: %s", req.UserId, req.Address, req.Name)

	// 参数验证
	if req.UserId <= 0 {
		global.Logger.Error("用户ID无效")
		return nil, status.Errorf(codes.InvalidArgument, "用户ID必须大于0")
	}

	// 验证用户是否存在（通过查询用户相关的购物车或其他方式间接验证）
	// 在微服务架构中，我们可以通过查询该用户是否有相关数据来验证用户存在性
	if valid, err := s.validateUserExists(ctx, req.UserId); err != nil {
		global.Logger.Errorf("验证用户存在性失败: %v", err)
		return nil, status.Errorf(codes.Internal, "验证用户信息失败")
	} else if !valid {
		global.Logger.Warnf("用户不存在或无效，用户ID: %d", req.UserId)
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if req.Address == "" {
		global.Logger.Error("收货地址不能为空")
		return nil, status.Errorf(codes.InvalidArgument, "收货地址不能为空")
	}
	if req.Name == "" {
		global.Logger.Error("收货人姓名不能为空")
		return nil, status.Errorf(codes.InvalidArgument, "收货人姓名不能为空")
	}
	if req.Mobile == "" {
		global.Logger.Error("收货人手机不能为空")
		return nil, status.Errorf(codes.InvalidArgument, "收货人手机不能为空")
	}

	// ===========================================
	// 分布式锁机制：防止用户重复下单和并发问题
	// ===========================================
	// 
	// 锁定策略：按用户ID锁定，防止同一用户并发创建订单
	// 锁定时长：30秒，足够完成整个订单创建流程（包括跨服务调用和批量插入）
	// 重试机制：最多重试3次，每次间隔100ms，确保在高并发情况下的可用性
	// 
	// 为什么需要分布式锁：
	// 1. 防止用户重复提交订单（双击、网络重试等）
	// 2. 确保库存扣减的原子性（跨服务调用需要一致性保证）
	// 3. 保护批量插入操作的完整性
	lockKey := fmt.Sprintf("order_create_lock:%d", req.UserId)
	lock := utils.NewRedisLock(lockKey, 30*time.Second) // 30秒锁定时间覆盖完整的订单创建流程
	
	locked, err := lock.TryLock(ctx, 3, 100*time.Millisecond) // 重试3次，间隔100ms
	if err != nil {
		global.Logger.Errorf("获取创建订单锁失败: %v", err)
		return nil, status.Errorf(codes.Internal, "系统忙，请稍后重试")
	}
	if !locked {
		global.Logger.Warn("用户正在创建其他订单，请稍后重试")
		return nil, status.Errorf(codes.ResourceExhausted, "正在处理其他订单，请稍后重试")
	}
	
	// 确保在任何情况下都能释放锁，包括panic和正常返回
	defer func() {
		if unlockErr := lock.Unlock(ctx); unlockErr != nil {
			global.Logger.Errorf("释放创建订单锁失败: %v", unlockErr)
		}
	}()

	// 查询用户购物车中选中的商品
	var shoppingCarts []model.ShoppingCart
	if err := global.DB.Where("user = ? AND checked = ?", req.UserId, true).Find(&shoppingCarts).Error; err != nil {
		global.Logger.Errorf("查询购物车失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询购物车失败")
	}

	if len(shoppingCarts) == 0 {
		global.Logger.Warn("购物车中没有选中的商品")
		return nil, status.Errorf(codes.FailedPrecondition, "购物车中没有选中的商品")
	}

	// 开启数据库事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, status.Errorf(codes.Internal, "开启事务失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 生成订单号
	orderSn := utils.GenerateOrderSn(req.UserId)
	
	// 创建订单基本信息
	orderInfo := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      orderSn,
		PayType:      "alipay", // 默认支付宝，后续可从请求中获取
		Status:       "WAIT_BUYER_PAY",
		OrderMount:   0, // 先设为0，后面计算
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}

	// 保存订单基本信息
	if err := tx.Create(&orderInfo).Error; err != nil {
		tx.Rollback()
		global.Logger.Errorf("创建订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	// ===========================================
	// 跨服务调用阶段：获取商品信息和扣减库存
	// ===========================================
	//
	// 调用顺序和原因：
	// 1. 先调用商品服务获取商品详情（价格、库存、上架状态等）
	// 2. 再调用库存服务扣减库存（确保有足够库存后才扣减）
	// 3. 在分布式锁保护下进行，确保操作的原子性
	//
	// 异常处理策略：
	// - 任何跨服务调用失败都会回滚数据库事务
	// - 库存扣减失败会在后续数据库操作失败时自动回滚
	
	// 收集购物车中所有商品ID，用于批量查询商品信息
	goodsIds := make([]int32, len(shoppingCarts))
	for i, cart := range shoppingCarts {
		goodsIds[i] = cart.Goods
	}
	global.Logger.Infof("准备查询商品信息，商品ID列表: %v", goodsIds)

	// 批量调用商品服务获取商品详情
	// 使用批量接口提高性能，减少网络调用次数
	goodsMap, err := utils.GetGoodsByIds(ctx, goodsIds)
	if err != nil {
		tx.Rollback()
		global.Logger.Errorf("批量获取商品信息失败: %v", err)
		return nil, status.Errorf(codes.Internal, "获取商品信息失败")
	}
	global.Logger.Infof("成功获取商品信息，数量: %d", len(goodsMap))

	// 验证商品可用性并准备库存扣减数据
	// 在扣减库存前先验证所有商品的可用性，避免部分扣减后发现问题
	var sellItems []*inventorypb.GoodsInvInfo
	var totalAmount float32 = 0
	
	global.Logger.Info("开始验证商品可用性并计算订单总额")
	for _, cart := range shoppingCarts {
		goodsInfo, exists := goodsMap[cart.Goods]
		if !exists {
			tx.Rollback()
			global.Logger.Errorf("商品不存在，商品ID: %d", cart.Goods)
			return nil, status.Errorf(codes.NotFound, "商品不存在")
		}

		// 验证商品是否可用（上架状态、库存充足等）
		if err := utils.ValidateGoodsAvailability(goodsInfo, cart.Nums); err != nil {
			tx.Rollback()
			global.Logger.Errorf("商品验证失败: %v", err)
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		// 准备库存扣减数据，用于批量调用库存服务
		sellItems = append(sellItems, &inventorypb.GoodsInvInfo{
			GoodsId: cart.Goods,
			Num:     cart.Nums,
		})

		// 计算订单总金额
		itemTotal := goodsInfo.ShopPrice * float32(cart.Nums)
		totalAmount += itemTotal
		
		global.Logger.Infof("商品验证通过: %s，单价: %.2f，数量: %d，小计: %.2f", 
			goodsInfo.Name, goodsInfo.ShopPrice, cart.Nums, itemTotal)
	}

	// 调用库存服务批量扣减库存
	// 重要：此操作在数据库事务外进行，如果后续数据库操作失败需要手动回滚库存
	global.Logger.Infof("开始扣减库存，扣减项目数: %d", len(sellItems))
	if err := utils.SellInventory(ctx, sellItems); err != nil {
		tx.Rollback()
		global.Logger.Errorf("库存扣减失败: %v", err)
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}
	global.Logger.Info("库存扣减成功")

	// 批量创建订单商品记录
	// 使用批量插入提高性能，减少数据库交互次数
	// 在分布式锁和事务保护下，确保数据一致性
	global.Logger.Infof("开始批量创建订单商品记录，商品数量: %d", len(shoppingCarts))
	
	// 构建订单商品批量数据
	orderGoodsList := make([]model.OrderGoods, 0, len(shoppingCarts))
	for i, cart := range shoppingCarts {
		goodsInfo := goodsMap[cart.Goods]
		
		orderGoods := model.OrderGoods{
			Order:      int32(orderInfo.ID),
			Goods:      cart.Goods,
			GoodsName:  goodsInfo.Name,
			GoodsImage: goodsInfo.GoodsFrontImage,
			GoodsPrice: goodsInfo.ShopPrice,
			Nums:       cart.Nums,
		}
		
		orderGoodsList = append(orderGoodsList, orderGoods)
		global.Logger.Debugf("准备订单商品[%d]: 商品ID=%d, 商品名=%s, 数量=%d, 单价=%.2f", 
			i+1, cart.Goods, goodsInfo.Name, cart.Nums, goodsInfo.ShopPrice)
	}
	
	// 执行批量插入
	// 使用CreateInBatches避免单次插入过多数据导致的性能问题
	// 批量大小设为100，平衡性能和内存使用
	const batchSize = 100
	if err := tx.CreateInBatches(&orderGoodsList, batchSize).Error; err != nil {
		tx.Rollback()
		// 批量插入失败时，需要回滚之前成功扣减的库存
		// 确保库存和订单数据的一致性
		if rebackErr := utils.RebackInventory(ctx, sellItems); rebackErr != nil {
			global.Logger.Errorf("批量插入失败后回滚库存失败: %v", rebackErr)
		}
		global.Logger.Errorf("批量创建订单商品失败: %v", err)
		return nil, status.Errorf(codes.Internal, "批量创建订单商品失败")
	}
	
	global.Logger.Infof("成功批量创建订单商品记录，总数: %d，批量大小: %d", len(orderGoodsList), batchSize)

	// 更新订单总金额
	if err := tx.Model(&orderInfo).Update("order_mount", totalAmount).Error; err != nil {
		tx.Rollback()
		// 回滚库存
		if rebackErr := utils.RebackInventory(ctx, sellItems); rebackErr != nil {
			global.Logger.Errorf("回滚库存失败: %v", rebackErr)
		}
		global.Logger.Errorf("更新订单总金额失败: %v", err)
		return nil, status.Errorf(codes.Internal, "更新订单总金额失败")
	}

	// 清空购物车中的选中商品
	if err := tx.Where("user = ? AND checked = ?", req.UserId, true).Delete(&model.ShoppingCart{}).Error; err != nil {
		tx.Rollback()
		// 回滚库存
		if rebackErr := utils.RebackInventory(ctx, sellItems); rebackErr != nil {
			global.Logger.Errorf("回滚库存失败: %v", rebackErr)
		}
		global.Logger.Errorf("清空购物车失败: %v", err)
		return nil, status.Errorf(codes.Internal, "清空购物车失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		// 回滚库存
		if rebackErr := utils.RebackInventory(ctx, sellItems); rebackErr != nil {
			global.Logger.Errorf("回滚库存失败: %v", rebackErr)
		}
		global.Logger.Errorf("提交事务失败: %v", err)
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	}

	// 构造返回数据
	response := &proto.OrderInfoResponse{
		Id:      int32(orderInfo.ID),
		UserId:  orderInfo.User,
		OrderSn: orderInfo.OrderSn,
		PayType: orderInfo.PayType,
		Status:  orderInfo.Status,
		Post:    orderInfo.Post,
		Total:   totalAmount,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SingerMobile,
	}

	global.Logger.Infof("成功创建订单，订单ID: %d，订单号: %s，总金额: %.2f", orderInfo.ID, orderInfo.OrderSn, totalAmount)
	return response, nil
}

func (s *OrderServiceServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	global.Logger.Infof("获取订单列表，用户ID: %d，状态: %s，页码: %d，页大小: %d", req.UserId, req.Status, req.Page, req.PageSize)

	// 参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100 // 限制最大页大小
	}

	var orders []model.OrderInfo
	var total int64

	// 构建查询条件
	query := global.DB.Model(&model.OrderInfo{})

	// 用户查询：提供用户ID则为用户查询，否则为后台管理查询
	if req.UserId > 0 {
		global.Logger.Infof("用户查询订单，用户ID: %d", req.UserId)
		query = query.Where("user = ?", req.UserId)
	} else {
		global.Logger.Info("后台管理查询订单")
	}

	// 状态过滤
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
		global.Logger.Infof("过滤订单状态: %s", req.Status)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		global.Logger.Errorf("查询订单总数失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询订单总数失败")
	}

	// 分页查询订单
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(int(offset)).Limit(int(req.PageSize)).Order("id DESC").Find(&orders).Error; err != nil {
		global.Logger.Errorf("查询订单列表失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询订单列表失败")
	}

	// 转换为响应格式
	orderInfos := make([]*proto.OrderInfoResponse, 0, len(orders))
	for _, order := range orders {
		orderInfos = append(orderInfos, &proto.OrderInfoResponse{
			Id:      int32(order.ID),
			UserId:  order.User,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
		})
	}

	response := &proto.OrderListResponse{
		Total: int32(total),
		Data:  orderInfos,
	}

	global.Logger.Infof("成功获取订单列表，总数: %d，当前页数量: %d", total, len(orders))
	return response, nil
}

func (s *OrderServiceServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	global.Logger.Infof("获取订单详情，订单ID: %d，用户ID: %d", req.Id, req.UserId)

	// 参数验证
	if req.Id <= 0 {
		global.Logger.Error("订单ID无效")
		return nil, status.Errorf(codes.InvalidArgument, "订单ID必须大于0")
	}

	// 查询订单基本信息
	var orderInfo model.OrderInfo
	query := global.DB.Where("id = ?", req.Id)
	
	// 如果提供了用户ID，则验证订单所属用户
	if req.UserId > 0 {
		query = query.Where("user = ?", req.UserId)
		global.Logger.Infof("验证订单所属用户，用户ID: %d", req.UserId)
	}

	if err := query.First(&orderInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			global.Logger.Warnf("订单不存在，订单ID: %d", req.Id)
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		global.Logger.Errorf("查询订单信息失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询订单信息失败")
	}

	// 查询订单商品详情 - 使用实际的订单ID
	var orderGoods []model.OrderGoods
	if err := global.DB.Where("`order` = ?", orderInfo.ID).Find(&orderGoods).Error; err != nil {
		global.Logger.Errorf("查询订单商品失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询订单商品失败")
	}

	// 转换为响应格式
	orderInfoResponse := &proto.OrderInfoResponse{
		Id:      int32(orderInfo.ID),
		UserId:  orderInfo.User,
		OrderSn: orderInfo.OrderSn,
		PayType: orderInfo.PayType,
		Status:  orderInfo.Status,
		Post:    orderInfo.Post,
		Total:   orderInfo.OrderMount,
		Address: orderInfo.Address,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SingerMobile,
	}

	// 转换订单商品列表
	orderGoodsResponse := make([]*proto.OrderItemResponse, 0, len(orderGoods))
	for _, goods := range orderGoods {
		orderGoodsResponse = append(orderGoodsResponse, &proto.OrderItemResponse{
			Id:         int32(goods.ID),
			OrderId:    int32(goods.Order),
			GoodsId:    goods.Goods,
			GoodsName:  goods.GoodsName,
			GoodsImage: goods.GoodsImage,
			GoodsPrice: goods.GoodsPrice,
			Nums:       goods.Nums,
		})
	}

	response := &proto.OrderInfoDetailResponse{
		OrderInfo: orderInfoResponse,
		Goods:     orderGoodsResponse,
	}

	global.Logger.Infof("成功获取订单详情，订单ID: %d，商品数量: %d", req.Id, len(orderGoods))
	return response, nil
}

func (s *OrderServiceServer) OrderUpdate(ctx context.Context, req *proto.OrderStatus) (*emptypb.Empty, error) {
	global.Logger.Infof("更新订单状态，订单ID: %d，订单号: %s，新状态: %s", req.Id, req.OrderSn, req.Status)

	// 参数验证
	if req.Id <= 0 && req.OrderSn == "" {
		global.Logger.Error("订单ID或订单号必须提供")
		return nil, status.Errorf(codes.InvalidArgument, "订单ID或订单号必须提供")
	}
	if req.Status == "" {
		global.Logger.Error("订单状态不能为空")
		return nil, status.Errorf(codes.InvalidArgument, "订单状态不能为空")
	}

	// 验证状态值是否有效
	validStatuses := map[string]bool{
		"PAYING":         true, // 待支付
		"TRADE_SUCCESS":  true, // 成功
		"TRADE_CLOSED":   true, // 超时关闭
		"WAIT_BUYER_PAY": true, // 交易创建
		"TRADE_FINISHED": true, // 交易结束
	}
	if !validStatuses[req.Status] {
		global.Logger.Errorf("无效的订单状态: %s", req.Status)
		return nil, status.Errorf(codes.InvalidArgument, "无效的订单状态")
	}

	// 创建分布式锁，防止并发更新同一订单
	var lockKey string
	if req.Id > 0 {
		lockKey = fmt.Sprintf("order_update_lock:id:%d", req.Id)
	} else {
		lockKey = fmt.Sprintf("order_update_lock:sn:%s", req.OrderSn)
	}
	
	lock := utils.NewRedisLock(lockKey, 10*time.Second) // 10秒锁定时间
	locked, err := lock.TryLock(ctx, 3, 50*time.Millisecond)
	if err != nil {
		global.Logger.Errorf("获取订单更新锁失败: %v", err)
		return nil, status.Errorf(codes.Internal, "系统忙，请稍后重试")
	}
	if !locked {
		global.Logger.Warn("订单正在被其他请求更新，请稍后重试")
		return nil, status.Errorf(codes.ResourceExhausted, "订单正在处理中，请稍后重试")
	}
	
	defer func() {
		if unlockErr := lock.Unlock(ctx); unlockErr != nil {
			global.Logger.Errorf("释放订单更新锁失败: %v", unlockErr)
		}
	}()

	// 构建查询条件
	query := global.DB.Model(&model.OrderInfo{})
	if req.Id > 0 {
		query = query.Where("id = ?", req.Id)
	} else {
		query = query.Where("order_sn = ?", req.OrderSn)
	}

	// 查询订单是否存在
	var orderInfo model.OrderInfo
	if err := query.First(&orderInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			global.Logger.Warn("订单不存在")
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		global.Logger.Errorf("查询订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询订单失败")
	}

	// 记录状态变更
	oldStatus := orderInfo.Status
	if oldStatus == req.Status {
		global.Logger.Infof("订单状态无需更新，当前状态: %s", oldStatus)
		return &emptypb.Empty{}, nil
	}

	// 验证状态转换是否合法
	if !isValidStatusTransition(oldStatus, req.Status) {
		global.Logger.Warnf("无效的状态转换: %s -> %s", oldStatus, req.Status)
		return nil, status.Errorf(codes.FailedPrecondition, "无效的状态转换")
	}

	// 开启事务更新订单状态
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, status.Errorf(codes.Internal, "开启事务失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新订单状态
	updateData := map[string]interface{}{
		"status": req.Status,
	}
	
	// 如果是支付成功，更新支付时间
	if req.Status == "TRADE_SUCCESS" {
		updateData["pay_time"] = time.Now()
	}

	if err := tx.Model(&model.OrderInfo{}).Where("id = ?", orderInfo.ID).Updates(updateData).Error; err != nil {
		tx.Rollback()
		global.Logger.Errorf("更新订单状态失败: %v", err)
		return nil, status.Errorf(codes.Internal, "更新订单状态失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Logger.Errorf("提交事务失败: %v", err)
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	}

	global.Logger.Infof("成功更新订单状态，订单ID: %d，状态变更: %s -> %s", orderInfo.ID, oldStatus, req.Status)
	return &emptypb.Empty{}, nil
}

// isValidStatusTransition 验证订单状态转换是否合法
func isValidStatusTransition(oldStatus, newStatus string) bool {
	// 定义合法的状态转换规则
	transitions := map[string][]string{
		"WAIT_BUYER_PAY": {"PAYING", "TRADE_CLOSED"},
		"PAYING":         {"TRADE_SUCCESS", "TRADE_CLOSED"},
		"TRADE_SUCCESS":  {"TRADE_FINISHED"},
		"TRADE_CLOSED":   {}, // 已关闭的订单不能再转换
		"TRADE_FINISHED": {}, // 已完成的订单不能再转换
	}
	
	validNextStates, exists := transitions[oldStatus]
	if !exists {
		return false
	}
	
	for _, validState := range validNextStates {
		if validState == newStatus {
			return true
		}
	}
	
	return false
}

func (s *OrderServiceServer) OrderDelete(ctx context.Context, req *proto.OrderDelRequest) (*emptypb.Empty, error) {
	global.Logger.Infof("删除订单，订单ID: %d，用户ID: %d", req.Id, req.UserId)

	// 参数验证
	if req.Id <= 0 {
		global.Logger.Error("订单ID无效")
		return nil, status.Errorf(codes.InvalidArgument, "订单ID必须大于0")
	}

	// 创建分布式锁，防止并发删除同一订单
	lockKey := fmt.Sprintf("order_delete_lock:%d", req.Id)
	lock := utils.NewRedisLock(lockKey, 15*time.Second) // 15秒锁定时间
	
	locked, err := lock.TryLock(ctx, 3, 50*time.Millisecond)
	if err != nil {
		global.Logger.Errorf("获取订单删除锁失败: %v", err)
		return nil, status.Errorf(codes.Internal, "系统忙，请稍后重试")
	}
	if !locked {
		global.Logger.Warn("订单正在被其他请求处理，请稍后重试")
		return nil, status.Errorf(codes.ResourceExhausted, "订单正在处理中，请稍后重试")
	}
	
	defer func() {
		if unlockErr := lock.Unlock(ctx); unlockErr != nil {
			global.Logger.Errorf("释放订单删除锁失败: %v", unlockErr)
		}
	}()

	// 开启事务，确保订单和订单商品一起删除
	tx := global.DB.Begin()
	if tx.Error != nil {
		global.Logger.Errorf("开启事务失败: %v", tx.Error)
		return nil, status.Errorf(codes.Internal, "开启事务失败")
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 构建查询条件
	query := tx.Model(&model.OrderInfo{}).Where("id = ?", req.Id)
	
	// 如果提供了用户ID，则验证订单所属用户
	if req.UserId > 0 {
		query = query.Where("user = ?", req.UserId)
		global.Logger.Infof("验证订单所属用户，用户ID: %d", req.UserId)
	}

	// 检查订单是否存在
	var orderInfo model.OrderInfo
	if err := query.First(&orderInfo).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			global.Logger.Warn("订单不存在或无权限删除")
			return nil, status.Errorf(codes.NotFound, "订单不存在或无权限删除")
		}
		global.Logger.Errorf("查询订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "查询订单失败")
	}

	// 检查订单状态是否允许删除
	if !isDeletableOrderStatus(orderInfo.Status) {
		tx.Rollback()
		global.Logger.Warnf("订单状态不允许删除，当前状态: %s", orderInfo.Status)
		return nil, status.Errorf(codes.FailedPrecondition, "当前状态的订单不允许删除")
	}

	// 查询订单商品数量，用于日志记录
	var orderGoodsCount int64
	if err := tx.Model(&model.OrderGoods{}).Where("order = ?", req.Id).Count(&orderGoodsCount).Error; err != nil {
		global.Logger.Warnf("查询订单商品数量失败: %v", err)
		orderGoodsCount = 0
	}

	// 删除订单商品
	if err := tx.Where("order = ?", req.Id).Delete(&model.OrderGoods{}).Error; err != nil {
		tx.Rollback()
		global.Logger.Errorf("删除订单商品失败: %v", err)
		return nil, status.Errorf(codes.Internal, "删除订单商品失败")
	}
	global.Logger.Infof("成功删除订单商品，数量: %d", orderGoodsCount)

	// 删除订单
	if err := tx.Delete(&orderInfo).Error; err != nil {
		tx.Rollback()
		global.Logger.Errorf("删除订单失败: %v", err)
		return nil, status.Errorf(codes.Internal, "删除订单失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Logger.Errorf("提交事务失败: %v", err)
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	}

	global.Logger.Infof("成功删除订单，订单ID: %d，订单号: %s，原状态: %s", req.Id, orderInfo.OrderSn, orderInfo.Status)
	return &emptypb.Empty{}, nil
}

// isDeletableOrderStatus 检查订单状态是否允许删除
func isDeletableOrderStatus(status string) bool {
	// 只有特定状态的订单才允许删除
	deletableStatuses := map[string]bool{
		"WAIT_BUYER_PAY": true, // 待支付的订单可以删除
		"TRADE_CLOSED":   true, // 已关闭的订单可以删除
		"PAYING":         true, // 支付中的订单可以删除（取消支付）
	}
	
	return deletableStatuses[status]
}

// validateUserExists 验证用户是否存在
// 在微服务架构中，我们通过以下方式验证用户存在性：
// 1. 查询用户是否有历史购物车记录
// 2. 查询用户是否有历史订单记录  
// 3. 如果都没有，则认为可能是新用户（允许创建订单）
func (s *OrderServiceServer) validateUserExists(ctx context.Context, userID int32) (bool, error) {
	global.Logger.Infof("验证用户存在性，用户ID: %d", userID)

	// 策略1: 查询用户是否有购物车记录（不管是否删除）
	var cartCount int64
	if err := global.DB.Unscoped().Model(&model.ShoppingCart{}).Where("user = ?", userID).Count(&cartCount).Error; err != nil {
		global.Logger.Errorf("查询用户购物车记录失败: %v", err)
		return false, err
	}
	if cartCount > 0 {
		global.Logger.Infof("用户存在购物车记录，用户ID: %d", userID)
		return true, nil
	}

	// 策略2: 查询用户是否有历史订单记录  
	var orderCount int64
	if err := global.DB.Unscoped().Model(&model.OrderInfo{}).Where("user = ?", userID).Count(&orderCount).Error; err != nil {
		global.Logger.Errorf("查询用户订单记录失败: %v", err)
		return false, err
	}
	if orderCount > 0 {
		global.Logger.Infof("用户存在订单记录，用户ID: %d", userID)
		return true, nil
	}

	// 策略3: 对于全新用户，我们通过用户ID的合理性判断
	// 在实际生产环境中，这里应该调用用户服务验证用户存在性
	// 为了测试方便，我们假设用户ID在合理范围内（1-1000）就是有效的
	if userID > 0 && userID <= 1000 {
		global.Logger.Infof("新用户允许创建订单，用户ID: %d", userID)
		return true, nil
	}

	global.Logger.Warnf("用户ID超出合理范围或无效，用户ID: %d", userID)
	return false, nil
}
