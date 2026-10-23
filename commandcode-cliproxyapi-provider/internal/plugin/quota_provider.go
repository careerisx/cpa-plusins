package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/router-for-me/CLIProxyAPI/v8/sdk/pluginabi"
	"github.com/router-for-me/CLIProxyAPI/v8/sdk/pluginapi"
)

type quotaCallbackContextKey struct{}

var _ pluginapi.QuotaProvider = (*Manager)(nil)

func (m *Manager) Identifier() string { return ProviderID }

func (m *Manager) DescribeQuota(context.Context, pluginapi.QuotaDescribeRequest) (pluginapi.QuotaDescribeResponse, error) {
	return pluginapi.QuotaDescribeResponse{SupportedProviders: []string{ProviderID}, DisplayName: "CommandCode", SupportsReset: false}, nil
}

func (m *Manager) FetchQuota(ctx context.Context, req pluginapi.QuotaFetchRequest) (pluginapi.QuotaFetchResponse, error) {
	if req.Provider != ProviderID {
		return pluginapi.QuotaFetchResponse{}, fmt.Errorf("unsupported quota provider")
	}
	key := selectedAuthKey(req.Attributes, req.Metadata)
	if key == "" {
		return pluginapi.QuotaFetchResponse{}, fmt.Errorf("selected credential has no API key")
	}
	m.mu.RLock()
	base, timeout := m.cfg.BaseURL, m.cfg.RequestTimeout
	m.mu.RUnlock()
	usage, _, _, err := fetchQuota(ctx, m.bridge, base, timeout, key)
	if err != nil {
		return pluginapi.QuotaFetchResponse{}, fmt.Errorf("quota query failed: %w", err)
	}
	return nativeQuota(usage), nil
}

func (m *Manager) ResetQuota(context.Context, pluginapi.QuotaResetRequest) (pluginapi.QuotaResetResponse, error) {
	return pluginapi.QuotaResetResponse{}, fmt.Errorf("CommandCode quota reset is not supported")
}

func nativeQuota(u quotaUsage) pluginapi.QuotaFetchResponse {
	out := pluginapi.QuotaFetchResponse{}
	if u.Plan != "" {
		out.Subscription = &pluginapi.QuotaSubscription{Plan: u.Plan}
	}
	if u.creditsPresent {
		out.Summary = append(out.Summary, pluginapi.QuotaMetric{Key: "monthly_credits_remaining", Label: "Monthly credits remaining", Value: u.CreditsLeft, Unit: "credits", Format: "number"})
	}
	if u.purchasedPresent {
		out.Summary = append(out.Summary, pluginapi.QuotaMetric{Key: "purchased_credits", Label: "Purchased credits", Value: u.PurchasedCredits, Unit: "credits", Format: "number"})
	}
	buckets := []pluginapi.QuotaBucket{}
	add := func(id string, w *quotaWindow) {
		// v8.0.0 has no unknown/unlimited representation. Omit those windows.
		if w == nil || w.Status == "unknown" || w.Status == "unlimited" || w.Cap <= 0 {
			return
		}
		buckets = append(buckets, pluginapi.QuotaBucket{Window: id, RemainingFraction: math.Max(0, math.Min(1, 1-w.Used/w.Cap)), ResetTime: w.ResetsAt})
	}
	add("monthly", u.Month)
	add("5h", &u.FiveHour)
	add("weekly", &u.Weekly)
	if len(buckets) > 0 {
		out.Groups = []pluginapi.QuotaGroup{{DisplayName: "CommandCode", Buckets: buckets}}
	}
	return out
}

func (m *Manager) handleQuota(method string, raw []byte) ([]byte, error) {
	if method == pluginabi.MethodQuotaIdentifier {
		return okEnvelope(map[string]string{"identifier": ProviderID}), nil
	}
	// Reject null as well as malformed JSON rather than accepting a zero request.
	if !json.Valid(raw) || !bytes.HasPrefix(bytes.TrimSpace(raw), []byte("{")) {
		return ErrEnvelope("invalid_request", "malformed quota request"), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), quotaMaxTimeout)
	defer cancel()
	var result any
	var err error
	switch method {
	case pluginabi.MethodQuotaDescribe:
		var req pluginapi.QuotaDescribeRequest
		if json.Unmarshal(raw, &req) != nil {
			return ErrEnvelope("invalid_request", "malformed quota request"), nil
		}
		result, err = m.DescribeQuota(ctx, req)
	case pluginabi.MethodQuotaFetch:
		var req struct {
			pluginapi.QuotaFetchRequest
			HostCallbackID string `json:"host_callback_id,omitempty"`
		}
		if json.Unmarshal(raw, &req) != nil {
			return ErrEnvelope("invalid_request", "malformed quota request"), nil
		}
		ctx = context.WithValue(ctx, quotaCallbackContextKey{}, req.HostCallbackID)
		result, err = m.FetchQuota(ctx, req.QuotaFetchRequest)
	case pluginabi.MethodQuotaReset:
		var req pluginapi.QuotaResetRequest
		if json.Unmarshal(raw, &req) != nil {
			return ErrEnvelope("invalid_request", "malformed quota request"), nil
		}
		result, err = m.ResetQuota(ctx, req)
	}
	if err != nil {
		return ErrEnvelope("quota_failure", err.Error()), nil
	}
	return okEnvelope(result), nil
}
