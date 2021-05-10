# Configurar a hospedagem
>Exemplos baseados em um sistema operacional Centos 7.

## Configurações iniciais
- Atualizamos o gerenciador de pacotes:
```
sudo yum update -y
```
- Instalamos o  repositório oficial de pacotes pré compilados do Centos:
```
sudo yum install epel-release
```
- Instalamos pré requisitos necessários para configurar repositórios externos:
```
sudo yum install yum-utils -y
```

## Instalar o servidor WEB Nginx

- Criamos um arquivo referenciando o repositório oficial do Nginx para o centos 7:
```
sudo echo "[nginx-stable]
name=nginx stable repo
baseurl=http://nginx.org/packages/centos/$releasever/$basearch/
gpgcheck=1
enabled=1
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true

[nginx-mainline]
name=nginx mainline repo
baseurl=http://nginx.org/packages/mainline/centos/$releasever/$basearch/
gpgcheck=1
enabled=0
gpgkey=https://nginx.org/keys/nginx_signing.key
module_hotfixes=true" >> /etc/yum.repos.d/nginx.repo
```
- Habilitamos o repositório no gerenciador de pacotes:
```
sudo yum-config-manager --enable nginx-mainline
```
- Instalamos o Nginx:
```
sudo yum install nginx -y
```
- Configuramos o servidor Web para ser iniciado junto ao sistema operacional:
```
sudo systemctl enable nginx 
```
- Iniciamos o nginx e verificamos o status do serviço:
```
sudo systemctl nginx start
sudo systemctl nginx status
```
> Ref: http://nginx.org/en/linux_packages.html#RHEL-CentOS

## instalar o PHP 7.4

- Instalamos o repositório "Remi's" para versões 7.x do PHP:
```
sudo yum install http://rpms.remirepo.net/enterprise/remi-release-7.rpm -y
```
- Habilitamos a versão 7.4 do PHP para ser instalado:
```
sudo yum-config-manager --enable remi-php74	
```
- Instalamos o PHP e suas extensões necessárias para rodar o Wordpress:
```
sudo yum install php74-php php74-php-fpm php74-php-gd php74-php-json php74-php-mbstring php74-php-mysqlnd php74-php-xml php74-php-xmlrpc php74-php-opcache -y
```
> Ref: https://make.wordpress.org/hosting/handbook/handbook/server-environment/#php-extensions

- Configuramos o PHP-fpm para ser executado com o usuário e grupo do nginx:
```
sudo vim /etc/opt/remi/php74/php-fpm.d/www.conf 
```
Encontre as linhas com o parâmetro "user","group" e "listen". Deixe essas opções com os mesmo valores do exemplo abaixo:
```
user = nginx
group = nginx
liste = 127.0.0.1:9000
```
- Habilite o serviço do php-fpm para ser iniciado com o sistema operacional:
```
sudo systemctl enable php74-php-fpm 
```
- Inicie o serviço do php-fpm e verifique o status:
```
sudo systemctl start php74-php-fpm 
sudo systemctl status php74-php-fpm 
```
## Instalar o MySql 5.7
 - Instalamos o repositório do MySQl 5.7:
 ```
sudo yum install https://dev.mysql.com/get/mysql57-community-release-el7-9.noarch.rpm -y
 ```
- Instalamos o MySql 5.7 Server:
```
sudo yum install mysql-server 
```
- Habilitamos o MySQL e iniciamos o serviço:
```
sudo systemctl enable mysqld
sudo systemctl start mysqld
sudo systemctl status mysqld
```
- Quando instalado o MySql 5.7 gera uma senha temporaria. Devemos coletar esta senha para que possamos seguir com a configuração inicial:
```
sudo grep --color=auto 'temporary password' /var/log/mysqld.log
```
Uma linha com o padrão abaixo deve ser apresentada ao executar o comando acima:
```
2016-12-01T00:22:31.416107Z 1 [Note] A temporary password is generated for root@localhost: mqRfBU_3Xk>r
```
Com a senha temporária em mãos podemos seguir com a configuração do MySql:
```
sudo mysql_secure_installation
```
- Será solicitada a senha do usuário root do MySql. Informe a senha temporária que foi encontrada anteriormente:

