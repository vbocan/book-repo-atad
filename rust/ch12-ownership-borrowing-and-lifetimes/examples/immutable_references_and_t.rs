fn calculate_length(s: &String) -> usize {
    s.len()
} // s goes out of scope, but it doesn't own the String, so nothing happens

fn main() {
    let s = String::from("hello");
    let len = calculate_length(&s); // Borrow s
    println!("'{}' has length {}", s, len); // s is still valid
}
