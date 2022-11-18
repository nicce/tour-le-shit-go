CREATE TABLE player (
	id VARCHAR(36),
	name VARCHAR(150),
	PRIMARY KEY(id)
);

CREATE TABLE scoreboard (
	points INT,
	player_id VARCHAR(36),
	holder_of_snek BOOLEAN,
	last_played VARCHAR(10),
	season INT,
	FOREIGN KEY(player_id) REFERENCES player(id)
);

CREATE TABLE score (
	id VARCHAR(36),
	player_id VARCHAR(36),
	points INT,
	holder_of_snek BOOLEAN,
	birdies INT,
	eagles INT,
	muligans INT,
	day VARCHAR(10),
	seasont INT,
	FOREIGN KEY(player_id) REFERENCES player(id)
);
