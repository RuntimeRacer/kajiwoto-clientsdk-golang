package graphql

import "github.com/google/go-cmp/cmp"

/*
  Helper functions for this package
*/

func (e *AiDialogueInput) IsDuplicate(c *AiDialogueInput) bool {
	if e.Message == c.Message &&
		e.UserMessage == c.UserMessage &&
		e.Conditions.Equals(&c.Conditions) &&
		cmp.Equal(e.History, c.History) {
		return true
	}
	return false
}

func (e *AITrainingCondition) Equals(c *AITrainingCondition) bool {
	if e.ASM != c.ASM {
		if e.ASM != nil && c.ASM == nil ||
			e.ASM == nil && c.ASM != nil ||
			*e.ASM != *c.ASM {
			return false
		}
	}

	if e.Endearment != c.Endearment {
		if e.Endearment != nil && c.Endearment == nil ||
			e.Endearment == nil && c.Endearment != nil ||
			*e.Endearment != *c.Endearment {
			return false
		}
	}

	if e.Time != c.Time {
		if e.Time != nil && c.Time == nil ||
			e.Time == nil && c.Time != nil ||
			*e.Time != *c.Time {
			return false
		}
	}

	if e.Recent != c.Recent {
		if e.Recent != nil && c.Recent == nil ||
			e.Recent == nil && c.Recent != nil ||
			*e.Recent != *c.Recent {
			return false
		}
	}

	return true
}
