package draft

// DraftChannelProvider defines the methods needed to get a draft selection channel.
type DraftChannelProvider interface {
	GetDraftSelectionChannel(leagueID uint) chan uint
}
