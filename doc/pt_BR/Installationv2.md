GoLearn é uma biblioteca Go que é relativamente padrão, exceto por algumas dependências opcionais em C. Para instalá-la, você precisará [instalar o Go 1.2 ou posterior primeiro](http://golang.org/doc/install).

## Instalação

### Dependências do sistema
GoLearn usa [Gonum's BLAS wrapper](https://github.com/gonum/blas), que requer que o OpenBLAS ou similar esteja instalado.

#### Instalando o Gonum's BLAS wrapper

##### Instalando o OpenBLAS no Ubuntu

```bash
sudo apt-get install libopenblas-dev
```

##### Instalando o OpenSUSE 13.1
```bash
sudo zypper in blas-devel cblas-devel
```

#### Completando a instalação do Gonum's BLAS wrapper
Depois de instalar o Gonum's BLAS wrapper, digite os seguintes comandos
```bash
go get github.com/gonum/blas
cd $GOPATH/src/github.com/gonum/blas
go install ./...
```
No Ubuntu 14.04, modifique `$GOPATH/src/github.com/gonum/blas/cblas/blas.go` e modifique a linha que começa com `#cgo linux LDFLAGS:` para `#cgo linux LDFLAGS: -L/usr/lib -lopenblas`. /usr/lib deve ser o caminho para a biblioteca. Se não for, use `dpkg -L libopenblas-dev` para encontrar o caminho.

No OpenSUSE 13.1, modifique `$GOPATH/src/github.com/gonum/blas/cblas/blas.go` e modifique a linha que começa com `#cgo linux LDFLAGS:` para `#cgo linux LDFLAGS: -lblas -lcblas`.

### Instalando o liblinear no OpenSUSE 13.1 (64 bits)

1. Visite [o site do liblinear](http://www.csie.ntu.edu.tw/~cjlin/liblinear/) e baixe a versão 1.94. ([Link para baixar o zip](http://www.csie.ntu.edu.tw/~cjlin/cgi-bin/liblinear.cgi?+http://www.csie.ntu.edu.tw/~cjlin/liblinear+zip))
2. Extraia o arquivo, edite o arquivo para digitar `make lib`, para então construir a biblioteca.
3. Copie `liblinear.so.1` para `/usr/lib64`
4. Copie `linear.h` para `/usr/include`
5. Crie um link simbólico para `liblinear.so.1` em `/usr/lib64` usando `ln -s /usr/lib64/liblinear.so.1 /usr/lib64/liblinear.so`
6. Execute `ldconfig -v | grep linear` para verificar se a biblioteca está instalada.


### Instalando o liblinear no Ubuntu 14.04 (64 bits)
```bash
cd /tmp
wget http://www.csie.ntu.edu.tw/~cjlin/liblinear/oldfiles/liblinear-1.94.tar.gz
tar xf liblinear-1.94.tar.gz
cd liblinear-1.94 
make lib
sudo install -vm644 linear.h /usr/include 
sudo install -vm755 liblinear.so.1 /usr/lib 
sudo ln -sfv liblinear.so.1 /usr/lib/liblinear.so
```
### Instalando as dependências internas do GoLearn
Depois de instalar a linguagem Go e as dependências do sistema, digite:
```bash
go get -t -u -v github.com/sjwhitworth/golearn
cd $GOPATH/src/github.com/sjwhitworth/golearn/ext
go run make.go
```

Em seguida, você precisará adicionar `$GOPATH/src/github.com/sjwhitworth/golearn/ext/lib` a uma variável apropriada.

#### Bash (Linux)
Adicione ao seu arquivo `~/.bashrc`:
```bash
export LD_LIBRARY_PATH=$GOPATH/src/github.com/sjwhitworth/golearn/ext/lib:$LD_LIBRARY_PATH
```
Certifique-se de executar `source ~/.bashrc` ou reiniciar o shell antes de continuar.

#### Bash (Mac OS X)
Adicione ao seu arquivo `~/.bashrc`:
```bash
export DYLD_LIBRARY_PATH=$GOPATH/src/github.com/sjwhitworth/golearn/ext/lib:$DYLD_LIBRARY_PATH
```
Certifique-se de executar `source ~/.bashrc` ou reiniciar o shell antes de continuar.


### Completando a instalação
Execute o seguinte para completar a instalação
```bash
cd $GOPATH/src/github.com/sjwhitworth/golearn
go get ./...
```

## Problemas comuns
No Linux e no Mac OS X, você pode verificar se o Go está instalado corretamente por meio do terminal. [Tente compilar o programa de exemplo](http://golang.org/doc/install) para verificar se sua instalação está funcionando corretamente.

**Seu diretório `go` deve existir em seu diretório principal e ser gravável.**
Se não estiver, digite `cd && mkdir go` para criá-lo.

**Suas variáveis `GOPATH` e `PATH` devem estar configuradas corretamente.**
* Para verificar, digite `echo $GOROOT` e `echo $GOPATH`.
* Sua variável `$GOPATH` deve incluir seu `$GOROOT`, mais um diretório `bin/`. Por exemplo, se `$GOROOT` estiver definido como `/home/sen/go`, `$GOPATH` deve ser definido como `/home/sen/go/bin`
* Para garantir que essas variáveis estejam configuradas corretamente, adicione `export GOROOT=$HOME/go` e `export PATH=$PATH:$GOROOT/bin` ao seu arquivo de configuração do Bash.

## Support status
<table>
<tr>
<td>Sistemas operacionais</td><td>Mac OS X 10.8 <br /> Ubuntu 14.04 <br /> OpenSUSE 13.1</td></tr>
<tr><td>Versão do Go</td><td>1.2</td></tr>
<tr><td>Versão do GoLearn</td><td>Atual</td></tr>
<tr><td>Status do suporte</td><td>Atual</td>
<tr><td>Próxima revisão</td><td>Na atualização da versão</td></tr>
</table>

### Mac OS X 10.8

* A confirmação do funcionamento da instalação do BLAS através do HomeBrew ainda não foi realizada.
