package wl

type RatedQuery interface {
	Data() RateLimitData
}

type RateLimitQuery struct {
	RateLimitData RateLimitData
}

func (r *RateLimitQuery) Data() RateLimitData {
	return r.RateLimitData
}

type ExpansionEncountersQuery struct {
	WorldData     WorldData
	RateLimitData RateLimitData
}

func (e *ExpansionEncountersQuery) Data() RateLimitData {
	return e.RateLimitData
}

type ExpansionPartitionsQuery struct {
	WorldData     WorldDataByExpansionID
	RateLimitData RateLimitData
}

func (e *ExpansionPartitionsQuery) Data() RateLimitData {
	return e.RateLimitData
}
