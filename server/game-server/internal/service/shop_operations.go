package service

import (
	"context"
	"fmt"

	"github.com/cultivation-world/shared/errors"
	"github.com/cultivation-world/shared/types"
)

// ── 商店系统（F1） ──

// executeShopList 列出当前区域的所有商店
func (s *OperationService) executeShopList(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "商店系统不可用")
	}

	shops, err := s.shopRepo.ListShopsByRegion(ctx, entity.Position.RegionID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "查询商店失败")
	}

	shopList := make([]map[string]interface{}, 0, len(shops))
	for _, shop := range shops {
		shopList = append(shopList, map[string]interface{}{
			"id":          shop.ID,
			"name":        shop.Name,
			"description": shop.Description,
			"shop_type":   shop.ShopType,
			"npc_owner":   shop.NPCOwner,
		})
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("当前区域共有 %d 家商店", len(shopList)),
		Effects: map[string]interface{}{
			"shops": shopList,
			"count": len(shopList),
		},
	}, nil
}

// executeShopItems 查看商店库存
func (s *OperationService) executeShopItems(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "商店系统不可用")
	}

	shopID, ok := op.Params["shop_id"].(string)
	if !ok || shopID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少商店ID")
	}

	shop, err := s.shopRepo.GetShopByID(ctx, shopID)
	if err != nil || shop == nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "商店不存在")
	}

	items, err := s.shopRepo.GetShopInventory(ctx, shopID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "查询库存失败")
	}

	itemList := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		stock := "充足"
		if item.Quantity == 0 {
			stock = "缺货"
		} else if item.Quantity > 0 && item.Quantity <= 5 {
			stock = fmt.Sprintf("仅剩%d", item.Quantity)
		}

		itemList = append(itemList, map[string]interface{}{
			"name":       item.ItemName,
			"type":       item.ItemType,
			"rarity":     item.Rarity,
			"price":      item.Price,
			"stock":      stock,
			"min_realm":  item.MinRealm,
		})
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("%s 的商品列表：", shop.Name),
		Effects: map[string]interface{}{
			"shop_id":   shopID,
			"shop_name": shop.Name,
			"items":     itemList,
			"count":     len(itemList),
		},
	}, nil
}

// executeBuy 从 NPC 商店购买物品
func (s *OperationService) executeBuy(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "商店系统不可用")
	}

	shopID, ok := op.Params["shop_id"].(string)
	if !ok || shopID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少商店ID")
	}

	itemName, ok := op.Params["item_name"].(string)
	if !ok || itemName == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少物品名称")
	}

	quantity := 1
	if q, ok := op.Params["quantity"].(float64); ok && q > 0 {
		quantity = int(q)
	}

	shop, err := s.shopRepo.GetShopByID(ctx, shopID)
	if err != nil || shop == nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "商店不存在")
	}

	if shop.RegionID != "" && shop.RegionID != entity.Position.RegionID {
		return &types.OperationResult{
			Success: false,
			Message: "该商店不在此区域",
		}, nil
	}

	shopItem, err := s.shopRepo.GetShopItemByName(ctx, shopID, itemName)
	if err != nil || shopItem == nil {
		return &types.OperationResult{
			Success: false,
			Message: "该商店没有此物品",
		}, nil
	}

	if shopItem.MinRealm != "" && shopItem.MinRealm != "mortal" {
		if types.CultivationRealmLevel(entity.Realm) < types.CultivationRealmLevel(types.CultivationRealm(shopItem.MinRealm)) {
			return &types.OperationResult{
				Success: false,
				Message: fmt.Sprintf("境界不足，需要 %s", shopItem.MinRealm),
			}, nil
		}
	}

	if shopItem.Quantity >= 0 && shopItem.Quantity < quantity {
		return &types.OperationResult{
			Success: false,
			Message: "库存不足",
		}, nil
	}

	totalPrice := shopItem.Price * int64(quantity)

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	if attr.SpiritStones.LowGrade < totalPrice {
		return &types.OperationResult{
			Success: false,
			Message: fmt.Sprintf("灵石不足，需要 %d 低阶灵石", totalPrice),
		}, nil
	}

	attr.SpiritStones.LowGrade -= totalPrice

	for i := 0; i < quantity; i++ {
		s.ensureItemInInventory(ctx, entity.ID, shopItem.ItemName, types.ItemType(shopItem.ItemType), shopItem.Rarity)
	}

	if shopItem.Quantity > 0 {
		s.shopRepo.DecrementShopStock(ctx, shopID, shopItem.ItemName, quantity)
	}

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功购买 %s x%d，花费 %d 灵石", shopItem.ItemName, quantity, totalPrice),
		Effects: map[string]interface{}{
			"shop_id":    shopID,
			"item_name":  shopItem.ItemName,
			"quantity":   quantity,
			"price":      totalPrice,
			"unit_price": shopItem.Price,
		},
	}, nil
}

