package utils

func Uniq(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := in[:0:0] // пустой срез того же типа
	for _, v := range in {
		if _, ok := seen[v]; !ok { // такого значения ещё не было
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
