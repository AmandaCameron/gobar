package main

func failMeMaybe(err error) {
	if err != nil {
		panic(err)
	}
}