// executeSell 向 NPC 商店出售物品
func (s *OperationService) executeSell(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "商店系统不可用")
	}

	shopID, ok := op.Params["shop_id"].(string)
	if !ok || shopID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少商店ID")
	}

	itemName, ok := op.Params["item_name"].(string)
	if !ok || itemName == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少物品名称")
	}

	quantity := 1
	if q, ok := op.Params["quantity"].(float64); ok && q > 0 {
		quantity = int(q)
	}

	shop, err := s.shopRepo.GetShopByID(ctx, shopID)
	if err != nil || shop == nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "商店不存在")
	}

	item, err := s.itemRepo.GetByName(ctx, itemName)
	if err != nil || item == nil {
		return &types.OperationResult{
			Success: false,
			Message: "未知物品",
		}, nil
	}

	invItem, err := s.inventoryRepo.GetItem(ctx, entity.ID, item.ID)
	if err != nil || invItem == nil || invItem.Quantity < quantity {
		return &types.OperationResult{
			Success: false,
			Message: "背包中没有该物品",
		}, nil
	}

	if invItem.Equipped {
		return &types.OperationResult{
			Success: false,
			Message: "该物品已装备，请先卸下",
		}, nil
	}

	shopItem, _ := s.shopRepo.GetShopItemByName(ctx, shopID, itemName)
	var sellPrice int64
	if shopItem != nil {
		sellPrice = int64(float64(shopItem.Price) * shop.BuyRate)
	} else {
		basePrices := map[int]int64{1: 10, 2: 50, 3: 200, 4: 1000, 5: 5000}
		sellPrice = basePrices[item.Rarity]
		if sellPrice == 0 {
			sellPrice = 10
		}
	}

	totalPrice := sellPrice * int64(quantity)

	s.inventoryRepo.RemoveItem(ctx, entity.ID, item.ID, quantity)

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	attr.SpiritStones.LowGrade += totalPrice
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	s.modifyKarma(ctx, entity.ID, 1, "出售物品")

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功出售 %s x%d，获得 %d 灵石", itemName, quantity, totalPrice),
		Effects: map[string]interface{}{
			"shop_id":    shopID,
			"item_name":  itemName,
			"quantity":   quantity,
			"price":      totalPrice,
			"unit_price": sellPrice,
		},
	}, nil
}

// ── 拍卖行系统（F4） ──

// executeAuctionList 上架拍卖物品
func (s *OperationService) executeAuctionList(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "拍卖行系统不可用")
	}

	itemName, ok := op.Params["item_name"].(string)
	if !ok || itemName == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少物品名称")
	}

	price, ok := op.Params["price"].(float64)
	if !ok || price <= 0 {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "价格无效")
	}

	quantity := 1
	if q, ok := op.Params["quantity"].(float64); ok && q > 0 {
		quantity = int(q)
	}

	item, err := s.itemRepo.GetByName(ctx, itemName)
	if err != nil || item == nil {
		return &types.OperationResult{
			Success: false,
			Message: "物品不存在",
		}, nil
	}

	invItem, err := s.inventoryRepo.GetItem(ctx, entity.ID, item.ID)
	if err != nil || invItem == nil || invItem.Quantity < quantity {
		return &types.OperationResult{
			Success: false,
			Message: "背包中没有足够的该物品",
		}, nil
	}

	if invItem.Equipped {
		return &types.OperationResult{
			Success: false,
			Message: "该物品已装备，请先卸下",
		}, nil
	}

	if invItem.Bound {
		return &types.OperationResult{
			Success: false,
			Message: "已绑定的物品无法拍卖",
		}, nil
	}

	deposit := int64(price * 0.1)
	if deposit < 1 {
		deposit = 1
	}

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	if attr.SpiritStones.LowGrade < deposit {
		return &types.OperationResult{
			Success: false,
			Message: fmt.Sprintf("押金不足，需要 %d 灵石", deposit),
		}, nil
	}

	attr.SpiritStones.LowGrade -= deposit
	s.inventoryRepo.RemoveItem(ctx, entity.ID, item.ID, quantity)

	auctionID, err := s.shopRepo.CreateAuction(ctx, string(entity.ID), string(item.ID), item.Name, quantity, int64(price), deposit)
	if err != nil {
		s.inventoryRepo.AddItem(ctx, entity.ID, item.ID, quantity)
		attr.SpiritStones.LowGrade += deposit
		return nil, errors.NewGameError(errors.ErrInternalError, "创建拍卖失败")
	}

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.modifyKarma(ctx, entity.ID, 1, "上架拍卖")

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功上架 %s，价格 %d 灵石，押金 %d", item.Name, int64(price), deposit),
		Effects: map[string]interface{}{
			"auction_id": auctionID,
			"item_name":  item.Name,
			"price":      int64(price),
			"deposit":    deposit,
			"quantity":   quantity,
		},
	}, nil
}

