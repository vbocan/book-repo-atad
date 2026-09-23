use axum::{
    extract::{Path, State},
    http::StatusCode,
    response::{IntoResponse, Redirect, Response},
    routing::{get, post},
    Json, Router,
};
use serde::{Deserialize, Serialize};
use sqlx::sqlite::SqlitePool;
use std::{
    collections::hash_map::DefaultHasher,
    hash::{Hash, Hasher},
    sync::{
        atomic::{AtomicU64, Ordering},
        Arc,
    },
    time::{SystemTime, UNIX_EPOCH},
};
use tower_http::trace::TraceLayer;
use tracing::{error, info, instrument};

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
struct UrlEntry {
    id: i64,
    short_code: String,
    original_url: String,
    clicks: i64,
    created_at: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct ShortenRequest {
    #[serde(rename = "url")]
    original_url: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    custom_code: Option<String>,
}

#[derive(Serialize)]
struct ShortenResponse {
    short_code: String,
    short_url: String,
}

#[derive(Serialize)]
struct ErrorResponse {
    error: String,
}

enum AppError {
    NotFound,
    Database(String),
}

impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            AppError::NotFound => (StatusCode::NOT_FOUND, "short code not found".to_string()),
            AppError::Database(msg) => (StatusCode::INTERNAL_SERVER_ERROR, msg),
        };
        (status, Json(ErrorResponse { error: message })).into_response()
    }
}

impl From<sqlx::Error> for AppError {
    fn from(err: sqlx::Error) -> Self {
        match err {
            sqlx::Error::RowNotFound => AppError::NotFound,
            other => {
                error!(error = %other, "database error");
                AppError::Database(other.to_string())
            }
        }
    }
}

struct AppState {
    db: SqlitePool,
}

// generate_code returns an 8-character hex short code, e.g. "a1b2c3d4".
// A process-wide counter is folded in alongside the current time so that
// two codes requested in the same nanosecond still differ.
static CODE_COUNTER: AtomicU64 = AtomicU64::new(0);

fn generate_code() -> String {
    let nanos = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_nanos();
    let count = CODE_COUNTER.fetch_add(1, Ordering::Relaxed);

    let mut hasher = DefaultHasher::new();
    nanos.hash(&mut hasher);
    count.hash(&mut hasher);

    format!("{:08x}", hasher.finish() as u32)
}

#[instrument(skip(state))]
async fn handle_shorten(
    State(state): State<Arc<AppState>>,
    Json(req): Json<ShortenRequest>,
) -> Result<(StatusCode, Json<ShortenResponse>), AppError> {
    let code = req.custom_code.unwrap_or_else(generate_code);

    sqlx::query("INSERT INTO urls (short_code, original_url) VALUES (?, ?)")
        .bind(&code)
        .bind(&req.original_url)
        .execute(&state.db)
        .await?;

    info!(short_code = %code, "URL shortened successfully");

    Ok((
        StatusCode::CREATED,
        Json(ShortenResponse {
            short_url: format!("http://localhost:8080/{code}"),
            short_code: code,
        }),
    ))
}

async fn handle_redirect(
    State(state): State<Arc<AppState>>,
    Path(code): Path<String>,
) -> Result<Redirect, AppError> {
    let original_url: (String,) =
        sqlx::query_as("SELECT original_url FROM urls WHERE short_code = ?")
            .bind(&code)
            .fetch_optional(&state.db)
            .await?
            .ok_or(AppError::NotFound)?;

    // Fire-and-forget: don't make the redirect wait on the click count update.
    let db = state.db.clone();
    let click_code = code.clone();
    tokio::spawn(async move {
        if let Err(e) = sqlx::query("UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?")
            .bind(&click_code)
            .execute(&db)
            .await
        {
            error!(short_code = %click_code, error = %e, "failed to increment clicks");
        }
    });

    Ok(Redirect::temporary(&original_url.0))
}

async fn handle_stats(
    State(state): State<Arc<AppState>>,
    Path(code): Path<String>,
) -> Result<Json<UrlEntry>, AppError> {
    let entry = sqlx::query_as::<_, UrlEntry>(
        "SELECT id, short_code, original_url, clicks, created_at FROM urls WHERE short_code = ?",
    )
    .bind(&code)
    .fetch_optional(&state.db)
    .await?
    .ok_or(AppError::NotFound)?;

    Ok(Json(entry))
}

async fn handle_list(
    State(state): State<Arc<AppState>>,
) -> Result<Json<Vec<UrlEntry>>, AppError> {
    let entries = sqlx::query_as::<_, UrlEntry>(
        "SELECT id, short_code, original_url, clicks, created_at FROM urls ORDER BY created_at DESC",
    )
    .fetch_all(&state.db)
    .await?;

    Ok(Json(entries))
}

fn init_logging() {
    tracing_subscriber::fmt()
        .with_env_filter(tracing_subscriber::EnvFilter::try_from_default_env()
            .unwrap_or_else(|_| tracing_subscriber::EnvFilter::new("urlshortener=info")))
        .with_target(false)
        .init();
}

#[tokio::main]
async fn main() {
    init_logging();

    let pool = SqlitePool::connect("sqlite:urlshort.db?mode=rwc")
        .await
        .expect("failed to connect to database");

    sqlx::query(
        "CREATE TABLE IF NOT EXISTS urls (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            short_code TEXT UNIQUE NOT NULL,
            original_url TEXT NOT NULL,
            clicks INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )",
    )
    .execute(&pool)
    .await
    .expect("failed to create schema");

    let state = Arc::new(AppState { db: pool });

    let app = Router::new()
        .route("/api/shorten", post(handle_shorten))
        .route("/api/stats/{code}", get(handle_stats))
        .route("/api/urls", get(handle_list))
        .route("/{code}", get(handle_redirect))
        .layer(TraceLayer::new_for_http())
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:8080")
        .await
        .expect("failed to bind port 8080");
    info!("Listening on http://localhost:8080");
    axum::serve(listener, app).await.expect("server error");
}
