package firewall

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/openforge-maintain/openforge-maintain/agent/model"
	"github.com/openforge-maintain/openforge-maintain/agent/repository"
)

// FirewallService provides firewall management operations using iptables/ufw.
type FirewallService struct {
	repo *repository.FirewallRepository
}

// NewFirewallService creates a new FirewallService.
func NewFirewallService(repo *repository.FirewallRepository) *FirewallService {
	return &FirewallService{repo: repo}
}

// ListRules returns all firewall rules.
func (s *FirewallService) ListRules() ([]model.FirewallRule, error) {
	return s.repo.List()
}

// AddRule adds a new firewall rule.
func (s *FirewallService) AddRule(rule *model.FirewallRule) error {
	rule.CreatedAt = time.Now()

	if err := s.repo.Create(rule); err != nil {
		return fmt.Errorf("failed to save rule: %w", err)
	}

	// Apply the rule to the system firewall
	return s.applyRule(rule)
}

// DeleteRule removes a firewall rule.
func (s *FirewallService) DeleteRule(id uint) error {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get rule: %w", err)
	}

	// Remove the rule from the system firewall
	s.removeRule(rule)

	return s.repo.Delete(id)
}

// EnableRule enables a previously disabled firewall rule.
func (s *FirewallService) EnableRule(id uint) error {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get rule: %w", err)
	}

	if err := s.applyRule(rule); err != nil {
		return fmt.Errorf("failed to apply rule: %w", err)
	}
	return nil
}

// DisableRule disables a firewall rule without deleting it.
func (s *FirewallService) DisableRule(id uint) error {
	rule, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get rule: %w", err)
	}

	return s.removeRule(rule)
}

// GetStatus returns the current firewall status.
func (s *FirewallService) GetStatus() (string, error) {
	cmd := exec.Command("ufw", "status")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown", fmt.Errorf("failed to get firewall status: %s: %w", string(output), err)
	}
	return strings.TrimSpace(string(output)), nil
}

// Enable enables the firewall.
func (s *FirewallService) Enable() error {
	cmd := exec.Command("ufw", "--force", "enable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable firewall: %s: %w", string(output), err)
	}
	return nil
}

// Disable disables the firewall.
func (s *FirewallService) Disable() error {
	cmd := exec.Command("ufw", "--force", "disable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable firewall: %s: %w", string(output), err)
	}
	return nil
}

// applyRule applies a firewall rule to the system.
func (s *FirewallService) applyRule(rule *model.FirewallRule) error {
	var cmd *exec.Cmd
	switch rule.Action {
	case "allow":
		cmd = exec.Command("ufw", "allow", rule.Protocol+"/"+rule.Port, "from", rule.Source)
	case "deny":
		cmd = exec.Command("ufw", "deny", rule.Protocol+"/"+rule.Port, "from", rule.Source)
	default:
		return fmt.Errorf("unknown action: %s", rule.Action)
	}

	if rule.Comment != "" {
		cmd.Args = append(cmd.Args, "comment", rule.Comment)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply rule: %s: %w", string(output), err)
	}
	return nil
}

// removeRule removes a firewall rule from the system.
func (s *FirewallService) removeRule(rule *model.FirewallRule) error {
	cmd := exec.Command("ufw", "delete", rule.Action, rule.Protocol+"/"+rule.Port, "from", rule.Source)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove rule: %s: %w", string(output), err)
	}
	return nil
}

// SyncRules synchronizes all rules in the database with the system firewall.
func (s *FirewallService) SyncRules() error {
	// Reset firewall rules
	resetCmd := exec.Command("ufw", "--force", "reset")
	resetCmd.CombinedOutput()

	// Re-enable and re-apply all rules
	enableCmd := exec.Command("ufw", "--force", "enable")
	enableCmd.CombinedOutput()

	rules, err := s.repo.List()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		_ = s.applyRule(&rule)
	}

	return nil
}
