-- name: CreateVote :one
insert into votes (voted_by, rank, voted_on_, voted_on_uuid)
values (@voted_by, @rank, @voted_on_, @voted_on_uuid)
returning *;

-- name: FindAllVotesByVotedOnUuid :many
select v.*, u.name as user_name
from votes v
join users u on v.voted_by = u.uuid
where v.voted_on_ = @voted_on_ 
    and v.voted_on_uuid = @voted_on_uuid;

-- name: UpdateVoteByUuid :exec
update votes
set rank = @rank
where uuid = @uuid;

-- name: DeleteVoteByUuid :exec
delete from votes
where uuid = @uuid;
