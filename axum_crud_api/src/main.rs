use axum::{
    routing::{get}, Router,
};

use std::net::SocketAddr;

#[tokio::main]
async fn main() {

    tracing_subscriber::fmt::init();

    let app = Router::new()
        .route("/hello", get(hello));

    let addr = SocketAddr::from(([127, 0, 0, 1], 3000));
    tracing::debug!("Listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();

}

async fn hello() -> &'static str {
    "Hello, World!"
}