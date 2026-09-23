//! Exercise 13.1 -- Library Catalog
//!
//! Models a small library catalog with structs, methods, and encapsulation.
//! `Book` keeps its four fields private and lets the rest of the program
//! reach them only through getters (`title`, `author`, `isbn`,
//! `is_available`) and the two methods allowed to change availability
//! (`check_out`, `return_book`) -- the same pattern Section 13.4.1 uses for
//! `User`'s `password_hash`. `Library` wraps a `Vec<Book>` behind the same
//! kind of narrow API: adding books, searching by title or author, listing
//! what is available, and checking a book out by ISBN, which reports
//! failure through a `Result` instead of panicking. Both types implement
//! `std::fmt::Display` for human-readable printing, the way Section 13.3.2
//! implements it for `Point`.
//!
//! Run with `cargo run --example library_catalog`.

use std::fmt;

/// A single catalog entry. All four fields are private; the rest of the
/// program reads them through getters and changes availability only
/// through `check_out` and `return_book`, never by writing a field
/// directly.
#[derive(Debug)]
struct Book {
    title: String,
    author: String,
    isbn: String,
    available: bool,
}

impl Book {
    /// Creates a new book. Every book starts out available.
    fn new(title: String, author: String, isbn: String) -> Book {
        Book {
            title,
            author,
            isbn,
            available: true,
        }
    }

    fn title(&self) -> &str {
        &self.title
    }

    fn author(&self) -> &str {
        &self.author
    }

    fn isbn(&self) -> &str {
        &self.isbn
    }

    fn is_available(&self) -> bool {
        self.available
    }

    /// Marks the book as checked out. Callers decide first whether that is
    /// legal (see `Library::check_out_by_isbn`); this method just mutates.
    fn check_out(&mut self) {
        self.available = false;
    }

    /// Marks the book as available again.
    fn return_book(&mut self) {
        self.available = true;
    }
}

impl fmt::Display for Book {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let status = if self.available {
            "available"
        } else {
            "checked out"
        };
        write!(
            f,
            "\"{}\" by {} (ISBN {}) - {status}",
            self.title, self.author, self.isbn
        )
    }
}

/// A collection of books. `Library` is the only way client code touches the
/// catalog: the `Vec<Book>` itself is private, so every addition, search,
/// or checkout goes through a method that can enforce the catalog's rules.
#[derive(Debug)]
struct Library {
    books: Vec<Book>,
}

impl Library {
    fn new() -> Library {
        Library { books: Vec::new() }
    }

    fn add_book(&mut self, book: Book) {
        self.books.push(book);
    }

    /// Books whose title contains `query`, case-insensitively.
    fn search_by_title(&self, query: &str) -> Vec<&Book> {
        let needle = query.to_lowercase();
        self.books
            .iter()
            .filter(|book| book.title().to_lowercase().contains(needle.as_str()))
            .collect()
    }

    /// Books whose author contains `query`, case-insensitively.
    fn search_by_author(&self, query: &str) -> Vec<&Book> {
        let needle = query.to_lowercase();
        self.books
            .iter()
            .filter(|book| book.author().to_lowercase().contains(needle.as_str()))
            .collect()
    }

    fn list_available(&self) -> Vec<&Book> {
        self.books
            .iter()
            .filter(|book| book.is_available())
            .collect()
    }

    /// Checks a book out by ISBN. Fails with an explanatory `Err` if no
    /// book carries that ISBN, or if the book is already checked out;
    /// otherwise flips its availability and returns `Ok(())`.
    fn check_out_by_isbn(&mut self, isbn: &str) -> Result<(), String> {
        match self.books.iter_mut().find(|book| book.isbn() == isbn) {
            Some(book) if book.is_available() => {
                book.check_out();
                Ok(())
            }
            Some(book) => Err(format!("\"{}\" is already checked out", book.title())),
            None => Err(format!("no book found with ISBN {isbn}")),
        }
    }
}

impl fmt::Display for Library {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        writeln!(f, "Library catalog ({} book(s)):", self.books.len())?;
        let lines: Vec<String> = self.books.iter().map(|book| format!("  {book}")).collect();
        write!(f, "{}", lines.join("\n"))
    }
}

fn main() {
    // `Book` on its own, showing `check_out`/`return_book` mutating
    // availability independent of any `Library`, and Display/Debug
    // printing side by side.
    let mut spare_copy = Book::new(
        "Spare Copy".into(),
        "Demo Author".into(),
        "000-0000000001".into(),
    );
    println!("Display: {spare_copy}");
    println!("Debug:   {spare_copy:?}");
    spare_copy.check_out();
    println!("After check_out():   {spare_copy}");
    spare_copy.return_book();
    println!("After return_book(): {spare_copy}");
    println!();

    // Build a small catalog.
    let mut library = Library::new();
    library.add_book(Book::new(
        "The Pragmatic Programmer".into(),
        "David Thomas and Andrew Hunt".into(),
        "978-0135957059".into(),
    ));
    library.add_book(Book::new(
        "Clean Code".into(),
        "Robert C. Martin".into(),
        "978-0132350884".into(),
    ));
    library.add_book(Book::new(
        "The Rust Programming Language".into(),
        "Steve Klabnik and Carol Nichols".into(),
        "978-1718503106".into(),
    ));
    library.add_book(Book::new(
        "Design Patterns".into(),
        "Erich Gamma, Richard Helm, Ralph Johnson, and John Vlissides".into(),
        "978-0201633610".into(),
    ));
    library.add_book(Book::new(
        "Effective Rust".into(),
        "David Drysdale".into(),
        "978-1098151402".into(),
    ));

    println!("{library}");
    println!();

    println!("Search for \"rust\" in titles:");
    for book in library.search_by_title("rust") {
        println!("  {book}");
    }
    println!();

    println!("Search for \"Martin\" in authors:");
    for book in library.search_by_author("Martin") {
        println!("  {book}");
    }
    println!();

    if let Some(book) = library.search_by_title("Effective Rust").first() {
        println!(
            "Details via getters: \"{}\" by {}, ISBN {}, available: {}",
            book.title(),
            book.author(),
            book.isbn(),
            book.is_available()
        );
        println!();
    }

    // Check one book out.
    match library.check_out_by_isbn("978-0132350884") {
        Ok(()) => println!("Checked out \"Clean Code\"."),
        Err(e) => println!("Could not check out: {e}"),
    }
    println!();

    println!("Available titles after checkout:");
    for book in library.list_available() {
        println!("  {}", book.title());
    }
    println!();

    // Error path: no book carries this ISBN.
    match library.check_out_by_isbn("000-0000000000") {
        Ok(()) => println!("Unexpected success."),
        Err(e) => println!("Error checking out an unknown ISBN: {e}"),
    }

    // Error path: the book is already checked out.
    match library.check_out_by_isbn("978-0132350884") {
        Ok(()) => println!("Unexpected success."),
        Err(e) => println!("Error checking out an already-checked-out book: {e}"),
    }
    println!();

    println!("Final catalog state:");
    println!("{library}");
}
