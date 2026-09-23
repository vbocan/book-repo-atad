use std::collections::HashMap;
use std::fmt;
use std::fs;

#[derive(Debug, Clone, PartialEq)]
enum ConfigValue {
    Text(String),
    Number(i64),
    Float(f64),
    Boolean(bool),
    List(Vec<ConfigValue>),
}

impl fmt::Display for ConfigValue {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ConfigValue::Text(s) => write!(f, "{s}"),
            ConfigValue::Number(n) => write!(f, "{n}"),
            ConfigValue::Float(x) => write!(f, "{x}"),
            ConfigValue::Boolean(b) => write!(f, "{b}"),
            ConfigValue::List(items) => {
                let rendered: Vec<String> = items.iter().map(ToString::to_string).collect();
                write!(f, "[{}]", rendered.join(", "))
            }
        }
    }
}

#[derive(Debug)]
enum ConfigError {
    Io(std::io::Error),
    InvalidLine { line_number: usize, content: String },
}

impl fmt::Display for ConfigError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            ConfigError::Io(e) => write!(f, "could not read config file: {e}"),
            ConfigError::InvalidLine { line_number, content } => {
                write!(f, "line {line_number} is not a valid \"key = value\" pair: {content:?}")
            }
        }
    }
}

impl std::error::Error for ConfigError {}

impl From<std::io::Error> for ConfigError {
    fn from(err: std::io::Error) -> Self {
        ConfigError::Io(err)
    }
}

/// Classifies a single raw value from the right-hand side of a "key = value"
/// line. This never fails: a value that isn't an integer, a float, or one
/// of the booleans "true"/"false" is kept as text.
fn parse_value(raw: &str) -> ConfigValue {
    let trimmed = raw.trim();

    if let Ok(n) = trimmed.parse::<i64>() {
        return ConfigValue::Number(n);
    }
    if let Ok(x) = trimmed.parse::<f64>() {
        return ConfigValue::Float(x);
    }
    match trimmed {
        "true" => return ConfigValue::Boolean(true),
        "false" => return ConfigValue::Boolean(false),
        _ => {}
    }

    ConfigValue::Text(trimmed.to_string())
}

/// Parses a simple key-value configuration file: one "key = value" pair per
/// line, blank lines ignored, and lines whose first non-whitespace
/// character is '#' treated as comments. Returns an error for a malformed
/// line (no '=') or if the file cannot be read.
fn parse_config(path: &str) -> Result<HashMap<String, ConfigValue>, ConfigError> {
    let content = fs::read_to_string(path)?;
    let mut config = HashMap::new();

    for (i, raw_line) in content.lines().enumerate() {
        let line = raw_line.trim();
        if line.is_empty() || line.starts_with('#') {
            continue;
        }

        let (key, value) = line.split_once('=').ok_or_else(|| ConfigError::InvalidLine {
            line_number: i + 1,
            content: raw_line.to_string(),
        })?;

        config.insert(key.trim().to_string(), parse_value(value));
    }

    Ok(config)
}

fn main() {
    let sample_path = "sample_config.txt";
    let sample = "\
# Server configuration
name = Hello World
port = 8080
timeout = 2.5
debug = true
verbose = false
";
    fs::write(sample_path, sample).expect("failed to write sample config");

    match parse_config(sample_path) {
        Ok(config) => {
            let mut keys: Vec<&String> = config.keys().collect();
            keys.sort();
            for key in keys {
                println!("{key} = {} ({:?})", config[key], config[key]);
            }
        }
        Err(e) => eprintln!("Failed to parse config: {e}"),
    }

    fs::remove_file(sample_path).ok();
}
