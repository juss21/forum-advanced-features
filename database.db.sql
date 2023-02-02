BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "session" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"key"	TEXT UNIQUE,
	"userId"	INTEGER UNIQUE
);
CREATE TABLE IF NOT EXISTS "posts" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"userId"	INTEGER,
	"title"	TEXT,
	"content"	TEXT,
	"date"	TEXT
);
CREATE TABLE IF NOT EXISTS "commentLikes" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"name"	TEXT,
	"userId"	INTEGER,
	"commentId"	INTEGER,
	UNIQUE("commentId","userId")
);
CREATE TABLE IF NOT EXISTS "postlikes" (
	"id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	"name"	TEXT,
	"userId"	INTEGER,
	"postId"	INTEGER,
	UNIQUE("postId","userId")
);
CREATE TABLE IF NOT EXISTS "users" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"username"	TEXT,
	"password"	TEXT,
	"email"	TEXT
);
CREATE TABLE IF NOT EXISTS "comments" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	"userId"	INTEGER,
	"content"	TEXT,
	"postId"	INTEGER
);
INSERT INTO "session" VALUES (1,'dsad',12);
INSERT INTO "session" VALUES (2,'dddfdssad',1);
INSERT INTO "session" VALUES (3,'0d8bc291-5743-4a90-be03-eb7162c90db2',3);
INSERT INTO "posts" VALUES (66,1,'hahaha','dsfsadf',NULL);
INSERT INTO "posts" VALUES (67,1,'fgfdg','dfgdfg',NULL);
INSERT INTO "posts" VALUES (68,1,'katsetus','Tere siin minu tekst',NULL);
INSERT INTO "posts" VALUES (69,1,'teine katsetus','eelmine pani kelbast',NULL);
INSERT INTO "posts" VALUES (70,1,'kOlm on kohtu seadus','eaeaeaeaeaeaeaeaeaea',NULL);
INSERT INTO "posts" VALUES (71,1,'','',NULL);
INSERT INTO "posts" VALUES (72,4,'uus kasutaja uus postitus!!','jehuhusageajehuhusageajehuhusagea',NULL);
INSERT INTO "posts" VALUES (73,1,'uus post','uus sisu ja uus sisu',NULL);
INSERT INTO "commentLikes" VALUES (96,'dislike',1,114);
INSERT INTO "commentLikes" VALUES (98,'like',1,126);
INSERT INTO "commentLikes" VALUES (107,'dislike',1,127);
INSERT INTO "commentLikes" VALUES (110,'dislike',1,128);
INSERT INTO "commentLikes" VALUES (131,'like',1,131);
INSERT INTO "commentLikes" VALUES (132,'dislike',1,132);
INSERT INTO "commentLikes" VALUES (133,'like',1,133);
INSERT INTO "commentLikes" VALUES (134,'dislike',1,135);
INSERT INTO "commentLikes" VALUES (135,'like',4,139);
INSERT INTO "commentLikes" VALUES (136,'like',1,140);
INSERT INTO "commentLikes" VALUES (137,'like',1,130);
INSERT INTO "commentLikes" VALUES (138,'like',1,137);
INSERT INTO "postlikes" VALUES (71,'dislike',1,0);
INSERT INTO "postlikes" VALUES (82,'like',1,61);
INSERT INTO "postlikes" VALUES (83,'dislike',1,62);
INSERT INTO "postlikes" VALUES (109,'like',1,63);
INSERT INTO "postlikes" VALUES (111,'dislike',1,64);
INSERT INTO "postlikes" VALUES (113,'like',1,67);
INSERT INTO "postlikes" VALUES (124,'dislike',1,66);
INSERT INTO "postlikes" VALUES (125,'like',4,72);
INSERT INTO "postlikes" VALUES (126,'dislike',1,73);
INSERT INTO "postlikes" VALUES (128,'like',1,70);
INSERT INTO "users" VALUES (1,'esimene','teine','first.last@mail.ee');
INSERT INTO "users" VALUES (3,'kolmas','neljas','adsfasdfadsfa');
INSERT INTO "users" VALUES (4,'uuskasutaja3','!Tere123','uuskasutaja3@rrr.rrr');
INSERT INTO "comments" VALUES (108,1,'dsf',0);
INSERT INTO "comments" VALUES (109,1,'',0);
INSERT INTO "comments" VALUES (110,1,'',59);
INSERT INTO "comments" VALUES (111,1,'',59);
INSERT INTO "comments" VALUES (112,1,'dsfdsf',59);
INSERT INTO "comments" VALUES (113,1,'AAAAAAAAAAAAAAAA',59);
INSERT INTO "comments" VALUES (114,1,'fdsfs',58);
INSERT INTO "comments" VALUES (115,1,'dgdsf',58);
INSERT INTO "comments" VALUES (116,1,'',58);
INSERT INTO "comments" VALUES (117,1,'',58);
INSERT INTO "comments" VALUES (118,1,'',58);
INSERT INTO "comments" VALUES (119,1,'',58);
INSERT INTO "comments" VALUES (120,1,'',58);
INSERT INTO "comments" VALUES (121,1,'',58);
INSERT INTO "comments" VALUES (122,1,'',58);
INSERT INTO "comments" VALUES (123,1,'HELLO THERE',58);
INSERT INTO "comments" VALUES (124,1,'fdsfd',60);
INSERT INTO "comments" VALUES (125,1,'fdsf',61);
INSERT INTO "comments" VALUES (126,1,'Väga hea',62);
INSERT INTO "comments" VALUES (127,1,'sdfadsf',63);
INSERT INTO "comments" VALUES (128,1,'jk',64);
INSERT INTO "comments" VALUES (129,1,'fdgfdg',67);
INSERT INTO "comments" VALUES (130,1,'fdsfdsf',66);
INSERT INTO "comments" VALUES (131,1,'jah',66);
INSERT INTO "comments" VALUES (132,1,'katse',66);
INSERT INTO "comments" VALUES (133,1,'test',66);
INSERT INTO "comments" VALUES (134,1,'testeter',0);
INSERT INTO "comments" VALUES (135,1,'jahhh',66);
INSERT INTO "comments" VALUES (136,1,'minu tekst',66);
INSERT INTO "comments" VALUES (137,1,'jälle rämpspost!',70);
INSERT INTO "comments" VALUES (138,1,'',68);
INSERT INTO "comments" VALUES (139,4,'rämpspostitususususususs',72);
INSERT INTO "comments" VALUES (140,1,'vaata kui tore!!!!',73);
COMMIT;
