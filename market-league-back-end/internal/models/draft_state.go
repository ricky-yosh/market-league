package models

type LeagueState string

const (
	PreDraft  LeagueState = "pre_draft"
	InDraft   LeagueState = "draft"
	PostDraft LeagueState = "post_draft"
)
