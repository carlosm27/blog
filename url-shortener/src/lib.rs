
use axum::{
    routing::{get, post},
    Router,
    response::Redirect,
    http::StatusCode,
    extract::Extension,
};
use shuttle_service::{error::CustomError, ShuttleAxum};
use sync_wrapper:: SyncWrapper;
use serde::Serialize;
use sqlx::migrate::Migrator;
use sqlx::{FromRow, PgPool};
use url::Url;



#[derive(Serialize, FromRow)]
struct StoredURL {
    pub id: String,
    pub url: String,
}

async fn redirect(id: String, Extension(pool): Extension<PgPool>) -> Result<Redirect, StatusCode> {
    let stored_url: StoredURL = sqlx::query_as("SELECT * FROM urls WHERE id = $1")
        .bind(id)
        .fetch_one(&pool)
        .await
        .map_err(|err| match err {
            sqlx::Error::RowNotFound => StatusCode::NOT_FOUND,
            _=> StatusCode::INTERNAL_SERVER_ERROR

        })?;

        Ok(Redirect::to(&stored_url.url))
}

async fn shorten(url:String, Extension(pool): Extension<PgPool>) -> Result<String, StatusCode> {
    let id = &nanoid::nanoid!(6);

    let parserd_url = Url::parse(&url).map_err(|_err| {
        StatusCode::UNPROCESSABLE_ENTITY
    })?;

    sqlx::query("INSERT INTO urls(id, url) VALUES ($1, $2)")
        .bind(id)
        .bind(parserd_url.as_str())
        .execute(&pool)
        .await
        .map_err(|_| {
            StatusCode::INTERNAL_SERVER_ERROR
        })?;

        Ok(format!("https://s.shuttleapp.rs/{id}"))
}

static MIGRATOR: Migrator = sqlx::migrate!();

#[shuttle_service::main]
async fn axum(pgpool: PgPool) -> ShuttleAxum {
    
    MIGRATOR.run(&pgpool).await.map_err(CustomError::new)?;

    let router = Router::new()
        .route("/", get(redirect))
        .route("/", post(shorten))
        .layer(Extension(pgpool));

    let sync_wrapper = SyncWrapper::new(router);

    Ok(sync_wrapper)

}






