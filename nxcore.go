package main

func main() {
	/*
		testString := "action:bake<;>arg:cake"

		testMessage := Message{}
		testMessage.FromString(testString)
		fmt.Println("cake")
		for k, v := range testMessage.Data {
			fmt.Printf("%v:%v\n", k, v)
		}
		fmt.Println(testMessage.Data["origami"] == "")
	*/
	testServer := Server{}
	testServer.Start()
}
