package normalizer

import "github.com/marcinbohm/search-index-lint/internal/model"

type FieldStats struct {
	Properties    int
	MultiFields   int
	RuntimeFields int
	TotalFields   int
}

func CountFields(corpus Corpus) FieldStats {
	var stats FieldStats
	WalkFields(corpus, func(visit model.FieldVisit) {
		switch visit.Role {
		case model.FieldRoleProperty:
			stats.Properties++
		case model.FieldRoleMultiField:
			stats.MultiFields++
		case model.FieldRoleRuntimeField:
			stats.RuntimeFields++
		}
	})
	stats.TotalFields = stats.Properties + stats.MultiFields + stats.RuntimeFields
	return stats
}
