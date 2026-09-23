# 🦀🐹 Advanced Techniques for Application Development — Code Examples

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![Rust](https://img.shields.io/badge/Rust-2024_edition-CE422B?logo=rust&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-bundled-003B57?logo=sqlite&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-2F855A)
![CI](https://img.shields.io/badge/build-via_GitHub_Actions-2088FF?logo=githubactions&logoColor=white)

Runnable companion code for the textbook **_Advanced Techniques for Application
Development_** by Valer Bocan (Editura Politehnica Timișoara). Every program here is
extracted from the book and **verified to compile and run** with the toolchains the book
targets, so you can clone, run, and tinker as you read.

```bash
git clone https://github.com/vbocan/book-repo-atad.git
```

## 📦 What's inside

| | |
|---|---|
| 🐹 [`go/`](go/) | Go examples — Part A, chapters 2–9 (one Go module) |
| 🦀 [`rust/`](rust/) | Rust examples — Part B, chapters 10–19 (a Cargo workspace) |

Each chapter folder is named after the book chapter; each example folder is named after
the listing it comes from.

## 🚀 Running an example

**Go** (needs Go 1.24+):
```bash
cd go
go run ./ch03-types-variables-and-control-flow-in-go/temperature-converter
```

**Rust** (needs a 2024-edition toolchain, Rust 1.85+):
```bash
cd rust
cargo run -p ch12-ownership-borrowing-and-lifetimes --example move-semantics
```

## 📚 Chapter map

| Ch. | Topic | Go examples | Rust examples |
|----:|-------|:-----------:|:-------------:|
| 01 | The Landscape Of Modern Systems Programming | 1 | 1 |
| 02 | Getting Started With Go | 3 | — |
| 03 | Types Variables And Control Flow In Go | 2 | — |
| 04 | Functions Methods And Error Handling In Go | 1 | — |
| 10 | Getting Started With Rust | — | 5 |
| 11 | Types Variables And Control Flow In Rust | — | 6 |
| 12 | Ownership Borrowing And Lifetimes | — | 8 |
| 16 | Iterators Closures And Collections In Rust | — | 2 |
| 19 | Building Http Services In Rust | — | 1 |

## 🧪 How these stay correct

The code is generated from the manuscript and checked in CI (`go build`/`go vet` and
`cargo build --examples`). See [`.github/workflows/ci.yml`](.github/workflows/ci.yml).

## ⚖️ License

Source code in this repository is released under the [MIT License](LICENSE). The book text
itself is © Valer Bocan and is **not** covered by this license.
