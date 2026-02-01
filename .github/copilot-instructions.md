# AI Coding Agent Instructions

## Project Overview
**current-state** is a Go-based Terminal User Interface (TUI) application displaying system information. It uses tview for UI rendering and gopsutil for cross-platform system calls.

### Architecture Layers
1. **UI Layer** (`ui/`): tview-based components (menu, modal) - handles all terminal rendering
2. **Handler Layer** (`handlers/`): Entry points for menu items that format system data into `ListItem` types
3. **System Layer** (`pkg/system/`): Low-level OS abstractions (disk, memory, user info)
4. **Type Definitions** (`types/`): Core data structures (`ListItem`, `DiskStatus`)

## Critical Workflows

### Build & Run
```bash
make build      # Compiles to ./bin/current-state
make run        # Executes the built binary
make all        # Build + run combined
```
Binary output directory: `./bin/`

### Adding New System Info Feature
1. Create handler function in `handlers/` returning `types.ListItem` (see memory.go, disk.go pattern)
2. Create system function in `pkg/system/` for OS-level calls (use syscall for Unix, not gopsutil when possible)
3. Add `ListItem` to `listItems` slice in `ui.CreateMenu()` function
4. Use `system.FormatSize()` for human-readable byte output

## Key Patterns & Conventions

### ListItem Pattern
Each menu item is a `ListItem` struct with:
- `PrimaryText`: Display title (e.g., "Get memory usage")
- `SecondaryText`: Dynamic data or error message
- `Shortcut`: Single rune (e.g., 'a' for memory)
- `Action`: Optional callback (always `nil` for data handlers)
- `Err`: Error field (triggers red error modal if non-nil)

**Handler convention**: Return `ListItem` with default error state, populate fields on success, return early on error.

### Dynamic Updates
The `updateSelectedItem()` goroutine refreshes selected item data every 5 seconds:
- Only updates handlers at indices 1, 2 (memory, disk)
- Skips updates if handler returns error
- Uses `app.QueueUpdateDraw()` for thread-safe UI updates

### Error Handling
Errors are user-visible: when a handler's `Err` field is non-nil, an error modal displays via `GetOKModal()`. Handlers default to informative secondary text messages when data unavailable.

**Error handling pattern** (see `handlers/memory.go`):
```go
func MemoryHandler() types.ListItem {
	listItem := types.ListItem{
		PrimaryText:   "Get memory usage",
		SecondaryText: "No memory usage information available",  // Default fallback
		Shortcut:      'a',
	}

	memUsage, err := mem.VirtualMemory()
	if err != nil {
		listItem.Err = err  // Error set; triggers modal on selection
		return listItem
	}

	listItem.SecondaryText = fmt.Sprintf(
		"All: %s | Used: %s | Available: %s",
		system.FormatSize(memUsage.Total),
		system.FormatSize(memUsage.Used),
		system.FormatSize(memUsage.Available),
	)
	return listItem
}
```

**Key principles**:
- Initialize `ListItem` with `PrimaryText`, default `SecondaryText`, and `nil` `Err`
- Check for errors immediately; if error occurs, set `listItem.Err` and return early
- Never populate `SecondaryText` when `Err` is non-nil
- Use informative default secondary text for when data retrieval fails

## tview/tcell Modal Workflow

The modal system handles user-visible errors and confirmations:

1. **Error Triggering** (`ui/menu.go` - `SetChangedFunc`):
   - When user selects a menu item, `SetChangedFunc` callback fires
   - If `ListItem.Err != nil`, the modal is immediately triggered
   - Secondary text color switches to red for error items

2. **Modal Display** (`ui/modal.go` - `GetOKModal`):
   ```go
   func GetOKModal(app *tview.Application, switchToPrimitive tview.Primitive, text string) *tview.Modal {
       okModal := tview.NewModal().
           SetText(text).
           AddButtons([]string{"OK"}).
           SetDoneFunc(func(int, string) {
               app.SetRoot(switchToPrimitive, true)  // Return to menu on OK
           })
       return okModal
   }
   ```
   - Modal displays error message from handler
   - User presses OK to dismiss and return to menu
   - `app.SetRoot()` switches focus back to menu primitive

3. **Thread-Safe Updates** (`ui/menu.go` - `updateSelectedItem`):
   ```go
   app.QueueUpdateDraw(func() {
       // UI modifications must happen here
       list.SetItemText(currentIndex, freshItem.PrimaryText, freshItem.SecondaryText)
   })
   ```
   - All UI state changes from goroutines use `app.QueueUpdateDraw()`
   - tview event loop processes queued updates atomically
   - Prevents race conditions when updating from dynamic refresh goroutine

**Important**: Never call tview methods directly from goroutines—always queue via `app.QueueUpdateDraw()`.

## Cross-Component Communication

- **main.go → ui.CreateMenu()**: Passes app pointer for Stop() and modal display
- **Handlers → system package**: All OS calls abstracted here (e.g., `system.GetDiskUsage("/")`)
- **UI → Handlers**: Menu queries handlers on selection change to display secondary text
- **System → types**: Handlers format raw system data into `ListItem` and `DiskStatus` structs

## Development Notes

- **Platform-specific code**: Disk usage uses `syscall.Statfs` (Unix) - avoid gopsutil when syscall suffices
- **UI thread safety**: Always queue updates via `app.QueueUpdateDraw()` when modifying UI from goroutines
- **Error propagation**: Handlers never panic; return errors in `ListItem.Err` for modal display
- **Testing patterns**: System functions are pure (no side effects) - unit tests easy to add in `pkg/system/`

## Security Architecture

### Path Validation (`pkg/system/config.go`)
All filesystem queries must validate paths against a whitelist:
```go
diskUsage, err := system.GetDiskUsage("/")  // OK - "/" is whitelisted
diskUsage, err := system.GetDiskUsage("/etc/shadow")  // ERROR - not whitelisted
```
- Prevents arbitrary filesystem traversal
- Edit `AllowedPaths` map to permit new paths
- Use `system.AddAllowedPath()` for runtime additions

### Error Sanitization (`pkg/system/errors.go`)
All handler errors must be sanitized before display:
```go
err := someSystemCall()
listItem.Err = system.SanitizeError(err)  // Maps OS errors to safe messages
```
- Maps `EACCES` → "permission denied: insufficient access..."
- Maps `ENOENT` → "filesystem not found: path may not exist..."
- Prevents information leakage about system configuration

### Goroutine Lifecycle (`ui/menu.go`)
Dynamic update goroutines use context for graceful cleanup:
```go
ctx, cancel := context.WithCancel(context.Background())
go updateSelectedItem(ctx, app, list, listItems)
app.SetDoneFunc(func() {
    cancel()  // Stops goroutine when app exits
})
```
- Prevents goroutine leaks on app termination
- Allows graceful shutdown of long-running operations
- Standard Go idiom for cancellation
