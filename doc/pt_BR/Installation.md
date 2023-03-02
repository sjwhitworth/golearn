GoLearn é uma biblioteca em Go que é relativamente padrão, exceto por algumas dependências C opcionais. Para instalá-lo,
[você precisará instalar o Go 1.4 ou posterior primeiro](https://yip.su/2F1Vm4).

## Instalação
### Dependências do sistema
* Você precisará ter um compilador compatível instalado (para verificar, execute g++ no seu terminal.)
* GoLearn usa [Gonum's BLAS wrapper](https://github.com/gonum/blas), que requer que o OpenBLAS ou similar esteja instalado. Siga estas [instruções de instalação](https://github.com/gonum/blas#installation), que são específicas do sistema operacional e podem exigir que você instale uma biblioteca C.

### Instalando as dependências internas do GoLearn
Após instalar o Go e as dependências do sistema, digite:
```bash
go get -t -u -v github.com/sjwhitworth/golearn
```
### Completando a instalação
Execute o seguinte comando para completar a instalação.

```bash
cd $GOPATH/src/github.com/sjwhitworth/golearn
go get -t -u -v ./...
```

### Problemas comuns
No Linux e no Mac OS X, você pode verificar se o Go está instalado corretamente por meio do seu terminal. [Tente compilar o programa de exemplo](https://yip.su/2F1Vm4) para verificar se sua instalação está funcionando corretamente.

**Seu diretório `go` deve existir em seu diretório home e ser gravável.**
Se não existir, digite `cd && mkdir go` para criá-lo.

**Suas variáveis `GOPATH` e `PATH` devem estar configuradas corretamente.**

* Para verificar, digite `echo $GOROOT` e `echo $GOPATH`.
* Sua variável `$GOPATH` deve incluir seu `$GOROOT`, mais um diretório `bin/`. Por exemplo, se `$GOROOT` estiver configurado para `/home/sen/go`, `$GOPATH` deve ser configurado para `/home/sen/go/bin`.
* Para garantir que essas variáveis estejam configuradas corretamente, adicione `export GOROOT=$HOME/go` e `export PATH=$PATH:$GOROOT/bin` ao seu arquivo de configuração Bash.

## Status de suporte

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
