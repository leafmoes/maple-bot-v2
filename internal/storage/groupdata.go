package storage

import (
	"encoding/json"
	"fmt"
	"maple-bot/internal/bot/model"
	"maple-bot/internal/config"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type KVStorage struct {
	kv *CFKVClient
}

func NewKVStorage() *KVStorage {
	cfg := config.App.CloudflareKV
	return &KVStorage{kv: NewCFClient(cfg.AccountId, cfg.NamespaceId, cfg.ApiToken)}
}

func (ks *KVStorage) GetGroupData(ctx *ext.Context) (*model.GroupData, bool) {
	groupId := ctx.EffectiveChat.Id
	key := fmt.Sprintf("group_data_%s", groupId)
	data, err := ks.kv.Read(key)
	if err != nil {
		fmt.Errorf("failed to get kv data:%w", err.Error())
		return nil, false
	}
	var v *model.GroupData
	if err := json.Unmarshal(data, v); err != nil {
		fmt.Errorf("failed to unmarshal kv data:%w", err.Error())
		return nil, false
	}
	return v, true
}

func (ks *KVStorage) SetGroupData(ctx *ext.Context, value *model.GroupData) bool {
	groupId := ctx.EffectiveChat.Id
	key := fmt.Sprintf("group_data_%s", groupId)
	v, err := json.Marshal(value)
	if err != nil {
		fmt.Errorf("failed to marshal kv data:%w", err.Error())
		return false
	}
	ok, err := ks.kv.Write(key, v)
	if err != nil || !ok {
		fmt.Errorf("failed to set kv data:%w", err.Error())
		return false
	}
	return true
}
