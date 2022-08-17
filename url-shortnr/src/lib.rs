use axum::{routing::{get,post}, 
    Router,
    response::Redirect,
    http::StatusCode,
    extract::Extension,
};



use sync_wrapper::SyncWrapper;

use shuttle_service::{error::CustomError, ShuttleAxum};
//use sync_wrapper::SyncWrapper;
use serde::{Serialize, Deserialize};
use sqlx::postgres::PgPoolOptions;
use sqlx::{FromRow, PgPool};
use url::Url;
use anyhow::Context;
use std::fs;

#[derive(Deserialize, Serialize, FromRow)]
struct StoredURL {
    pub id: String,
    pub url: String,
}

async fn redirect(id: String, Extension(pool): Extension<PgPool>) -> Result<Redirect, StatusCode> {
    let stored_url: StoredURL = sqlx::query_as("SELECT * FROM url WHERE id = $1")
        .bind(id)
        .fetch_one(&pool)
        .await
        .map_err(|err| match err {
            sqlx::Error::RowNotFound => StatusCode::NOT_FOUND,
            _=> StatusCode::INTERNAL_SERVER_ERROR

        })?;

        Ok(Redirect::to(&stored_url.url))
}

async fn shorten(url: String, Extension(pool): Extension<PgPool>) -> Result<String, StatusCode> {
    let id = &nanoid::nanoid!(6);

    let parserd_url = Url::parse(&url).map_err(|_err| {
        StatusCode::UNPROCESSABLE_ENTITY
    })?;

    sqlx::query("INSERT INTO url(id, url) VALUES ($1, $2)")
        .bind(id)
        .bind(parserd_url.as_str())
        .execute(&pool)
        .await
        .map_err(|_| {
            StatusCode::INTERNAL_SERVER_ERROR
        })?;

        Ok(format!("https://url-shrtnr.shuttleapp.rs/{id}"))
}


async fn hello_world() -> &'static str {
    "Hello, world!"
}


#[shuttle_service::main]
async fn axum() -> ShuttleAxum {

    let env = fs::read_to_string(".env").unwrap();
    let (key, database_url) = env.split_once('=').unwrap();

    assert_eq!(key, "DATABASE_URL"); 


    let pgpool = PgPoolOptions::new()
        .max_connections(50)
        .connect(database_url)
        .await
        .context("could not connect to database_url")?;
 

    let router = Router::new()
        .route("/hello", get(hello_world))
        .route("/:id", get(redirect))
        .route("/", post(shorten))
        .layer(Extension(pgpool));


    let sync_wrapper = SyncWrapper::new(router);

    Ok(sync_wrapper)
}