package app

import "github.com/KolYanCHiKK/dl_integration_service/internal/shared/types/app"

const (
	SubjectTypePublic   app.SubjectType = "public"
	SubjectTypePairwise app.SubjectType = "pairwise"
)

const (
	SigningAlgES256 app.SigningAlg = "ES256"
)

const (
	CodeChallengeAlgS256 app.CodeChallengeAlg = "S256"
)
