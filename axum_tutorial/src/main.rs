use axum::{
    routing::{get},
    Router,
};

use std::net::SocketAddr;

#[tokio::main]
async fn main() {


    let app = Router::new()
        .route("/", get(root));

    let addr = SocketAddr::from(([127, 0 , 0, 1], 8000));
    println!("listening on {}", addr);
    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();

}

async fn root() -> &'static str {
    "Hello, World!"
}
