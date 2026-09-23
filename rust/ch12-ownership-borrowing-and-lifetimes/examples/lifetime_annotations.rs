fn longest_with_announcement<'a, T: std::fmt::Display>(
    x: &'a str,
    y: &'a str,
    ann: T,
) -> &'a str {
    println!("Announcement: {}", ann);
    if x.len() >= y.len() {
        x
    } else {
        y
    }
}

fn main() {
    let s1 = String::from("Rust");
    let result;
    {
        let s2 = String::from("Go");
        result = longest_with_announcement(
            s1.as_str(), s2.as_str(), "Comparing languages!");
        println!("Longest: {}", result);
    }
    // Why can't we print result here?
    // Think about which input has the shorter lifetime.
}
