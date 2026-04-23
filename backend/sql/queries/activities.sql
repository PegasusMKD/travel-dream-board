-- name: CreateActivity :one
insert into activities (url, title, image_url, board_uuid, user_uuid)
values (@url, @title, @image_url, @board_uuid, @user_uuid)
returning *;

-- name: GetActivityByUuid :one
select a.*,
    count(case when v.rank = 1 then 1 end)::int as likes,
    count(case when v.rank < 1 then 1 end)::int as dislikes
from activities a
left join votes v on v.voted_on_uuid = a.uuid and v.voted_on_ = 'activities'
where a.uuid = @uuid
group by a.uuid;

-- name: FindAllActivitiesByBoardUuid :many
select a.*,
    count(case when v.rank = 1 then 1 end)::int as likes,
    count(case when v.rank < 1 then 1 end)::int as dislikes
from activities a
left join votes v on v.voted_on_uuid = a.uuid and v.voted_on_ = 'activities'
where a.board_uuid = @board_uuid
group by a.uuid;

-- name: UpdateActivityByUuid :exec
update activities
set url = @url,
    title = @title,
    image_url = @image_url,
    notes = @notes,
    status = @status,
    booking_reference = @booking_reference,
    selected = @selected
where uuid = @uuid;


-- name: DeleteActivityByUuid :exec
delete from activities
where uuid = @uuid;
