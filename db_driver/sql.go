package db_driver




var userTable = `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username varchar(100) DEFAULT "",
    password varchar(100) DEFAULT "",
    first_name varchar(100) DEFAULT "",
    last_name varchar(100) DEFAULT "",
    email varchar(254) DEFAULT "",
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    picture TEXT DEFAULT "static/ava/pomog_ava.png",
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)`

var threadTable = `CREATE TABLE IF NOT EXISTS thread (
    id INTEGER PRIMARY KEY,
    subject TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    userID INTEGER,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE
)`

var postTable = `CREATE TABLE IF NOT EXISTS post (
    id INTEGER PRIMARY KEY,
    subject TEXT,
    content TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    threadID INTEGER,
    userID INTEGER,
    FOREIGN KEY (threadID) REFERENCES thread(id) ON DELETE CASCADE,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE
)`

var votesTable = `CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY,
    upCount INTEGER,
    downCount INTEGER,
    postID INTEGER,
    userID INTEGER,
    FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE
)`

// var fk_user_thread = `
// 	ALTER TABLE thread 	
// 	ADD CONSTRAINT fk_user_thread
// 	FOREIGN KEY (userID) REFERENCES users(id)
// 	ON DELETE CASCADE;`



// var fk_user_post = `
// 	ALTER TABLE post
// 	ADD CONSTRAINT fk_user_post
// 	FOREIGN KEY (userID) REFERENCES users(id)
// 	ON DELETE CASCADE`

// var fk_user_votes = `
// 	ALTER TABLE votes
// 	ADD CONSTRAINT fk_user_votes
// 	FOREIGN KEY (userID) REFERENCES users(id)
// 	ON DELETE CASCADE`

// var fk_post_votes = `
// 	ALTER TABLE votes
// 	ADD CONSTRAINT fk_post_votes
// 	FOREIGN KEY (userID) REFERENCES users(id)
// 	ON DELETE CASCADE`

func getQuerys() []string {
	var sqlQuerys []string
	sqlQuerys = append(sqlQuerys, userTable)
	sqlQuerys = append(sqlQuerys, threadTable)
	sqlQuerys = append(sqlQuerys, postTable)
	sqlQuerys = append(sqlQuerys, votesTable)
	return sqlQuerys

}

// func getFKQuerys() []string {
// 	var sqlFKQuerys []string
// 	sqlFKQuerys = append(sqlFKQuerys, fk_user_thread)
// 	// sqlFKQuerys = append(sqlFKQuerys, fk_user_post)
// 	// sqlFKQuerys = append(sqlFKQuerys, fk_user_votes)
// 	// sqlFKQuerys = append(sqlFKQuerys, fk_post_votes)

// 	return sqlFKQuerys
// }
