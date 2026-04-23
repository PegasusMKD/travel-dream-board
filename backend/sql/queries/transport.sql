-- name: CreateTransport :one
insert into transport (url, title, image_url, board_uuid, user_uuid)
values (@url, @title, @image_url, @board_uuid, @user_uuid)
returning *;

-- name: GetTransportByUuid :one
select t.*,
    count(case when v.rank = 1 then 1 end)::int as likes,
    count(case when v.rank < 1 then 1 end)::int as dislikes
from transport t
left join votes v on v.voted_on_uuid = t.uuid and v.voted_on_ = 'transport'
where t.uuid = @uuid
group by t.uuid;

-- name: FindAllTransportByBoardUuid :many
select t.*,
    count(case when v.rank = 1 then 1 end)::int as likes,
    count(case when v.rank < 1 then 1 end)::int as dislikes
from transport t
left join votes v on v.voted_on_uuid = t.uuid and v.voted_on_ = 'transport'
where t.board_uuid = @board_uuid
group by t.uuid;

-- name: UpdateTransportByUuid :exec
update transport
set url = @url,
    title = @title,
    image_url = @image_url,
    notes = @notes,
    status = @status,
    booking_reference = @booking_reference,
    selected = @selected
where uuid = @uuid;


-- name: DeleteTransportByUuid :exec
delete from transport
where uuid = @uuid;
