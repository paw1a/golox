package main

import "github.com/paw1a/golox/internal/interpreter"

func main() {
	interpreter.Run()
	//keysEvents, err := keyboard.GetKeys(10)
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	_ = keyboard.Close()
	//}()
	//
	//fmt.Println("Press ESC to quit")
	//for {
	//	event := <-keysEvents
	//	if event.Err != nil {
	//		panic(event.Err)
	//	}
	//	fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
	//	if event.Key == keyboard.KeyEsc {
	//		break
	//	}
	//}
}
