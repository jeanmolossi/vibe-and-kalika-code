---
name: go-atomic-binary-replace
description: Use this skill when a Go CLI must atomically replace its own running binary (e.g., self-update). Call it before implementing any file-staging or os.Rename logic. It prevents the EXDEV cross-device link error that occurs when the temp dir lives on a different filesystem than the target binary.
---

# Go Atomic Binary Replace

## Purpose

Stage and atomically swap a running binary on the same filesystem to avoid the
`EXDEV: invalid cross-device link` error that `os.Rename` returns when source
and destination are on different filesystems (e.g., staging in `/tmp` while the
binary lives in `/usr/local/bin`).

The fix is: resolve the running binary path first with `os.Executable()`, then
create the temp directory under `filepath.Dir(execPath)` so source and destination
are guaranteed to be on the same device.

## Trigger

Call this skill when:

- A Go CLI implements a self-update that downloads a new binary and must replace the running executable.
- Any Go code uses `os.Rename` to atomically replace a file and the staging directory is not guaranteed to be on the same filesystem.
- A code review or CI failure reports `rename ... invalid cross-device link` or `EXDEV`.
- You are about to call `os.MkdirTemp("", ...)` to stage a file that will be renamed into a system directory (`/usr/local/bin`, `/opt/...`, etc.).

## Do not use when

Do not call this skill when:

- The source and destination are always on the same known filesystem (e.g., replacing a file within the same project directory).
- You do not use `os.Rename` for the final swap (e.g., you copy and delete instead).
- The binary is not being replaced at runtime.

## Required inputs

The agent must provide:

- The path of the file to be replaced (destination).
- The code location where `os.MkdirTemp` is called or will be called.

## Procedure

1. **Resolve the running executable path before creating any temp dir:**
   ```go
   execPath, err := os.Executable()
   if err != nil {
       return fmt.Errorf("resolve executable: %w", err)
   }
   execDir := filepath.Dir(execPath)
   ```

2. **Create the staging temp dir under the same directory as the binary:**
   ```go
   // Stage on the same filesystem as the binary to prevent EXDEV on os.Rename.
   tmpDir, err := os.MkdirTemp(execDir, "myapp-update-*")
   if err != nil {
       return fmt.Errorf("create temp dir: %w", err)
   }
   defer os.RemoveAll(tmpDir)
   ```

3. **Write / extract the new binary into the temp dir:**
   ```go
   binPath := filepath.Join(tmpDir, "myapp")
   // ... write or extract downloaded binary to binPath ...
   if err := os.Chmod(binPath, 0o755); err != nil {
       return fmt.Errorf("chmod: %w", err)
   }
   ```

4. **Atomically swap with os.Rename:**
   ```go
   // os.Rename is atomic on Unix when src and dst are on the same filesystem.
   if err := os.Rename(binPath, execPath); err != nil {
       return fmt.Errorf("replace binary: %w", err)
   }
   ```

## Expected output

```md
# Atomic Binary Replace Result

## Status
PASS | BLOCKED | NOT_APPLICABLE

## Findings
- execPath resolved: <path>
- tmpDir staged under: <execDir>
- os.Rename: succeeded / failed with <error>

## Next step
- If PASS: binary replaced; caller should print version confirmation.
- If BLOCKED: document why same-filesystem staging is not possible; consider os.Link + os.Remove fallback.
```

## Stop conditions

Stop and return `BLOCKED` when:

- `os.Executable()` fails (e.g., the binary was deleted before the update ran).
- The target directory is not writable by the running process (requires privilege escalation outside scope of this skill).
- The OS is Windows: `os.Rename` on Windows fails if the destination file is open/locked; use a reboot-rename strategy instead.

## Notes

- **Why not `/tmp`?**: On Linux, `/tmp` is often on `tmpfs` while system directories use `ext4`/`xfs`. `os.Rename` is a single syscall (`rename(2)`) that only works within the same mounted filesystem. Crossing mount points returns `EXDEV`.
- **Why not `io.Copy` + `os.Remove`?**: Non-atomic. A crash between copy and remove leaves a corrupt or missing binary.
- **Windows caveat**: Windows locks open executables; `os.Rename` returns `Access is denied`. A common workaround is to move the old binary to a `.old` name and place the new one, then delete `.old` on next startup.
