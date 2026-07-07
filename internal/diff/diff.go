package diff

import "github.com/marcinbohm/search-index-preflight/internal/model"

func Compare(base model.Corpus, current model.Corpus) Result {
	baseResources := collectResources(base)
	currentResources := collectResources(current)

	resourceIDs := make(map[ResourceID]struct{}, len(baseResources)+len(currentResources))
	for resourceID := range baseResources {
		resourceIDs[resourceID] = struct{}{}
	}
	for resourceID := range currentResources {
		resourceIDs[resourceID] = struct{}{}
	}

	var changes []FieldChange
	for resourceID := range resourceIDs {
		changes = append(changes, compareResource(resourceID, baseResources[resourceID], currentResources[resourceID])...)
	}

	sortFieldChanges(changes)
	return Result{FieldChanges: changes}
}

func compareResource(resourceID ResourceID, baseFields, currentFields map[FieldID]FieldSnapshot) []FieldChange {
	fieldIDs := make(map[FieldID]struct{}, len(baseFields)+len(currentFields))
	for fieldID := range baseFields {
		fieldIDs[fieldID] = struct{}{}
	}
	for fieldID := range currentFields {
		fieldIDs[fieldID] = struct{}{}
	}

	var changes []FieldChange
	for fieldID := range fieldIDs {
		before, hadBefore := baseFields[fieldID]
		after, hasAfter := currentFields[fieldID]

		switch {
		case !hadBefore && hasAfter:
			afterCopy := after
			changes = append(changes, FieldChange{
				Kind:     ChangeFieldAdded,
				Resource: resourceID,
				Field:    fieldID,
				After:    &afterCopy,
			})
		case hadBefore && !hasAfter:
			beforeCopy := before
			changes = append(changes, FieldChange{
				Kind:     ChangeFieldRemoved,
				Resource: resourceID,
				Field:    fieldID,
				Before:   &beforeCopy,
			})
		case hadBefore && hasAfter && before.Type != after.Type:
			beforeCopy := before
			afterCopy := after
			changes = append(changes, FieldChange{
				Kind:     ChangeFieldTypeChanged,
				Resource: resourceID,
				Field:    fieldID,
				Before:   &beforeCopy,
				After:    &afterCopy,
			})
		}
	}
	return changes
}
