use axum::{Router, extract::State, response::IntoResponse, routing::get};
use std::sync::Arc;

struct AppState {
    db: sqlx::SqlitePool,
}

async fn list_urls(State(state): State<Arc<AppState>>) -> impl IntoResponse {
    // A real handler would query state.db; this one reports the pool size
    format!("{} open connection(s)", state.db.size())
}

#[tokio::main]
async fn main() {
    let pool = sqlx::SqlitePool::connect("sqlite:urlshort.db?mode=rwc").await.unwrap();
    let state = Arc::new(AppState { db: pool });

    let app = Router::new()
        .route("/api/urls", get(list_urls))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