```
Securing the MySQL server deployment.

Enter password for user root:
```
- Após isso será solicitada uma nova senha para o usuário root que irá substituir a senha temporária, crie uma senha com mais de 10 caracteres com números letras e ao menos um caracter especial:

```
The existing password for the user account root has expired. Please set a new password.

New password:

Re-enter new password:
```
> Outras perguntas serão feitas pelo script de configuração. Responda todas as perguntas com "yes".

- Por fim, podemos fazer login no MySql e criar um banco de dados a ser utilizado pelo Wordpress. Para fazer login basta digitar o comando abaixo no terminal e informar a senha configurada anteriormente quando solicitado:

```
mysql -u root -p
```
- Criamos um banco de dados chamado "fatecnologia" porém você pode utilizar o nome que achar mais apropriado a sua necessidade:
```
create database fatecnologia;
```
- Configuramos um usuário chamado "fatecnologia" com a senha de acesso "bzwYI483uGMs/" e permissão de acesso ao banco de dados "fatecnologia" que acabou de ser criado:
```
create user 'fatecnologia'@'%' identified by 'bzwYI483uGMs/';
grant all privileges on fatecnologia.* to 'fatecnologia'@'%';
flush privileges;  
```
>ref: https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-centos-7
## Instalar o Wordpress
- Baixamos o wordpress em **"/var/www/html"**:
```
sudo cd /var/www/html/
sudo wget https://br.wordpress.org/latest-pt_BR.tar.gz
```
- Descompactamos o download e renomeamos a pasta do Wordpress:

```
sudo tar -xvf latest-pt_BR.tar.gz
sudo mv wordpress fatecnologia
```
- Configuramos as permissões do Wordpress.
```
sudo chown nginx:nginx fatecnologia/ -R
sudo sudo chmod 755 fatecnologia/
```
## Configurar um website Wordpress no Nginx:
> Ref: https://www.nginx.com/resources/wiki/start/topics/recipes/wordpress/

Para finalizar a configuração do servidor Web devemos criar um arquivo de configuração do Nginx que vai definir os parâmetros para o funcionamento do site Wordpress.

- Criamos o arquivo "fatecnologia.gocdn.com.br.conf" na pasta "/etc/nginx/conf.d":
```
sudo vim /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf
```
O conteúdo do arquivo arquivo deve seguir o padrão abaixo:

```
# Informa o IP e a porta do servidor PHP
upstream php {
        server 127.0.0.1:9000;
}

server {

        # Porta que o Nginx vai utilizar para receber as requisições:
        listen 80;

        # Configura uma pasta para os logs de acesso/erro do site 
        error_log /var/log/fatecnologia/fatecnologia.log ;
        access_log /var/log/fatecnologia/fatecnologia.log ;

        ## O nome do seu website
        server_name fatecnologia.gocdn.com.br ;
        ## A pasta contendo os arquivos do site.
        root /var/www/html/fatecnologia ; 
        ## Configura qual o arquivo de index padrão do seu site
        index index.php;

        location = /favicon.ico {
                log_not_found off;
                access_log off;
        }

        location = /robots.txt {
                allow all;
                log_not_found off;
                access_log off;
        }

        location / {
                # This is cool because no php is touched for static content.
                # include the "?$args" part so non-default permalinks doesn't break when using query string
                try_files $uri $uri/ /index.php?$args;
        }

        location ~ \.php$ {
                #NOTE: You should have "cgi.fix_pathinfo = 0;" in php.ini
                include fastcgi_params;
                fastcgi_intercept_errors on;
                fastcgi_pass php;
                #The following parameter can be also included in fastcgi_params file
                fastcgi_param  SCRIPT_FILENAME $document_root$fastcgi_script_name;
        }

}
```

