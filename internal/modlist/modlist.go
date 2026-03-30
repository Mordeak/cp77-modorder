// Package modlist writes the modlist.txt file consumed by the CP2077 mod loader.
package modlist

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Mordeak/cp77-modorder-gui/internal/conflict"
)

// Write writes modlist.txt to dir, one archive name per line in load order.
// The first entry loads first and wins all conflicts.
// Any existing modlist.txt is backed up into a modlist.old subfolder before overwriting.
// tag is appended to the backup filename when non-empty (e.g. "re-order").
func Write(dir string, mods []*conflict.ModInfo, tag string) error {
	dest := filepath.Join(dir, "modlist.txt")

	// Back up existing file into modlist.old/.
	if _, err := os.Stat(dest); err == nil {
		backupDir := filepath.Join(dir, "modlist.old")
		if err := os.MkdirAll(backupDir, 0o755); err != nil {
			return fmt.Errorf("create modlist.old: %w", err)
		}
		ts := time.Now().Format("2006-01-02_15-04-05")
		backupName := "modlist.txt." + ts
		if tag != "" {
			backupName += "." + tag
		}
		if err := os.Rename(dest, filepath.Join(backupDir, backupName)); err != nil {
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

// PruneBackups removes the oldest backup files from backupDir so that at most
// limit files remain. Files are sorted lexicographically; because backup names
// embed a timestamp (2006-01-02_15-04-05) the lexicographic order equals
// chronological order, so the oldest are removed first.
// A limit of 0 or less is a no-op.
func PruneBackups(backupDir string, limit int) error {
	if limit <= 0 {
		return nil
	}
	entries, err := os.ReadDir(backupDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read backup dir: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	if len(names) <= limit {
		return nil
	}
	sort.Strings(names) // oldest first
	for _, name := range names[:len(names)-limit] {
		_ = os.Remove(filepath.Join(backupDir, name))
	}
	return nil
}
