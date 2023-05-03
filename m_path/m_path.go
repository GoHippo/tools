package m_path

import "strings"

// Func Slash_Moder(путь) исправляет слэш в пути, под linux формат
func ToSlash(path_original string) (p_slash string) {
	path_arr := strings.Split(path_original, "\\")
	for _, p := range path_arr {
		p_slash += p + "/"
	}
	return p_slash[:len(p_slash)-1]
}
