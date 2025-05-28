//frontend

package seoscraper

import basicgroupsprotocol "github.com/futura-platform/protocol/basicgroups/protocol"

type SearchTerm string

// Equals implements basicgroupsprotocol.Parsable.
func (s SearchTerm) Equals(s2 SearchTerm) bool {
	return s == s2
}

// GetGroupConfig implements basicgroupsprotocol.Parsable.
func (SearchTerm) GetGroupConfig() basicgroupsprotocol.GroupConfig {
	return basicgroupsprotocol.GroupConfig{
		EntryTypeSingular: "Search Term",
		EntryTypePlural:   "Search Terms",
		EntryPlaceholder:  "Enter a search term",
		Icon: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-6">
  <path stroke-linecap="round" stroke-linejoin="round" d="m15.75 15.75-2.489-2.489m0 0a3.375 3.375 0 1 0-4.773-4.773 3.375 3.375 0 0 0 4.774 4.774ZM21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z" />
</svg>
`,
	}
}

// ParseEntry implements basicgroupsprotocol.Parsable.
func (SearchTerm) ParseEntry(s string) (SearchTerm, error) {
	return SearchTerm(s), nil
}

// SerializeEntry implements basicgroupsprotocol.Parsable.
func (s SearchTerm) SerializeEntry() string {
	return string(s)
}
