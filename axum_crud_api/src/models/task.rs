use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize, sqlx::FromRow)]
pub struct Task {
    pub id: i32,
    pub task: String,
}

#[derive(Deserialize, Serialize, sqlx::FromRow)]
pub struct NewTask {
    pub task: String,
}