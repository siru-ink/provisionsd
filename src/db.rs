use rocket::fairing::AdHoc;
use rocket_db_pools::Database;
use rocket_db_pools::sqlx::{Executor, SqlitePool};

#[derive(Database)]
#[database("maindb")]
pub struct MainDB(SqlitePool);

pub fn database_setup() -> AdHoc {
    AdHoc::try_on_ignite("Database setup.", |rocket| async {
        let pool = MainDB::fetch(&rocket).expect("Database setup failed.");

        pool.execute(
            "
            create table if not exists 'users' (
	        'id' integer not null unique,
	        'user' text not null,
	        'passwd' text not null,
	        primary key('id' autorincrement)
            );",
        )
        .await
        .expect("Database table users should exist or be creatable.");

        pool.execute(
            "
            create table if not exists 'stores' (
            'id' integer not null unique,
            'name' text not null unique,
            primary key('id' autoincrement)
            );",
        )
        .await
        .expect("Database table stores should exist or be creatable.");

        pool.execute(
            "
            create table if not exists 'ingredients' (
            'id' integer not null unique,
            'name' text not null,
            primary key('id' autoincrement)
            );",
        )
        .await
        .expect("Database table ingredients should exist or be creatable.");

        pool.execute(
            "
            create table if not exists 'meal' (
            'id' integer not null unique,
            'name' text not null,
            primary key('id' autoincrement)
            );",
        )
        .await
        .expect("Database table meal should exist or be creatable.");

        Ok(rocket)
    })
}
