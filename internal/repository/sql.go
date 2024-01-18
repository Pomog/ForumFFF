package repository

var userTable = `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username varchar(100) DEFAULT "",
    password varchar(100) DEFAULT "",
    first_name varchar(100) DEFAULT "",
    last_name varchar(100) DEFAULT "",
    email varchar(254) DEFAULT "",
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    picture TEXT DEFAULT "static/ava/ava1.png",
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

var threadTable = `CREATE TABLE IF NOT EXISTS thread (
    id INTEGER PRIMARY KEY,
    subject TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    userID INTEGER,
    threadImage TEXT DEFAULT "",
    category TEXT DEFAULT "",
    FOREIGN KEY (userID) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);`

var postTable = `CREATE TABLE IF NOT EXISTS post (
    id INTEGER PRIMARY KEY,
    subject TEXT,
    content TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    threadID INTEGER,
    userID INTEGER,
    postImage TEXT DEFAULT "",
    FOREIGN KEY (threadID) REFERENCES thread(id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);`

var votesTable = `CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY,
    like BOOLEAN DEFAULT FALSE,
    dislike BOOLEAN DEFAULT FALSE,
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

var privateMessageTable = `CREATE TABLE IF NOT EXISTS pm (
    id INTEGER PRIMARY KEY,
    content TEXT,
    created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    senderUserID INTEGER,
    receiverUserID INTEGER,
    FOREIGN KEY (senderUserID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
    FOREIGN KEY (receiverUserID) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);`

var guestUser = `INSERT INTO users (username, password, first_name, last_name, email, type)
VALUES ('guest', '123456', 'Guest', 'User', 'guest@gmail.com', 'guest');`

var addClassificationToPost = `ALTER TABLE post
ADD COLUMN classification VARCHAR(50) DEFAULT 'unsorted';`

var addClassificationToThread = `ALTER TABLE thread
ADD COLUMN classification VARCHAR(50) DEFAULT 'unsorted';`

var addUserType = `ALTER TABLE users
ADD COLUMN type VARCHAR(50) DEFAULT 'user';`

func getQuerys() []string {
	var sqlQuerys []string
	sqlQuerys = append(sqlQuerys, userTable)
	sqlQuerys = append(sqlQuerys, threadTable)
	sqlQuerys = append(sqlQuerys, postTable)
	sqlQuerys = append(sqlQuerys, votesTable)
	sqlQuerys = append(sqlQuerys, sessionIdTable)
	sqlQuerys = append(sqlQuerys, privateMessageTable)

	sqlQuerys = append(sqlQuerys, addClassificationToPost)
	sqlQuerys = append(sqlQuerys, addClassificationToThread)
	sqlQuerys = append(sqlQuerys, addUserType)

	return sqlQuerys
}
