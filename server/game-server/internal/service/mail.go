package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cultivation-world/shared/types"
)

// MailRepository defines mail data access for the operation service.
type MailRepository interface {
	Create(ctx context.Context, mail *types.Mail) error
	GetByID(ctx context.Context, id string) (*types.Mail, error)
	GetByReceiver(ctx context.Context, receiverID types.EntityID, limit int) ([]*types.Mail, error)
	GetUnreadByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.Mail, error)
	GetUnclaimedByReceiver(ctx context.Context, receiverID types.EntityID) ([]*types.Mail, error)
	MarkAsRead(ctx context.Context, mailID string) error
	MarkAsClaimed(ctx context.Context, mailID string) error
	Delete(ctx context.Context, mailID string) error
}

// executeMailList 获取邮件列表
func (s *OperationService) executeMailList(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	mailType, _ := op.Params["mail_type"].(string)
	limit := 50

	var mails []*types.Mail
	var err error

	switch mailType {
	case "unread":
		mails, err = s.mailRepo.GetUnreadByReceiver(ctx, entity.ID)
	case "unclaimed":
		mails, err = s.mailRepo.GetUnclaimedByReceiver(ctx, entity.ID)
	default:
		mails, err = s.mailRepo.GetByReceiver(ctx, entity.ID, limit)
	}

	if err != nil {
		return nil, err
	}

	mailList := make([]map[string]interface{}, 0, len(mails))
	unreadCount := 0
	unclaimedCount := 0

	for _, m := range mails {
		entry := map[string]interface{}{
			"id":           m.ID,
			"sender_name":  m.SenderName,
			"title":        m.Title,
			"mail_type":    m.MailType,
			"is_read":      m.IsRead,
			"has_attachment": m.HasAttachment,
			"is_claimed":   m.IsClaimed,
			"created_at":   m.CreatedAt,
		}
		if !m.IsRead {
			unreadCount++
		}
		if m.HasAttachment && !m.IsClaimed {
			unclaimedCount++
		}
		mailList = append(mailList, entry)
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("邮箱共有 %d 封信件", len(mailList)),
		Effects: map[string]interface{}{
			"mails":          mailList,
			"count":          len(mailList),
			"unread_count":   unreadCount,
			"unclaimed_count": unclaimedCount,
		},
	}, nil
}

// executeMailRead 读取邮件详情
func (s *OperationService) executeMailRead(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	mailID, ok := op.Params["mail_id"].(string)
	if !ok || mailID == "" {
		return nil, fmt.Errorf("缺少邮件ID")
	}

	mail, err := s.mailRepo.GetByID(ctx, mailID)
	if err != nil {
		return nil, fmt.Errorf("邮件不存在")
	}
	if mail == nil {
		return &types.OperationResult{
			Success: false,
			Message: "邮件不存在",
		}, nil
	}

	if mail.ReceiverID != string(entity.ID) {
		return &types.OperationResult{
			Success: false,
			Message: "这不是您的邮件",
		}, nil
	}

	// 标记为已读
	if !mail.IsRead {
		s.mailRepo.MarkAsRead(ctx, mailID)
	}

	return &types.OperationResult{
		Success: true,
		Message: mail.Title,
		Effects: map[string]interface{}{
			"id":                  mail.ID,
			"sender_id":           mail.SenderID,
			"sender_name":         mail.SenderName,
			"title":               mail.Title,
			"content":             mail.Content,
			"mail_type":           mail.MailType,
			"has_attachment":      mail.HasAttachment,
			"attachment_item_name": mail.AttachmentItemName,
			"attachment_quantity": mail.AttachmentQuantity,
			"attachment_spirit_stones": mail.AttachmentSpiritStones,
			"is_claimed":          mail.IsClaimed,
			"created_at":          mail.CreatedAt,
		},
	}, nil
}

