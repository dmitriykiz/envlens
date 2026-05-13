// Package reporter provides formatting and output utilities for envlens audit
// results. It supports human-readable text output and per-file summary lines
// suitable for CI pipelines or terminal use.
//
// # Usage
//
// Build a slice of [Report] values — each pairing a file path with its
// [analyzer.Result] — then pass them to [WriteText] to render a structured
// report to any [io.Writer]:
//
//	reports := []reporter.Report{
//		{FilePath: ".env", Result: result},
//	}
//	reporter.WriteText(os.Stdout, reports)
//
// For a concise single-line summary per file, use [Summary]:
//
//	fmt.Println(reporter.Summary(reports[0]))
package reporter
