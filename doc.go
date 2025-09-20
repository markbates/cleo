// Package cleo provides a modern, plugin-based command-line interface framework.
//
// Cleo allows building CLI applications with:
//   - Hierarchical command structures with sub-commands
//   - Plugin-based extensibility for custom functionality
//   - Flexible I/O and filesystem abstraction
//   - Context-aware execution for cancellation and deadlines
//   - Structured logging with slog integration
//   - Memory-efficient operations using object pooling
//   - Modern Go idioms including functional options pattern
//
// # Basic Usage
//
// Create a new command using the functional options pattern:
//
//	cmd, err := cleo.NewCmd("myapp",
//		cleo.WithStdio(iox.Default()),
//		cleo.WithFS(os.DirFS(".")),
//		cleo.WithDescription("My awesome CLI app"),
//		cleo.WithAliases("ma", "myapp"),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Initialize the command with context:
//
//	ctx := context.Background()
//	if err := cmd.InitWithContext(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// # Plugin System
//
// Cleo supports a flexible plugin system. Plugins can implement various interfaces
// to provide functionality:
//
//   - Exiter: Handle exit operations
//   - plugins.FSSetable: Receive filesystem instances
//   - plugins.IOSetable: Receive I/O configurations
//   - plugins.Needer: Receive plugin collections
//
// # Context Support
//
// All major operations support context for cancellation, deadlines, and tracing:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	err := cmd.MainWithContext(ctx, "/current/dir", []string{"arg1", "arg2"})
//
// # Structured Logging
//
// Cleo integrates with Go's structured logging (slog) package:
//
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	cmd, err := cleo.NewCmd("myapp", cleo.WithLogger(logger))
//
// # Memory Efficiency
//
// Cleo uses object pooling internally to reduce garbage collection pressure
// when dealing with plugin collections and other frequently allocated objects.
//
// # Error Handling
//
// Cleo provides typed errors for common conditions:
//
//	if errors.Is(err, cleo.ErrNilCommand) {
//		// Handle nil command error
//	}
//
// Available error types:
//   - ErrNilCommand: Command is nil
//   - ErrNilFS: Filesystem is nil
//   - ErrEmptyCommandName: Command name is empty
//   - ErrNilSubCommand: Sub-command is nil
//   - ErrNoCommand: No command specified
//   - ErrNoCommands: No commands registered
//   - ErrUnknownCommand: Unknown command specified
package cleo