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

        pool.execute(
            "
            create table if not exists 'units' (
            'id' integer not null unique,
            'name' text not null unique,
            'abbr' text not null unique,
            primary key('id' autoincrement)
            );",
        )
        .await
        .expect("Database table units should exist or be creatable.");

        pool.execute(
            "
            create table if not exists 'meal_ingredients' (
            'mealid' integer not null,
            'ingredientid' integer not null,
            'quantity' real not null,
            'unit' id not null,
            foreign key('mealid') references 'meal'('id') on delete cascade,
            foreign key('ingredientid') references 'ingredients'('id') on delete cascade,
            foreign key('unit') references 'units'('id') on delete cascade,
            primary key('meal_id', 'ingredient_id')
            );",
        )
        .await
        .expect("Database table meal_ingredients should exist or be creatable.");

        Ok(rocket)
    })
}
