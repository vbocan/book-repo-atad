//! Exercise 15.1: A Pluggable Logging System
//!
//! Reimplements the Go logging system from Exercise 6.3 (Chapter 6) using
//! Rust traits instead of Go interfaces. `Logger<F: Formatter>` drives a
//! formatter through static dispatch; a `Vec<Box<dyn Formatter>>` drives a
//! heterogeneous collection of formatters through dynamic dispatch.

/// The severity of a log entry, from least to most severe.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Level {
    Info,
    Warn,
    Error,
}

impl Level {
    /// The level's canonical uppercase name, e.g. `"INFO"`.
    fn as_str(&self) -> &'static str {
        match self {
            Level::Info => "INFO",
            Level::Warn => "WARN",
            Level::Error => "ERROR",
        }
    }
}

/// A single log record: a severity level, a human-readable message, and
/// the Unix timestamp (seconds since the epoch) at which it was created.
#[derive(Debug, Clone)]
struct LogEntry {
    level: Level,
    message: String,
    timestamp: u64,
}

impl LogEntry {
    /// Builds a new entry, stamping it with the current time. `message`
    /// accepts anything convertible to `String`, so callers can pass a
    /// `&str` literal or an owned `String` without extra ceremony.
    fn new(level: Level, message: impl Into<String>) -> Self {
        LogEntry {
            level,
            message: message.into(),
            timestamp: current_unix_timestamp(),
        }
    }
}

/// Seconds elapsed since the Unix epoch, used to stamp each `LogEntry`.
fn current_unix_timestamp() -> u64 {
    std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .expect("system clock is set before the Unix epoch")
        .as_secs()
}

/// Renders a `LogEntry` as a string. `Logger<F>` (static dispatch) and the
/// `Vec<Box<dyn Formatter>>` demo (dynamic dispatch) both drive
/// implementations of this trait through the same single method, so
/// neither cares how a given formatter builds its output.
trait Formatter {
    fn format(&self, entry: &LogEntry) -> String;
}

/// Formats an entry as one line of plain text: `[LEVEL] message (timestamp)`.
struct TextFormatter;

impl Formatter for TextFormatter {
    fn format(&self, entry: &LogEntry) -> String {
        format!(
            "[{}] {} ({})",
            entry.level.as_str(),
            entry.message,
            entry.timestamp
        )
    }
}

/// Formats an entry as one line of JSON. Hand-rolled with `format!` rather
/// than `serde_json`, since this example is standard-library only; the
/// `escape` helper takes care of the characters that would otherwise
/// produce invalid JSON.
struct JsonFormatter;

impl JsonFormatter {
    /// Escapes the characters JSON forbids inside a string literal.
    fn escape(text: &str) -> String {
        let mut escaped = String::with_capacity(text.len());
        for c in text.chars() {
            match c {
                '"' => escaped.push_str("\\\""),
                '\\' => escaped.push_str("\\\\"),
                '\n' => escaped.push_str("\\n"),
                '\r' => escaped.push_str("\\r"),
                '\t' => escaped.push_str("\\t"),
                c if (c as u32) < 0x20 => escaped.push_str(&format!("\\u{:04x}", c as u32)),
                c => escaped.push(c),
            }
        }
        escaped
    }
}

impl Formatter for JsonFormatter {
    fn format(&self, entry: &LogEntry) -> String {
        format!(
            "{{\"level\":\"{}\",\"message\":\"{}\",\"timestamp\":{}}}",
            entry.level.as_str(),
            JsonFormatter::escape(&entry.message),
            entry.timestamp
        )
    }
}

/// Drives logging through a formatter fixed at compile time. `F` is a
/// concrete type parameter, not `dyn Formatter`, so the compiler
/// monomorphizes a distinct `Logger<F>` — and a distinct `log` — for every
/// formatter type it is instantiated with below. The call to
/// `self.formatter.format(entry)` resolves to a direct, inlinable function
/// call with no vtable involved: this is static dispatch.
struct Logger<F: Formatter> {
    formatter: F,
}

impl<F: Formatter> Logger<F> {
    fn new(formatter: F) -> Self {
        Logger { formatter }
    }

    fn log(&self, entry: &LogEntry) {
        println!("{}", self.formatter.format(entry));
    }
}

/// Demonstrates static dispatch: one `Logger<TextFormatter>` and one
/// `Logger<JsonFormatter>`, each a different monomorphized type even
/// though both are "a `Logger`". The compiler generates specialized code
/// for each instantiation, so there is no runtime cost to the abstraction.
fn run_static_dispatch_demo(entries: &[LogEntry]) {
    println!("=== Static dispatch: Logger<F: Formatter> ===");

    println!("\n-- Logger<TextFormatter> --");
    let text_logger = Logger::new(TextFormatter);
    for entry in entries {
        text_logger.log(entry);
    }

    println!("\n-- Logger<JsonFormatter> --");
    let json_logger = Logger::new(JsonFormatter);
    for entry in entries {
        json_logger.log(entry);
    }
}

/// Demonstrates dynamic dispatch: a single `Vec<Box<dyn Formatter>>` holds
/// a `TextFormatter` and a `JsonFormatter` side by side, something no
/// `Logger<F>` can do because `F` must be one fixed type for the whole
/// collection. Each `formatter.format(entry)` call is resolved at runtime
/// through the trait object's vtable, at the cost of that indirection.
fn run_dynamic_dispatch_demo(entries: &[LogEntry]) {
    println!("=== Dynamic dispatch: Vec<Box<dyn Formatter>> ===\n");

    let formatters: Vec<Box<dyn Formatter>> = vec![
        Box::new(TextFormatter),
        Box::new(JsonFormatter),
    ];

    for entry in entries {
        for formatter in &formatters {
            println!("{}", formatter.format(entry));
        }
    }
}

fn main() {
    let entries = vec![
        LogEntry::new(Level::Info, "server started on port 8080"),
        LogEntry::new(Level::Warn, "cache miss rate above 20%"),
        LogEntry::new(Level::Error, "failed to connect to database"),
    ];

    run_static_dispatch_demo(&entries);
    println!();
    run_dynamic_dispatch_demo(&entries);
}
