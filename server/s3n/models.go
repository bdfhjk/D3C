package s3n

type initmodelfunc func() error

func InitModel(fnc initmodelfunc) {
	err := fnc()
	if err != nil {
		Log.Print(err)
		panic(err)
	}
}
