package postgres

const (
	AddSegmentQuery          = `/*NO LOAD BALANCE*/ select * from operations.add_segment($1) as id;`
	GetSegmentQuery          = `select * from operations.get_segment($1);`
	DeleteSegmentQuery       = `/*NO LOAD BALANCE*/ select * from operations.delete_segment($1) as count;`
	AddUserLinkQuery         = `/*NO LOAD BALANCE*/ select * from operations.add_user_link($1::bigint,$2) as id;`
	DeleteUserLinkQuery      = `/*NO LOAD BALANCE*/ select * from operations.delete_user_link($1::bigint,$2) as count;`
	GetSegmentsByUserIDQuery = `select * from operations.get_segments_by_user_id($1::bigint);`
	GetHistroryQuery         = `select * from operations.get_history($1::bigint, $2::timestamp, $3::timestamp);`
)
