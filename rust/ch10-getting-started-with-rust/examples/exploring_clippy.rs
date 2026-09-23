fn main() {
    // 1. Manual Option mapping (clippy suggests .map())
    let x: Option<i32> = Some(5);
    let _y = match x {
        Some(v) => Some(v * 2),
        None => None,
    };

    // 2. Single-arm match (clippy suggests if let)
    let opt: Option<i32> = Some(7);
    match opt {
        Some(v) => println!("{}", v),
        _ => {}
    }

    // 3. Using .len() == 0 instead of .is_empty()
    let v: Vec<i32> = vec![];
    if v.len() == 0 {
        println!("empty");
    }

    // 4. Unnecessary return keyword
    let _result = add(2, 3);
}

fn add(a: i32, b: i32) -> i32 {
    return a + b;
}
