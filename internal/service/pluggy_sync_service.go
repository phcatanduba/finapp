package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"finapp/internal/model"
	"finapp/internal/pluggy"
	"finapp/internal/repository"
)

type pluggySyncService struct {
	pluggyClient   *pluggy.Client
	itemRepo       repository.PluggyItemRepository
	accountRepo    repository.AccountRepository
	txRepo         repository.TransactionRepository
	webhookLogRepo repository.WebhookLogRepository
}

func NewPluggySyncService(
	pluggyClient *pluggy.Client,
	itemRepo repository.PluggyItemRepository,
	accountRepo repository.AccountRepository,
	txRepo repository.TransactionRepository,
	webhookLogRepo repository.WebhookLogRepository,
) PluggySyncService {
	return &pluggySyncService{
		pluggyClient:   pluggyClient,
		itemRepo:       itemRepo,
		accountRepo:    accountRepo,
		txRepo:         txRepo,
		webhookLogRepo: webhookLogRepo,
	}
}

func (s *pluggySyncService) GenerateConnectToken(ctx context.Context, userID uuid.UUID) (*model.ConnectTokenResponse, error) {
	resp, err := s.pluggyClient.GenerateConnectToken(ctx, map[string]interface{}{
		"options": map[string]interface{}{
			"clientUserId": userID.String(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("generate connect token: %w", err)
	}
	return &model.ConnectTokenResponse{AccessToken: resp.AccessToken}, nil
}

func (s *pluggySyncService) ListItems(ctx context.Context, userID uuid.UUID) ([]model.PluggyItem, error) {
	return s.itemRepo.FindByUserID(ctx, userID)
}

func (s *pluggySyncService) DisconnectItem(ctx context.Context, itemID uuid.UUID, userID uuid.UUID) error {
	item, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return ErrNotFound
	}
	if item.UserID != userID {
		return ErrForbidden
	}

	// best-effort delete from Pluggy
	if err := s.pluggyClient.DeleteItem(ctx, item.PluggyItemID); err != nil {
		slog.Warn("failed to delete item from Pluggy", "err", err, "pluggy_item_id", item.PluggyItemID)
	}

	return s.itemRepo.Delete(ctx, itemID, userID)
}

func (s *pluggySyncService) SyncItem(ctx context.Context, pluggyItemID string, userID uuid.UUID) error {
	// Fetch item from Pluggy
	pluggyItem, err := s.pluggyClient.GetItem(ctx, pluggyItemID)
	if err != nil {
		return fmt.Errorf("get item from pluggy: %w", err)
	}

	// Upsert local item record
	now := time.Now()
	localItem := &model.PluggyItem{
		ID:            uuid.New(),
		UserID:        userID,
		PluggyItemID:  pluggyItemID,
		ConnectorName: pluggyItem.Connector.Name,
		ConnectorID:   &pluggyItem.ConnectorID,
		Status:        pluggyItem.Status,
		LastSyncedAt:  &now,
	}
	if err := s.itemRepo.Upsert(ctx, localItem); err != nil {
		return fmt.Errorf("upsert item: %w", err)
	}

	// Fetch accounts
	accountsResp, err := s.pluggyClient.GetAccounts(ctx, pluggyItemID)
	if err != nil {
		return fmt.Errorf("get accounts: %w", err)
	}

	for _, pa := range accountsResp.Results {
		acc := &model.Account{
			ID:              uuid.New(),
			UserID:          userID,
			ItemID:          localItem.ID,
			PluggyAccountID: pa.ID,
			Name:            pa.Name,
			Type:            mapAccountType(pa.Type),
			CurrencyCode:    pa.CurrencyCode,
			Balance:         pa.Balance,
		}
		if pa.Subtype != "" {
			acc.Subtype = &pa.Subtype
		}
		if pa.CreditData != nil {
			acc.CreditLimit = &pa.CreditData.CreditLimit
		}
		if err := s.accountRepo.Upsert(ctx, acc); err != nil {
			slog.Error("upsert account", "err", err, "account_id", pa.ID)
			continue
		}

		// Fetch transactions for this account (first page)
		if err := s.syncTransactions(ctx, acc, userID); err != nil {
			slog.Error("sync transactions", "err", err, "account_id", pa.ID)
		}
	}

	return s.itemRepo.UpdateStatus(ctx, pluggyItemID, pluggyItem.Status, &now)
}

func (s *pluggySyncService) syncTransactions(ctx context.Context, acc *model.Account, userID uuid.UUID) error {
	page := 1
	for {
		resp, err := s.pluggyClient.GetTransactions(ctx, acc.PluggyAccountID, page)
		if err != nil {
			return err
		}

		var txs []model.Transaction
		for _, pt := range resp.Results {
			tx := model.Transaction{
				ID:                  uuid.New(),
				UserID:              userID,
				AccountID:           acc.ID,
				Description:         pt.Description,
				Amount:              pt.Amount,
				Date:                pt.Date,
				Type:                model.TransactionType(pt.Type),
				Tags:                []string{},
			}
			if pt.ID != "" {
				tx.PluggyTransactionID = &pt.ID
			}
			if pt.Category != "" {
				tx.PluggyCategory = &pt.Category
			}
			txs = append(txs, tx)
		}

		if err := s.txRepo.BulkUpsert(ctx, txs); err != nil {
			return err
		}

		if len(resp.Results) < 500 {
			break
		}
		page++
	}
	return nil
}

func (s *pluggySyncService) SyncAllItems(ctx context.Context, userID uuid.UUID) error {
	items, err := s.itemRepo.FindByUserID(ctx, userID)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := s.SyncItem(ctx, item.PluggyItemID, userID); err != nil {
			slog.Error("sync item", "err", err, "pluggy_item_id", item.PluggyItemID)
		}
	}
	return nil
}

func (s *pluggySyncService) HandleWebhook(ctx context.Context, payload model.WebhookPayload) error {
	// Log the webhook
	errMsg := (*string)(nil)
	log := &model.WebhookLog{
		Event:     payload.Event,
		Processed: false,
		Payload:   payload,
	}
	if payload.ItemID != "" {
		log.PluggyItemID = &payload.ItemID
	}
	if err := s.webhookLogRepo.Create(ctx, log); err != nil {
		slog.Error("log webhook", "err", err)
	}

	// Find the local item to get the userID
	localItem, err := s.itemRepo.FindByPluggyItemID(ctx, payload.ItemID)
	if err != nil || localItem == nil {
		slog.Warn("webhook for unknown item", "pluggy_item_id", payload.ItemID)
		_ = errMsg
		return nil
	}

	switch payload.Event {
	case "item/updated", "item/created":
		if err := s.SyncItem(ctx, payload.ItemID, localItem.UserID); err != nil {
			slog.Error("sync after webhook", "err", err)
			return err
		}
	case "item/error":
		status := "ERROR"
		if payload.Error != nil {
			status = payload.Error.Code
		}
		_ = s.itemRepo.UpdateStatus(ctx, payload.ItemID, status, nil)
	}

	return nil
}

func mapAccountType(t string) model.AccountType {
	switch t {
	case "BANK":
		return model.AccountTypeBank
	case "CREDIT":
		return model.AccountTypeCredit
	case "INVESTMENT":
		return model.AccountTypeInvestment
	case "LOAN":
		return model.AccountTypeLoan
	default:
		return model.AccountTypeOther
	}
}
