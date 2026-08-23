// модуль создан для объявления моделей и типов данных 
// для последующего обращения из других точек программы

package domain


type Duplicate struct {
	Id 		int16 			`json:"id"`
	Files 	[]*FoundFile 	`json:"files"`
}

type FoundFile struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
	Path string `json:"path"`
}