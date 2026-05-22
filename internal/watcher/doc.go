// Package watcher provides lightweight polling-based file watching for .env
// files within a monorepo.
//
// Instead of relying on OS-level inotify or FSEvents (which require CGO or
// platform-specific builds), watcher uses a configurable tick interval to
// stat each registered path and compare modification times against a
// previously captured snapshot.
//
// # Basic usage
//
//	paths := []string{".env", "services/api/.env"}
//	w := watcher.New(paths, 500*time.Millisecond)
//	w.Start()
//	defer w.Stop()
//
//	for ev := range w.Events {
//		fmt.Printf("%s: %s\n", ev.Op, ev.Path)
//	}
//
// # Event operations
//
// Each Event carries one of three Op strings:
//   - "created"  — the file did not exist at startup and now does
//   - "modified" — the file's modification time has advanced
//   - "deleted"  — the file existed at startup and has been removed
//
// Non-fatal stat errors (e.g. permission issues on first encounter) are
// forwarded to the Errors channel without stopping the watcher.
package watcher
