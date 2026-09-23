fn main() {
    let text = String::from("The quick brown fox jumps over the lazy dog");

    let count = word_count(&text);
    println!("'{}' has {} words", text, count); // text is still valid

    // Calling it again is fine: we only ever borrowed.
    let count2 = word_count(&text);
    println!("Still {} words", count2);
}
