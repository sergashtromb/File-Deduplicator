// модуль создан для объявления моделей и типов данных 
// для последующего обращения из других точек программы

package domain


type Duplicate struct {
	Id 		int16 `json:"id"`
	Paths 	[]string `json:"paths"`
}