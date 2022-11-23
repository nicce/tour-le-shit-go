CREATE TABLE player (
	id VARCHAR(36),
	name VARCHAR(150),
	PRIMARY KEY(id)
);

CREATE TABLE score (
	id VARCHAR(36),
	player_id VARCHAR(36),
	points INT,
	birdies INT,
	eagles INT,
	muligans INT,
	day VARCHAR(10),
	season INT,
	FOREIGN KEY(player_id) REFERENCES player(id) ON DELETE CASCADE
);
