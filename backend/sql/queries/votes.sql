-- name: CreateVote :one
insert into votes (voted_by, rank, voted_on_, voted_on_uuid)
values (@voted_by, @rank, @voted_on_, @voted_on_uuid)
returning *;

-- name: FindAllVotesByVotedOnUuid :many
select *
from votes
where voted_on_ = @voted_on_ 
    and voted_on_uuid = @voted_on_uuid;

-- name: UpdateVoteByUuid :exec
update votes
set rank = @rank
where uuid = @uuid;

-- name: DeleteVoteByUuid :exec
delete from votes
where uuid = @uuid;
