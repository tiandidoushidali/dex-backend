package utility

func Setup() {
	err := initSequence()
	if err != nil {
		panic(err)
	}
}
