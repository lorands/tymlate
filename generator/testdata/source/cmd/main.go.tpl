package main

func main() {

    println("Hi, {{.Env.USER}}! Your new password for today is: {{ randAlphaNum 27 }}")
}
