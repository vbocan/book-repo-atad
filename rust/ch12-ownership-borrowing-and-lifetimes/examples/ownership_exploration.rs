fn main() {
    // Scenario 1: Move semantics
    let s1 = String::from("hello");
    let s2 = s1;
    // Does this print? Why or why not?
    // println!("{}", s1);

    // Scenario 2: Copy semantics
    let x = 42;
    let y = x;
    println!("x={}, y={}", x, y); // Does this work?

    // Scenario 3: Function ownership
    let s3 = String::from("world");
    takes_ownership(s3);
    // Does this print? Why or why not?
    // println!("{}", s3);

    // Scenario 4: Returning ownership
    let s4 = gives_ownership();
    println!("{}", s4); // Does this work?

    // Scenario 5: Taking and giving back
    let s5 = String::from("boomerang");
    let s5 = takes_and_gives_back(s5);
    println!("{}", s5); // Does this work?
}

fn takes_ownership(s: String) {
    println!("Took: {}", s);
}

fn gives_ownership() -> String {
    String::from("gift")
}

fn takes_and_gives_back(s: String) -> String {
    s
}
