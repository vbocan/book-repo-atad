use std::collections::HashMap;
use std::fs;

fn analyze(text: &str) -> Vec<(String, usize)> {
    let mut freq: HashMap<String, usize> = HashMap::new();

    text.split_whitespace()
        .map(|w| w.trim_matches(|c: char| !c.is_alphanumeric()).to_lowercase())
        .filter(|w| !w.is_empty())
        .for_each(|w| {
            *freq.entry(w).or_insert(0) += 1;
        });

    let mut sorted: Vec<(String, usize)> = freq.into_iter().collect();
    sorted.sort_by(|a, b| b.1.cmp(&a.1).then(a.0.cmp(&b.0)));
    sorted
}

fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.len() != 2 {
        eprintln!("Usage: wordfreq <file>");
        std::process::exit(1);
    }

    let text = fs::read_to_string(&args[1]).expect("Cannot read file");
    let freq = analyze(&text);

    println!("{:<20} {}", "WORD", "COUNT");
    println!("{}", "-".repeat(30));
    for (word, count) in freq.iter().take(20) {
        println!("{:<20} {}", word, count);
    }
    println!("\nTotal unique words: {}", freq.len());
}
