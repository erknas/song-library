CREATE TABLE IF NOT EXISTS songs (
	id SERIAL PRIMARY KEY,
	song VARCHAR(255) NOT NULL,
	group_name VARCHAR(255) NOT NULL,
	release_date DATE NOT NULL,
	text TEXT,
	link VARCHAR(255)
);

CREATE INDEX idx_song_group_date ON songs(song, group_name, release_date);

CREATE INDEX idx_group_date ON songs(group_name, release_date);

CREATE INDEX idx_date ON songs(release_date);
