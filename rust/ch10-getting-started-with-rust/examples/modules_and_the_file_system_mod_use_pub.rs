mod math {
    pub fn add(a: i32, b: i32) -> i32 {
        a + b
    }

    fn internal_helper() -> i32 {
        42
    }
}

fn main() {
    println!("{}", math::add(2, 3));
    // math::internal_helper(); // Error: function is private
}
