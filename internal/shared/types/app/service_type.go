package app

type SubjectType = string /* C какими типами subject работает сервис авторизации:
одинаковые для все или есть отличия под клиента */

type SigningAlg = string // Поддерживаемые алгоритмы шифрования токенов

type CodeChallengeAlg = string // Поддерживаемые алгоритмы для PKCE
