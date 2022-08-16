use axum::{
    extract::{Extension},routing::{get, post}, Router,
};

use sqlx::postgres::PgPoolOptions;
use std::net::SocketAddr;
use std::fs;
use anyhow::Context;
mod models;
mod controllers;

#[tokio::main]
async fn main() -> anyhow::Result<()> {

    let env = fs::read_to_string(".env").unwrap();
    let (key, database_url) = env.split_once('=').unwrap();


    assert_eq!(key, "DATABASE_URL");

    tracing_subscriber::fmt::init();

    let pool = PgPoolOptions::new()
    .max_connections(50)
    .connect(&database_url)
    .await
    .context("could not connect to database_url")?;

    let app = Router::new()
        .route("/hello", get(hello))
        .route("/tasks", get(controllers::task::all_tasks))
        .route("/task", post(controllers::task::new_task))
        .layer(Extension(pool));

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    tracing::debug!("Listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await?;

        Ok(())

}

async fn hello() -> &'static str {
    "Hello, World!"
}