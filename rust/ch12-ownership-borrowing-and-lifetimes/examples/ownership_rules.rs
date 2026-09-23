fn main() {
    let s = String::from("hello"); // s owns the String
    println!("{}", s);
} // s goes out of scope — the String is dropped, memory freed