// executeMailClaim 领取邮件附件
func (s *OperationService) executeMailClaim(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	mailID, ok := op.Params["mail_id"].(string)
	if !ok || mailID == "" {
		return nil, fmt.Errorf("缺少邮件ID")
	}

	mail, err := s.mailRepo.GetByID(ctx, mailID)
	if err != nil {
		return nil, fmt.Errorf("邮件不存在")
	}
	if mail == nil {
		return &types.OperationResult{
			Success: false,
			Message: "邮件不存在",
		}, nil
	}

	if mail.ReceiverID != string(entity.ID) {
		return &types.OperationResult{
			Success: false,
			Message: "这不是您的邮件",
		}, nil
	}

	if !mail.HasAttachment {
		return &types.OperationResult{
			Success: false,
			Message: "此邮件没有附件",
		}, nil
	}

	if mail.IsClaimed {
		return &types.OperationResult{
			Success: false,
			Message: "附件已领取",
		}, nil
	}

	// 发放灵石
	if mail.AttachmentSpiritStones > 0 {
		attr, err := s.entityRepo.GetAttributes(ctx, entity.ID)
		if err != nil {
			return nil, err
		}
		attr.SpiritStones.LowGrade += mail.AttachmentSpiritStones
		s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	}

	// 发放物品
	if mail.AttachmentItemID != "" && mail.AttachmentQuantity > 0 {
		s.inventoryRepo.AddItem(ctx, entity.ID, types.ItemID(mail.AttachmentItemID), mail.AttachmentQuantity)
	}

	// 标记已领取
	s.mailRepo.MarkAsClaimed(ctx, mailID)

	return &types.OperationResult{
		Success: true,
		Message: "附件领取成功！",
		Effects: map[string]interface{}{
			"spirit_stones": mail.AttachmentSpiritStones,
			"item_name":     mail.AttachmentItemName,
			"item_quantity": mail.AttachmentQuantity,
		},
	}, nil
}

// executeMailDelete 删除邮件
func (s *OperationService) executeMailDelete(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	mailID, ok := op.Params["mail_id"].(string)
	if !ok || mailID == "" {
		return nil, fmt.Errorf("缺少邮件ID")
	}

	mail, err := s.mailRepo.GetByID(ctx, mailID)
	if err != nil || mail == nil {
		return &types.OperationResult{
			Success: false,
			Message: "邮件不存在",
		}, nil
	}

	if mail.ReceiverID != string(entity.ID) {
		return &types.OperationResult{
			Success: false,
			Message: "这不是您的邮件",
		}, nil
	}

	s.mailRepo.Delete(ctx, mailID)

	return &types.OperationResult{
		Success: true,
		Message: "邮件已删除",
		Effects: map[string]interface{}{
			"mail_id": mailID,
		},
	}, nil
}

// SendSystemMail 发送系统邮件（内部调用，突破奖励/活动奖励等）
func (s *OperationService) SendSystemMail(ctx context.Context, receiverID types.EntityID, title, content string, spiritStones int64, itemID string, itemName string, quantity int) error {
	mail := &types.Mail{
		ID:                     fmt.Sprintf("mail_%d", time.Now().UnixNano()),
		ReceiverID:             string(receiverID),
		SenderName:             "系统",
		Title:                  title,
		Content:                content,
		MailType:               "system",
		IsRead:                 false,
		HasAttachment:          spiritStones > 0 || (itemID != "" && quantity > 0),
		AttachmentItemID:       itemID,
		AttachmentItemName:     itemName,
		AttachmentQuantity:     quantity,
		AttachmentSpiritStones: spiritStones,
		IsClaimed:              false,
		CreatedAt:              time.Now().Unix(),
	}
	return s.mailRepo.Create(ctx, mail)
}

// SendRewardMail 发送奖励邮件（突破成功自动发放）
func (s *OperationService) SendRewardMail(ctx context.Context, receiverID types.EntityID, realm string) error {
	title := fmt.Sprintf("突破奖励 - 晋升%s", realm)
	content := fmt.Sprintf("恭喜突破至%s！系统发放突破奖励，请查收附件。", realm)
	spiritStones := int64(500)
	return s.SendSystemMail(ctx, receiverID, title, content, spiritStones, "", "", 0)
}
