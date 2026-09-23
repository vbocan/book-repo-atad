fn main() {
    let x = 5;
    let x = x + 1;    // Shadows the previous x
    let x = x * 2;    // Shadows again
    println!("x = {}", x); // 12
}
