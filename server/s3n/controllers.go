package s3n

type initcontrollerfunc func() error

func InitController(fnc initcontrollerfunc) {
	err := fnc()
	if err != nil {
		Log.Print(err)
		panic(err)
	}
}
