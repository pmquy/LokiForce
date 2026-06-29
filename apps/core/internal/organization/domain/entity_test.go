package domain_test

import (
	"testing"

	"lokiforce.com/apps/core/internal/organization/domain"
)

func TestNewOrganization_Valid(t *testing.T) {
	org, err := domain.NewOrganization("org-123", "LokiForce Org", "Desc", "user-123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if org.Name != "LokiForce Org" {
		t.Errorf("Expected name to be LokiForce Org, got %s", org.Name)
	}
}

func TestNewOrganization_EmptyName(t *testing.T) {
	_, err := domain.NewOrganization("org-123", "", "Desc", "user-123")
	if err == nil {
		t.Errorf("Expected error for empty name, got nil")
	}
}
