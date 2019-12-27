package scryfall

// A Set object represents a group of related Magic cards.
// All Card objects on Scryfall belong to exactly one set.
//
// Due to Magic’s long and complicated history, Scryfall includes many
// un-official sets as a way to group promotional or outlier cards together.
// Such sets will likely have a code that begins with 'p' or 't',
// such as 'pcel' or 'tori'.
//
// Official sets always have a three-letter set code, such as 'zen'.
type Set struct {
}
