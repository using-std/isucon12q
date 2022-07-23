CREATE INDEX player_tenant_id_created_at_index on player(tenant_id, created_at DESC);
CREATE INDEX player_score_ranking_index on player_score(tenant_id, competition_id, row_num DESC);
CREATE INDEX player_score_individual_index on player_score(player_id, tenant_id,competition_id, row_num DESC);