- Após configurar este arquivo basta criar a pasta de destino dos logs de acesso/erro do website e verificar se a sintaxe está correta para recarregar o Nginx:

```
sudo mkdir -p /var/log/fatecnologia/
sudo chown nginx:nginx /var/log/fatecnologia/ -R
sudo chmod 755 /var/log/fatecnologia/ -R
sudo nginx -T && sudo nginx -s reload
```
Com essa operação finalizada, você pode fazer um acesso ao site no seu navegador para verificar se o acesso está operacional.

Prossiga com a configuração inicial do Wordpress antes de prosseguir.

> Ref: https://wordpress.org/support/article/how-to-install-wordpress/

## Gerar um certificado SSL Let's Encrypt

> Para gerar um certificado Let's Encrypt é necessário que a sua hospedagem tenha um domínio/subdomínio que esteja apontando para o IP público da sua hospedagem.


Para gerar o certificado Let 's Encrypt, utilizaremos o cert bot que é um programa de linha de comando que pode verificar de maneira automática os sites configurados no seu Nginx e gerar e configurar um certificado SSL para estes Websites.


> Ref: https://certbot.eff.org/

A instalação pode ser feita diretamente utilizando o comando abaixo:

```
sudo yum install certbot certbot-nginx -y
```

Para configurar gerar e configurar o certificado SSL:
```
sudo certbot --nginx
```

Algumas informações iniciais serão solicitadas pelo certbot. Responda de acordo com suas preferências:

1. E-mail para contato quando o SSL estiver para expirar.
```
Enter email address (used for urgent renewal and security notices)
 (Enter 'c' to cancel): mpatz@gocache.com.br
```
2. Aceitar ou negar os termos do serviço.
```
Please read the Terms of Service at
https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf. You must
agree in order to register with the ACME server. Do you agree?
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: y
```
3. Compartilhar o seu email com a "Electronic Frontier Foundation" para receber novidades e campanhas da organização.
```
Would you be willing, once your first certificate is successfully issued, to
share your email address with the Electronic Frontier Foundation, a founding
partner of the Let's Encrypt project and the non-profit organization that
develops Certbot? We'd like to send you email about our work encrypting the web,
EFF news, campaigns, and ways to support digital freedom.
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
(Y)es/(N)o: n
```
4. Informe o numero do domínio listado que deseja gerar o certificado SSL:
```
Which names would you like to activate HTTPS for?
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
1: fatecnologia.gocdn.com.br
- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
Select the appropriate numbers separated by commas and/or spaces, or leave input
blank to select all options shown (Enter 'c' to cancel): 1
```

Uma vez que o domínio é informado o SSL deve ser gerado e configurado automaticamente.

## Configurar o Rate Limit no Nginx

Copiamos o arquivo de configuração atual:
```
sudo cp /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf-ratelimit
```
Movemos o arquivo de configuração da hospedagem:
```
sudo mv /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf-default 
```
Agora devemos ser capazes de visualizar os dois arquivos de configuração:
```
sudo ls -l /etc/nginx/conf.d/
default.conf
fatecnologia.gocdn.com.br.conf-default
fatecnologia.gocdn.com.br.conf-ratelimit
```
Editamos o arquivo finalizado em "-ratelimit" e adicionamos as configurações de limite por IP:
```
sudo vim fatecnologia.gocdn.com.br.conf-ratelimit
```
Acima da linha que define a configuração do servidor web **"server {"**, adicionamos as linhas abaixo:
```
#Configura uma zona de limitação de acessos chamada "limiteComum".
limit_req_zone $binary_remote_addr zone=limiteComum:10m rate=20r/s;
```
1. **limit_req_zone**:
   * Modulo que limita o processamento de requisições baseado em uma variavel que for definida.
2. **$binary_remote_addr**:
   * Variavel com o tamanho de 4 Bytes para endereços IPv4.
3. **zone=limiteComum:10m**:
   * Nome da zona de Rate Limit "limiteComum" configurada("limiteComum") e tamanho em memória RAM dedicada para armazenar o estado das requisições que serão processadas.
