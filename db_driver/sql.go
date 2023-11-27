package db_driver

var sqlQurys []string

var userTable = `CREATE TABLE IF NOT EXISTS users 
		(id INTEGER PRIMARY KEY, 
		username varchar(100) DEFAULT "", 
		password varchar(100) DEFAULT "",
		first_name varchar(100) DEFAULT "",
		last_name varchar(100) DEFAULT "",
		email varchar(254) DEFAULT "",
		created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		picture  TEXT DEFAULT "static/ava/pomog_ava.png", 
		last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

var threadTable = `CREATE TABLE IF NOT EXISTS thread 
	(id INTEGER PRIMARY KEY, 
	subject TEXT,
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	userID INTEGER)`

var postTable = `CREATE TABLE IF NOT EXISTS post 
	(id INTEGER PRIMARY KEY, 
	subject TEXT,
	content TEXT,
	created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	threadID INTEGER,
	userID INTEGER)`

var votesTable = `CREATE TABLE IF NOT EXISTS post 
	(id INTEGER PRIMARY KEY, 
	upCount INTEGER,
	downCount INTEGER,
	postID INTEGER)`

func getStatements() []string {
	sqlQurys = append(sqlQurys, userTable)
	sqlQurys = append(sqlQurys, threadTable)
	sqlQurys = append(sqlQurys, postTable)
	sqlQurys = append(sqlQurys, votesTable)

	return sqlQurys
}
