# Plan: Add 5-Second CPU Profiling to DumpPprofToBytes

## Context
`DumpPprofToBytes()` in `ForwardToFile.go` currently returns only a heap memory snapshot. We want to add a 5-second CPU profile to the returned data so callers get both memory and CPU analysis in a single call.

## Key Consideration
Heap and CPU profiles are **separate pprof formats** — they cannot be merged into a single `[]byte` blob. A pprof consumer (e.g., `go tool pprof`) expects one profile type per file. Therefore the function must return both profiles separately.

## Approach
Return a struct (or two byte slices) containing both the heap profile and the CPU profile. The CPU profile is collected by calling `pprof.StartCPUProfile`, sleeping 5 seconds, then calling `pprof.StopCPUProfile`.

## Changes

### 1. Modify `ForwardToFile.go` (non-windows)

- Add `"time"` to imports.
- Add a new function `DumpFullPprofToBytes() (heap []byte, cpu []byte, err error)`:
  1. Start CPU profiling into a `bytes.Buffer` via `pprof.StartCPUProfile(&cpuBuf)`.
  2. Sleep for 5 seconds (`time.Sleep(5 * time.Second)`).
  3. Stop CPU profiling via `pprof.StopCPUProfile()`.
  4. Capture heap profile into a separate `bytes.Buffer` via `pprof.WriteHeapProfile(&heapBuf)`.
  5. Return both `heapBuf.Bytes()` and `cpuBuf.Bytes()`.
- Keep the existing `DumpPprofToBytes()` unchanged (backward compatibility).

### 2. Modify `ForwardToFile_windows.go`

- Add the same `DumpFullPprofToBytes()` signature so the package compiles on Windows.
- Check if `pprof.StartCPUProfile` / `pprof.StopCPUProfile` work on Windows (they do — standard library, no OS-specific constraint). Implement identically.

### 3. Update `DumpPprofToFile` (optional, for parity)

- Add a companion `DumpFullPprofToFile(path, alias string) error` that writes both `{alias}.dat` (heap) and `{alias}-cpu.dat` (CPU) files.
- Or skip this if file-based dumping of CPU profiles is not needed right now.

## Notes
- `pprof.StartCPUProfile` returns an error if CPU profiling is already active — the function should check and return this error.
- The 5-second sleep blocks the calling goroutine. Callers should be aware this is a blocking call.
- The heap snapshot is taken **after** the CPU profiling window, so it reflects memory state at the end of the 5-second window.
