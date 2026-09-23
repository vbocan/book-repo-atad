use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() != 3 {
        eprintln!("Usage: tempconv <value> <unit>");
        eprintln!("Units: C, F, K");
        std::process::exit(1);
    }

    let value: f64 = args[1].parse().unwrap_or_else(|_| {
        eprintln!("Invalid temperature: {}", args[1]);
        std::process::exit(1);
    });

    let unit = args[2].to_uppercase();

    let celsius = match unit.as_str() {
        "C" => value,
        "F" => (value - 32.0) * 5.0 / 9.0,
        "K" => value - 273.15,
        _ => {
            eprintln!("Unknown unit: {}", unit);
            std::process::exit(1);
        }
    };

    let fahrenheit = celsius * 9.0 / 5.0 + 32.0;
    let kelvin = celsius + 273.15;

    println!("{:.2}°C = {:.2}°F = {:.2}K", celsius, fahrenheit, kelvin);
}
