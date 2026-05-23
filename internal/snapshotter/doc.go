// Package snapshotter provides point-in-time capture and comparison of .env
// file contents across a monorepo.
//
// # Overview
//
// A Snapshot records the parsed key-value entries from a set of .env files at
// a specific moment in time. Snapshots can be serialised to JSON and reloaded
// later, enabling historical auditing and drift detection.
//
// # Usage
//
//	snap, errs := snapshotter.Take("before-deploy", paths)
//	if len(errs) > 0 { /* handle */ }
//	_ = snapshotter.Save(snap, "snap-before.json")
//
//	old, _ := snapshotter.Load("snap-before.json")
//	new, _ := snapshotter.Take("after-deploy", paths)
//	delta := snapshotter.Diff(old, new)
//
package snapshotter
