fn take_ownership(s: String) {
    println!("Got: {}", s);
} // s is dropped here

fn main() {
    let greeting = String::from("hello");
    take_ownership(greeting);

    // println!("{}", greeting); // Error: greeting was moved
}
