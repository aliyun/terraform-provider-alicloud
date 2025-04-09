package connectivity

import "strings"

func ramSamlProviderEncodedSamlMetadataDocumentEquivalence(document1, document2 string) (bool, error) {
	document1 = strings.ReplaceAll(document1, "\n", "")
	document2 = strings.ReplaceAll(document2, "\n", "")
	if document1 == document2 {
		return true, nil
	}
	return false, nil
}
