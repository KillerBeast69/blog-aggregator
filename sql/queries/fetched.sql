-- name: MarkFeedFetched :exec
update feeds 
set last_fetched_at = $2, updated_at = $3
where id = $1;

-- name: GetNextFeedToFetch :one
select * from feeds
order by last_fetched_at asc nulls first
limit 1;