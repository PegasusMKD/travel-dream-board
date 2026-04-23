-- Dedupe: keep the most recent row per (voted_by, voted_on_, voted_on_uuid)
delete from votes
where uuid in (
    select uuid from (
        select uuid,
            row_number() over (
                partition by voted_by, voted_on_, voted_on_uuid
                order by updated_at desc, created_at desc
            ) as rn
        from votes
    ) ranked
    where rn > 1
);

-- Clamp existing rank values into the new 1..5 range before the CHECK.
-- Old semantics used 1 (upvote) and -1 (downvote); map downvotes to 1 star.
update votes set rank = 1 where rank < 1;
update votes set rank = 5 where rank > 5;

alter table votes
    add constraint votes_unique_voter_per_item
    unique (voted_by, voted_on_, voted_on_uuid);

alter table votes
    add constraint votes_rank_range check (rank between 1 and 5);