4. **rate=20r/s**:
   * Define a taxa limite de requisições por segundo.

> Ref: http://nginx.org/en/docs/http/ngx_http_limit_req_module.html

Dentro do bloco **server {}**, adicionamos as duas linhas abaixo:
```
limit_req zone=limiteComum ;
limit_req_status 429 ;
```
1. **limit_req zone=limiteComum;**:
   * Informa qual zona de Rate Limit deve ser processada para este Website.
2. **limit_req_status 429;**:
   * Configura o status HTTP 429 como resposta de requisições que ultrapassarem o limite configurado.

Para habilitar e desabilitar a configuração de Rate Limit do Website configurado, basta criar um link simbolico para a configuração desejada e recarregar o serviço.

- Para ativar a configuração de ratelimit:
```
sudo ln -sf /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf-ratelimit /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf
sudo nginx -t && nginx -s reload
```
- Para desativar a configuração de Rate Limit:
```
sudo ln -sf /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf-default /etc/nginx/conf.d/fatecnologia.gocdn.com.br.conf
sudo nginx -t && nginx -s reload
```

## Instalar e configurar o Netdata para gerar métricas

- O netdata possui diversos módulos compilados durante a sua instalação e por esse motivo dispõe de um script de instalação automática com uma linha. Para instalar o netdata basta executar o comando abaixo:


```
sudo bash <(curl -Ss https://my-netdata.io/kickstart.sh)
```
> Ref: https://learn.netdata.cloud/docs/agent/packaging/installer
Aguarde a instalação finaliza...

Após a instalação finalizar é necessário configurar o Netdata para verificar os arquivos de log do Website. Para isso utilizamos o plugin baseado na linguagem Golang que vem instalado no Netdata por padrão.

1. Devemos desativar o módulo de verificação de logs baseado em Python. Para isso basta acessar a pasta de configuração do Netdata e utilizar a ferramenta de edição para encontrar a linha que contém o conteúdo **"# web_log: yes"** e alterá-la para **"web_log: no"**
```
sudo cd /etc/netdata/
sudo ./edit-config python.d.conf
.
.
.
web_log: no
```

2. Ativamos o módulo **"web_log"** no arquivo **"go.d.conf"**. Para isso é necessário abrir o arquivo para edição e encontrar a linha **"# web_log: no"** e alterá-la para **"web_log: yes"**:
```
sudo ./edit-config go.d.conf
.
.
.
web_log: yes
```
3. Por último devemos configurar um "job" no Netdata informando um nome e o arquivo de log que deve ser verificado. Abra o arquivo **"go.d/web_log.conf"** para edição e no fim do arquivo adicione o conteúdo abaixo:
```
# Abrir o arquivo para edição:
sudo ./edit-config go.d/web_log.conf
.
.
.
# Conteúdo a ser adicionado no fim do arquivo:

  - name: nginx
    path: /var/log/fatecnologia/fatecnologia.log 

```

4. Agora basta habilitar e reiniciar o serviço do Netdata:
```
systemctl enable netdata
systemctl restart netdata
```

Com isso você pode acessar "http://IP_DA_HOSPDAGEM:19999" para verificar as métricas do Netdata sendo geradas.

> Ref: https://learn.netdata.cloud/guides/collect-apache-nginx-web-logs

### Configurar um dashboard customizado no netdata

Para este laboratório criamos um dashboard customizado para o netdata. Caso tenha interesse em configurá-lo basta colocar o conteúdo do arquivo [custom.html](./custom.html) na pasta **"/usr/share/netdata/web/"** em um arquivo chamado **custom.html**:

```
vim /usr/share/netdata/web/custom.html
sudo chown netdata:netdata /usr/share/netdata/web/custom.html
```
Agora basta acessar o Netdata referenciando o arquivo em seu navegador:

**"http://IP_DA_HOSPDAGEM:19999/custom.html"**
>Ref: https://learn.netdata.cloud/guides/step-by-step/step-08