verifier=$(openssl rand -base64 48 | tr '+/' '-_' | tr -d '=\r\n')
challenge=$(printf '%s' "$verifier" | openssl dgst -sha256 -binary | openssl base64 -A | tr '+/' '-_' | tr -d '=\r\n')
printf 'verifier:  %s\nchallenge: %s\n' "$verifier" "$challenge"