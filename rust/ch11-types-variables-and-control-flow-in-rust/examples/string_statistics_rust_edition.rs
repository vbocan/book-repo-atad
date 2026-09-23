use std::env;

fn main() {
    let args: Vec<String> = env::args().skip(1).collect();
    if args.is_empty() {
        eprintln!("Usage: strstat <text>");
        std::process::exit(1);
    }

    let text = args.join(" ");

    let bytes = text.len();
    let chars = text.chars().count();
    let words = text.split_whitespace().count();
    let upper = text.chars().filter(|c| c.is_uppercase()).count();
    let lower = text.chars().filter(|c| c.is_lowercase()).count();
    let digits = text.chars().filter(|c| c.is_ascii_digit()).count();

    println!("Text:       {:?}", text);
    println!("Bytes:      {}", bytes);
    println!("Characters: {}", chars);
    println!("Words:      {}", words);
    println!("Uppercase:  {}", upper);
    println!("Lowercase:  {}", lower);
    println!("Digits:     {}", digits);
}
