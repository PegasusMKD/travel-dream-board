-- name: CreateActivity :one
insert into activities (
    url, title, image_url, board_uuid, user_uuid,
    start_at, end_at, location
)
values (
    @url, @title, @image_url, @board_uuid, @user_uuid,
    @start_at, @end_at, @location
)
returning *;

-- name: GetActivityByUuid :one
select a.*,
    coalesce(avg(v.rank)::float, 0)::float as avg_rating,
    count(v.uuid)::int as rating_count
from activities a
left join votes v on v.voted_on_uuid = a.uuid and v.voted_on_ = 'activities'
where a.uuid = @uuid
group by a.uuid;

-- name: FindAllActivitiesByBoardUuid :many
select a.*,
    coalesce(avg(v.rank)::float, 0)::float as avg_rating,
    count(v.uuid)::int as rating_count
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
    selected = @selected,
    start_at = @start_at,
    end_at = @end_at,
    location = @location
where uuid = @uuid;


-- name: DeleteActivityByUuid :exec
delete from activities
where uuid = @uuid;
