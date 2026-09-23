use std::fs;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let contents = fs::read_to_string("config.txt")?;
    println!("Config: {}", contents);
    Ok(())
}