// executeAuctionBuy 购买拍卖物品
func (s *OperationService) executeAuctionBuy(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "拍卖行系统不可用")
	}

	auctionID, ok := op.Params["auction_id"].(string)
	if !ok || auctionID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少拍卖ID")
	}

	auction, err := s.shopRepo.GetAuctionByID(ctx, auctionID)
	if err != nil || auction == nil {
		return &types.OperationResult{
			Success: false,
			Message: "拍卖不存在",
		}, nil
	}

	if auction.Status != "active" {
		return &types.OperationResult{
			Success: false,
			Message: "该拍卖已结束",
		}, nil
	}

	if auction.SellerID == string(entity.ID) {
		return &types.OperationResult{
			Success: false,
			Message: "不能购买自己的拍卖品",
		}, nil
	}

	buyerAttr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	if buyerAttr.SpiritStones.LowGrade < auction.Price {
		return &types.OperationResult{
			Success: false,
			Message: fmt.Sprintf("灵石不足，需要 %d", auction.Price),
		}, nil
	}

	seller, err := s.entityRepo.GetByID(ctx, types.EntityID(auction.SellerID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "卖家不存在")
	}

	buyerAttr.SpiritStones.LowGrade -= auction.Price
	sellerAttr, _ := s.entityRepo.GetAttributes(ctx, seller.ID)

	platformFee := int64(float64(auction.Price) * 0.05)
	sellerAttr.SpiritStones.LowGrade += auction.Price - platformFee

	// 返还卖家押金
	sellerAttr.SpiritStones.LowGrade += auction.Deposit

	s.entityRepo.UpdateAttributes(ctx, entity.ID, buyerAttr)
	s.entityRepo.UpdateAttributes(ctx, seller.ID, sellerAttr)

	s.shopRepo.BuyAuction(ctx, auctionID, string(entity.ID))
	s.inventoryRepo.AddItem(ctx, entity.ID, types.ItemID(auction.ItemID), auction.Quantity)

	s.modifyKarma(ctx, entity.ID, 2, "购买拍卖品")

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功购买 %s，花费 %d 灵石", auction.ItemName, auction.Price),
		Effects: map[string]interface{}{
			"auction_id":   auctionID,
			"item_name":    auction.ItemName,
			"price":        auction.Price,
			"platform_fee": platformFee,
			"seller_id":    auction.SellerID,
		},
	}, nil
}

// executeAuctionCancel 取消拍卖
func (s *OperationService) executeAuctionCancel(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "拍卖行系统不可用")
	}

	auctionID, ok := op.Params["auction_id"].(string)
	if !ok || auctionID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少拍卖ID")
	}

	auction, err := s.shopRepo.GetAuctionByID(ctx, auctionID)
	if err != nil || auction == nil {
		return &types.OperationResult{
			Success: false,
			Message: "拍卖不存在",
		}, nil
	}

	if auction.Status != "active" {
		return &types.OperationResult{
			Success: false,
			Message: "该拍卖已结束，无法取消",
		}, nil
	}

	if auction.SellerID != string(entity.ID) {
		return &types.OperationResult{
			Success: false,
			Message: "只能取消自己的拍卖",
		}, nil
	}

	err = s.shopRepo.CancelAuction(ctx, auctionID, string(entity.ID))
	if err != nil {
		return &types.OperationResult{
			Success: false,
			Message: "取消拍卖失败",
		}, nil
	}

	s.inventoryRepo.AddItem(ctx, entity.ID, types.ItemID(auction.ItemID), auction.Quantity)

	refundDeposit := int64(float64(auction.Deposit) * 0.8)
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	attr.SpiritStones.LowGrade += refundDeposit
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("已取消拍卖 %s，物品和押金（%d 灵石）已退还", auction.ItemName, refundDeposit),
		Effects: map[string]interface{}{
			"auction_id":     auctionID,
			"item_name":      auction.ItemName,
			"deposit_refund": refundDeposit,
		},
	}, nil
}


// executeAuctionView 查看活跃拍卖列表
func (s *OperationService) executeAuctionView(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.shopRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "拍卖行系统不可用")
	}

	auctions, err := s.shopRepo.ListActiveAuctions(ctx)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "查询拍卖列表失败")
	}

	auctionList := make([]map[string]interface{}, 0, len(auctions))
	for _, a := range auctions {
		auctionList = append(auctionList, map[string]interface{}{
			"id":        a.ID,
			"item_name": a.ItemName,
			"quantity":  a.Quantity,
			"price":     a.Price,
			"seller_id": a.SellerID,
			"status":    a.Status,
			"created_at": a.CreatedAt,
		})
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("当前有 %d 件拍卖品", len(auctionList)),
		Effects: map[string]interface{}{
			"auctions": auctionList,
			"count":    len(auctionList),
		},
	}, nil
}
