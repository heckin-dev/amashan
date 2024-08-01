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

type CharacterParsesQuery struct {
	CharacterData CharacterData
	RateLimitData RateLimitData
}

func (c *CharacterParsesQuery) Data() RateLimitData {
	return c.RateLimitData
}

func (c *CharacterParsesQuery) ToDTO() *CharacterParseDTO {
	dto := &CharacterParseDTO{}
	if c.CharacterData.Character.Hidden != nil {
		dto.Hidden = Bool(bool(*c.CharacterData.Character.Hidden))
	}

	dto.ZoneRankings = c.CharacterData.Character.ZoneRankings.DTO()

	return dto
}

func Bool(v bool) *bool {
	return &v
}
