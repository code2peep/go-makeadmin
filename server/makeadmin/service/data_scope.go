package service

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

	"gorm.io/gorm"

	"go-makeadmin/makeadmin/repository"
	"go-makeadmin/model/makeadmin"
)

func (srv authService) resolveDataScope(ctx context.Context, tenantID uint64, adminID uint64, roleIDs []uint64) (repository.DataScopeFilter, error) {
	scopeFilter := repository.DataScopeFilter{Enabled: true, AdminID: adminID}
	if len(roleIDs) == 0 {
		scopeFilter.Self = true
		return scopeFilter, nil
	}

	scopes, err := srv.repo.ListDataScopesByRoleIDs(ctx, tenantID, roleIDs)
	if err != nil {
		return repository.DataScopeFilter{}, err
	}
	if len(scopes) == 0 {
		scopeFilter.Self = true
		return scopeFilter, nil
	}

	var orgIDs []uint64
	for _, scope := range scopes {
		switch scope.ScopeType {
		case makeadmin.ScopeTypeAll:
			return repository.DataScopeFilter{Enabled: true, All: true, AdminID: adminID}, nil
		case makeadmin.ScopeTypeSelf:
			scopeFilter.Self = true
		case makeadmin.ScopeTypeOrg:
			orgID, err := srv.primaryOrgID(ctx, tenantID, adminID)
			if err != nil {
				return repository.DataScopeFilter{}, err
			}
			orgIDs = appendOrgID(orgIDs, orgID)
		case makeadmin.ScopeTypeOrgTree:
			ids, err := srv.primaryOrgTreeIDs(ctx, tenantID, adminID)
			if err != nil {
				return repository.DataScopeFilter{}, err
			}
			orgIDs = append(orgIDs, ids...)
		case makeadmin.ScopeTypeCustomOrg:
			orgIDs = append(orgIDs, customOrgIDs(scope.ScopeValue)...)
		}
	}

	scopeFilter.OrgIDs = uniqueUint64(orgIDs)
	if !scopeFilter.Self && len(scopeFilter.OrgIDs) == 0 {
		scopeFilter.NoAccess = true
	}
	return scopeFilter, nil
}

func (srv authService) primaryOrgID(ctx context.Context, tenantID uint64, adminID uint64) (uint64, error) {
	adminOrg, err := srv.repo.FindPrimaryAdminOrg(ctx, tenantID, adminID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return adminOrg.OrgID, nil
}

func (srv authService) primaryOrgTreeIDs(ctx context.Context, tenantID uint64, adminID uint64) ([]uint64, error) {
	rootID, err := srv.primaryOrgID(ctx, tenantID, adminID)
	if err != nil || rootID == 0 {
		return []uint64{}, err
	}
	orgs, err := srv.repo.ListOrgUnits(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	children := make(map[uint64][]uint64, len(orgs))
	for _, org := range orgs {
		children[org.ParentID] = append(children[org.ParentID], org.ID)
	}
	result := []uint64{rootID}
	queue := []uint64{rootID}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, childID := range children[current] {
			result = append(result, childID)
			queue = append(queue, childID)
		}
	}
	return uniqueUint64(result), nil
}

func customOrgIDs(value string) []uint64 {
	var raw interface{}
	if err := json.Unmarshal([]byte(value), &raw); err != nil {
		return []uint64{}
	}
	return orgIDsFromJSON(raw)
}

func orgIDsFromJSON(value interface{}) []uint64 {
	switch item := value.(type) {
	case []interface{}:
		result := make([]uint64, 0, len(item))
		for _, rawID := range item {
			if id := uint64FromJSON(rawID); id > 0 {
				result = append(result, id)
			}
		}
		return result
	case map[string]interface{}:
		for _, key := range []string{"org_ids", "orgIds", "ids"} {
			if ids := orgIDsFromJSON(item[key]); len(ids) > 0 {
				return ids
			}
		}
	}
	return []uint64{}
}

func uint64FromJSON(value interface{}) uint64 {
	switch item := value.(type) {
	case float64:
		if item > 0 {
			return uint64(item)
		}
	case int:
		if item > 0 {
			return uint64(item)
		}
	case uint64:
		return item
	}
	return 0
}

func appendOrgID(values []uint64, id uint64) []uint64 {
	if id == 0 {
		return values
	}
	return append(values, id)
}

func uniqueUint64(values []uint64) []uint64 {
	if len(values) == 0 {
		return []uint64{}
	}
	seen := make(map[uint64]struct{}, len(values))
	result := make([]uint64, 0, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}
