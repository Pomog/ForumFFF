package repository

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
);`

var threadTable = `CREATE TABLE IF NOT EXISTS thread (
    id INTEGER PRIMARY KEY,
    subject TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    userID INTEGER,
    FOREIGN KEY (userID) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);`

var postTable = `CREATE TABLE IF NOT EXISTS post (
    id INTEGER PRIMARY KEY,
    subject TEXT,
    content TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    threadID INTEGER,
    userID INTEGER,
    FOREIGN KEY (threadID) REFERENCES thread(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);`

var votesTable = `CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY,
    upCount INTEGER,
    downCount INTEGER,
    postID INTEGER,
    userID INTEGER,
    FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);`

var sessionIdTable = `CREATE TABLE IF NOT EXISTS sessionId (
    id INTEGER PRIMARY KEY,
    userID INTEGER,
    sessionID TEXT,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);`

var guestUser = `INSERT INTO users (username, password, first_name, last_name, email)
VALUES ('guest', '123456', 'Guest', 'User', 'guest@gmail.com');
);`

func getQuerys() []string {
	var sqlQuerys []string
	sqlQuerys = append(sqlQuerys, userTable)
	sqlQuerys = append(sqlQuerys, threadTable)
	sqlQuerys = append(sqlQuerys, postTable)
	sqlQuerys = append(sqlQuerys, votesTable)
	sqlQuerys = append(sqlQuerys, sessionIdTable)
	return sqlQuerys
}
