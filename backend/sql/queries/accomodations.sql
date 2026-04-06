-- name: FindAllAccomodationsByBoardUuid :many
select *
from accomodations
where board_uuid = @board_uuid;
