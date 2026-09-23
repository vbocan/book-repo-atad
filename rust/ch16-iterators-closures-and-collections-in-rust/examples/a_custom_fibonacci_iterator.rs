struct Fibonacci {
    a: u64,
    b: u64,
}

impl Fibonacci {
    fn new() -> Fibonacci {
        Fibonacci { a: 0, b: 1 }
    }
}

impl Iterator for Fibonacci {
    type Item = u64;

    fn next(&mut self) -> Option<u64> {
        let result = self.a;
        let new_b = self.a.checked_add(self.b)?; // Returns None on overflow
        self.a = self.b;
        self.b = new_b;
        Some(result)
    }
}

fn main() {
    // First 20 Fibonacci numbers
    let fibs: Vec<u64> = Fibonacci::new().take(20).collect();
    println!("First 20: {:?}", fibs);

    // Sum of even Fibonacci numbers below 4 million
    let sum: u64 = Fibonacci::new()
        .take_while(|&n| n < 4_000_000)
        .filter(|n| n % 2 == 0)
        .sum();
    println!("Sum of even Fibs below 4M: {}", sum);
    // Expected: 4613732 (Project Euler #2)
}
