package main

import (
  "net"
  "fmt"
  "strings"
  "time"
  "net/http"
  "io/ioutil"
	//linuxproc "github.com/c9s/goprocinfo/linux"
)

func sendToChannel(conn net.Conn, message string, channel string){
  fmt.Fprintf(conn, "PRIVMSG %s :%s\r\n", channel, message)
}

func urbanDict(conn net.Conn, message string) {
  resp, err := http.Get("http://api.urbandictionary.com/v0/define?term=test")
  if err != nil {
    sendToChannel(conn, "An error occurred reaching Urban Dictionary", "#bullpen")
    return
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  fmt.Printf("%s", string(body))
  sendToChannel(conn, string(body), "#bullpen")
}

func handleMessage(conn net.Conn, message string) {
  if strings.Contains(string(message), "PING"){
    fmt.Fprintf(conn, "PONG localhost.localdomain\r\n")
  }
  if strings.Contains(string(message), "ud?"){
    go urbanDict(conn, message)
  }
}

func receive(conn net.Conn) {
  for {
    message := make([]byte, 4096)
    length, err := conn.Read(message)
    if err != nil {
      conn.Close()
      break
    }
    if length > 0 {
      fmt.Printf("LENGTH: %d  RECEIVED: %s", length, string(message))
      handleMessage(conn, string(message))
    }
  }
}

func main(){


  fmt.Printf("Connecting...\n")
  conn, err := net.Dial("tcp", "sig2noise.net:6667")
  if err != nil {
    fmt.Printf("Error %s", err)
  }
  go receive(conn)
  
  fmt.Printf("Sending USER\n")
  fmt.Fprintf(conn, "USER gopher 0 * :Gopher\r\n")
  
  fmt.Printf("Sending NICK\n")
  fmt.Fprintf(conn, "NICK gopher\r\n")
  
  
  fmt.Printf("Sending JOIN\n")
  fmt.Fprintf(conn, "JOIN #bullpen\r\n")
  for {
    time.Sleep(10)
  }
}
