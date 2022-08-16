use axum::response::IntoResponse;
//use axum::extract::Path;
use axum::http::StatusCode;

use axum::{Extension, Json};
use sqlx::PgPool;

use crate::{
    models::task
};


pub async fn all_tasks(Extension(pool): Extension<PgPool>) -> impl IntoResponse {
    let sql = "SELECT * FROM task ".to_string();
    //let mut all = pool.acquire().await;

    let task = sqlx::query_as::<_, task::Task>(&sql).fetch_all(&pool).await.unwrap();

    (StatusCode::OK, Json(task))
}


pub async fn new_task(Json(task): Json<task::NewTask>, Extension(pool): Extension<PgPool>) -> impl IntoResponse {
    

    let sql = "INSERT INTO task (task) values ($1)";

    let _ = sqlx::query(&sql)
        .bind(&task.task)
        .execute(&pool)
        .await
        .unwrap();

    (StatusCode::OK, Json(task))
}