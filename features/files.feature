# language: pt
Funcionalidade: Gerenciamento de Arquivos

  Contexto:
    Dado que estou autenticado com a API Key "minha-chave-secreta"

  Cenário: Listar todos os arquivos
    Quando faço uma requisição GET para "/files"
    Então a resposta deve ter status 200
    E a resposta deve conter uma lista de arquivos

  Cenário: Fazer upload de um arquivo
    Quando faço uma requisição POST para "/files/sendFiles" com o arquivo "teste.txt"
    Então a resposta deve ter status 201
    E a resposta deve conter o nome do arquivo "teste.txt"

  Cenário: Buscar arquivo por ID existente
    Dado que existe um arquivo com ID "algum"
    Quando faço uma requisição GET para "/files/ultimo"
    Então a resposta deve ter status 200
    E a resposta deve conter o arquivo com ID "${ultimo_id}"

  Cenário: Atualizar um arquivo existente
    Dado que existe um arquivo com ID "algum"
    Quando faço uma requisição PUT para "/files/ultimo" com o corpo:
      """
      {"fileName": "Novo Nome"}
      """
    Então a resposta deve ter status 200
    E a resposta deve conter o arquivo com nome "Novo Nome"

  Cenário: Remover um arquivo existente
    Dado que existe um arquivo com ID "algum"
    Quando faço uma requisição DELETE para "/files/ultimo"
    Então a resposta deve ter status 204
