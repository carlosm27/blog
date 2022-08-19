use axum::{
    extract::{Extension},routing::{get, post, put, delete}, Router,
};

use sqlx::postgres::PgPoolOptions;
use std::net::SocketAddr;
use std::fs;
use anyhow::Context;
use tower_http::trace::TraceLayer;
use tracing_subscriber::{layer::SubscriberExt, util::SubscriberInitExt};

mod errors;
mod models;
mod controllers;

#[tokio::main]
async fn main() -> anyhow::Result<()> {

    let env = fs::read_to_string(".env").unwrap();
    let (key, database_url) = env.split_once('=').unwrap();


    assert_eq!(key, "DATABASE_URL");

    tracing_subscriber::registry()
        .with(tracing_subscriber::EnvFilter::new(
            std::env::var("tower_http=trace")
                .unwrap_or_else(|_| "example_tracing_aka_logging=debug,tower_http=debug".into()),
        ))
        .with(tracing_subscriber::fmt::layer())
        .init();

    let pool = PgPoolOptions::new()
    .max_connections(50)
    .connect(&database_url)
    .await
    .context("could not connect to database_url")?;

    let app = Router::new()
        .route("/hello", get(root))
        .route("/tasks", get(controllers::task::all_tasks))
        .route("/task", post(controllers::task::new_task))
        .route("/task/:id",get(controllers::task::task))
        .route("/task/:id", put(controllers::task::update_task))
        .route("/task/:id", delete(controllers::task::delete_task))
        .layer(Extension(pool))
        .layer(TraceLayer::new_for_http());

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    tracing::debug!("Listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await?;

        Ok(())

}

async fn root() -> &'static str {
    "Hello, World!"
}