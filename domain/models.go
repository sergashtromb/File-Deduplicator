// модуль создан для объявления моделей и типов данных 
// для последующего обращения из других точек программы

package domain


type Duplicate struct {
	Hash 	[]byte 			`json:"hash"`
	Size 	int64 			`json:"size"`
	Files 	[]*FoundFile 	`json:"files"`
}

type FoundFile struct {
	Name string `json:"name"`
	Hash []byte `json:"hash"`
	Path string `json:"path"`
	Size int64 	`json:"size"`
}

type ParamsWorker struct {
	Path 			string 
	IsSaved 		bool
	UsRecurSeach 	bool
}
