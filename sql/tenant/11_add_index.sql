CREATE INDEX competition_id_index on player_score(competition_id);
CREATE INDEX competition_id_player_id_index on player_score(competition_id, player_id);
