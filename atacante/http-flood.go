package main

import (
  "crypto/tls"
  "flag"
  "fmt"
  "math/rand"
  "net/http"
  "net/url"
  "os"
  "regexp"
  "strings"
  "time"
)

// Define amostragem de caracters para gerar requisições aleatorias
var lettersDigitsHyphen []byte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-")

func stringGenerator() string {
  minStringSize := 4
  maxStringSize := 10

  stringSize := rand.Intn(maxStringSize-minStringSize) + minStringSize

  bytesCharacters := make([]byte, stringSize)

  for i := range bytesCharacters {
    bytesCharacters[i] = lettersDigitsHyphen[rand.Intn(len(lettersDigitsHyphen))]
  }

  return string(bytesCharacters)
}

func makeHttpConnections(clientBase *http.Client, requestBase *http.Request, connection int, httpMode string, domainBase string, pathBase string) {
  for {
    for i := 0; i < connection; i++ {
      randomParameter := stringGenerator()

      urlRequest, err := url.Parse(httpMode + "://" + domainBase + "/" + pathBase + randomParameter)

      if err != nil {
        fmt.Println(err)
        continue
      } else {
        requestBase.URL = urlRequest
        go clientBase.Do(requestBase)
        fmt.Println(requestBase.URL.String())
      }
      time.Sleep(time.Duration(1000000000 / connection))
    }
    //time.Sleep(time.Second * 1)
  }
}

func main() {

  // Define parametros do programa
  threads := flag.Int("threads", 1, "Numero de threads a serem utilizados no flood. Padrão 1")
  conections := flag.Int("connections", 3, "numero de conexões por thread")
  timeToRun := flag.Int("duration", 30, "tempo em segundos que a ferramenta vai executar(zero = infinito)")
  url := flag.String("url", " ", "Url base do teste de sobrecarga(*OBRIGATÓRIO*)")

  flag.Parse()

  if *url == " " {
    fmt.Println("Parametro '--url' não foi definido")
    flag.PrintDefaults()
    os.Exit(2)
  } else {
    *url = strings.ToLower(*url)
  }

  // Variaveis padrão do cliente http
  headerAccept := [2]string{"Accept", "*/*"}
  headerUserAgent := [2]string{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"}
  headerReferers := [2]string{"Referer", "https://google.com/"}
  headerNoCache := [2]string{"Cache-Control", "no-cache"}
  headerKeepAlive := [2]string{"Connection", "Keep-Alive"}

  httpMode := strings.Split(*url, ":")[0]

  domain := *url

  httpModePattern := regexp.MustCompile("http.*://")
  // pathPattern := regexp.MustCompile("/.+$")

  domain = httpModePattern.ReplaceAllString(domain, "")
  urlSplited := strings.Split(domain, "/")

  domain = urlSplited[0]
  path := urlSplited[1]

  // fmt.Println(domain)
  // fmt.Println(path)

  switch httpMode {
  case "http":
    fmt.Println("Conexão HTTP será utilizada.")
  case "https":
    fmt.Println("Conexão HTTPS utilizada.")
  default:
    fmt.Println("Não foi possivel encontrar um modo http valido na URL. HTTPS será utilizado como padrão.")
    httpMode = "https"
  }

  client := &http.Client{}

  if httpMode == "https" {
    tr := &http.Transport{
      TLSClientConfig: &tls.Config{
        InsecureSkipVerify: true,
        ServerName:         domain,
      },
    }
    client = &http.Client{Transport: tr}
  }

  req, err := http.NewRequest("GET", domain, nil)

  if err != nil {
    panic(err)
  }

  req.Header.Add(headerAccept[0], headerAccept[1])
  req.Header.Add(headerUserAgent[0], headerUserAgent[1])
  req.Header.Add(headerReferers[0], headerReferers[1])
  req.Header.Add(headerNoCache[0], headerNoCache[1])
  req.Header.Add(headerKeepAlive[0], headerKeepAlive[1])

  for i := 0; i < *threads; i++ {
    go makeHttpConnections(client, req, *conections, httpMode, domain, path)
  }

  if *timeToRun == 0 {
    for {
    }
  } else {
    time.Sleep(time.Duration(*timeToRun) * time.Second)
  }
}
