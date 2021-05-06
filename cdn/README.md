# Configuração na CDN

Para preparar o ambiente dentro de uma CDN é necessário apontar o domínio para a mesma. Neste tutorial será mostrado como apontar para a CDN utilizada no Workshop, ou seja a GoCache, e também criar a regra de Rate Limit.

Após criada a conta e adicionado o domínio foi criada uma entrada raíz na aba “Websites & DNS”, no caso do projeto “fatecnologia.gocdn.com.br”, da seguinte forma:
![entrada](./imgs/cdn/entrada.png)

Após criado, o apontamento ficará da seguinte forma: 
![apontamento](./imgs/cdn/apontamento.png)

OBS.: o IP utilizado é apenas para demonstração, no caso é necessário apontar para o IP público de sua máquina (virtual ou não).

Após isso foi realizado o apontamento no provedor de DNS, no caso do ambiente realizado Registro.br, do modo que a entrada raiz apontasse para o conteúdo do “Apontamento de DNS” do painel GoCache (marcado em vermelho na imagem acima):
![registro-fatec](./imgs/cdn/registro-fatec.png)

Na aba de segurança do Painel GoCache, foi acessada a aba Rate Limit e foi criada uma regra para bloquear a máquina atacante. A regra de Rate Limit foi criada na CDN de forma que caso um determinado IP realize mais de 100 requisições em qualquer parte do site dentro de um período de 10 segundos, ele será bloqueado por 10 segundos. 
![regra-rl](./imgs/cdn/regra-rl.png)

Dessa forma o ambiente na CDN se encontra pronto para realização das simulações de ataques.

Inicialmente a regra pode ser criada como rascunho (Salvar como rascunho), para ativação posterior apenas no momento de utilização do Rate Limit da CDN.
