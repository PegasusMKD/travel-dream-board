alter table votes drop constraint if exists votes_rank_range;
alter table votes drop constraint if exists votes_unique_voter_per_item;
