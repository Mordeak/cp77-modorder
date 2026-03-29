// Package modlist writes the modlist.txt file consumed by the CP2077 mod loader.
package modlist

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mordeak/cp77-modorder-gui/internal/conflict"
)

// Write writes modlist.txt to dir, one archive name per line in load order.
// The first entry loads first and wins all conflicts.
// Any existing modlist.txt is backed up into a modlist.old subfolder before overwriting.
func Write(dir string, mods []*conflict.ModInfo) error {
	dest := filepath.Join(dir, "modlist.txt")

	// Back up existing file into modlist.old/.
	if _, err := os.Stat(dest); err == nil {
		backupDir := filepath.Join(dir, "modlist.old")
		if err := os.MkdirAll(backupDir, 0o755); err != nil {
			return fmt.Errorf("create modlist.old: %w", err)
		}
		ts := time.Now().Format("2006-01-02_15-04-05")
		backup := filepath.Join(backupDir, "modlist.txt."+ts)
		if err := os.Rename(dest, backup); err != nil {
			return fmt.Errorf("backup modlist.txt: %w", err)
		}
	}

	var sb strings.Builder
	for _, m := range mods {
		sb.WriteString(m.Name)
		sb.WriteByte('\n')
	}

	if err := os.WriteFile(dest, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("write modlist.txt: %w", err)
	}
	return nil
}
