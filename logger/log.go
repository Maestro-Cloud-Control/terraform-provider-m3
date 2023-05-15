package logger

type Log interface {
    // Error emit a message and key/value pairs at the ERROR level
    Error(msg string, args ...interface{})
    // Debug Emit a message and key/value pairs at the DEBUG level
    Debug(msg string, args ...interface{})
    // Info Emit a message and key/value pairs at the INFO level
    Info(msg string, args ...interface{})
}
